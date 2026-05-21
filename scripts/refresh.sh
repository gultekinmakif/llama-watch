#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# Throwaway pretty-print helpers for screenshot output.
c_blue=$'\033[1;34m'
c_cyan=$'\033[36m'
c_green=$'\033[1;32m'
c_yellow=$'\033[33m'
c_dim=$'\033[2m'
c_reset=$'\033[0m'

section() { printf '\n%s==>%s %s%s%s\n' "$c_blue" "$c_reset" "$c_blue" "$1" "$c_reset"; }
step()    { printf '    %s·%s %s\n' "$c_cyan" "$c_reset" "$1"; }
done_()   { printf '    %s✓%s %s%s%s\n' "$c_green" "$c_reset" "$c_dim" "$1" "$c_reset"; }
warn()    { printf '    %s!%s %s\n' "$c_yellow" "$c_reset" "$1"; }
T_START=$SECONDS

# Fail loud on missing toolchain rather than midway through a parallel fetch.
command -v bun >/dev/null 2>&1 || { echo "refresh.sh: bun not on PATH; install bun and retry" >&2; exit 1; }

# 1. Load env (DATABASE_URL, UPSTREAM_DIR, PROTOCOLS_JSON, SNAPSHOT_OUT).
[ -f .env ] && source .env

: "${UPSTREAM_REMOTE:=https://github.com/DefiLlama}"
: "${UPSTREAM_DIR:=./var/upstream}"
: "${PROTOCOLS_JSON:=./var/snapshot/protocols.json}"
: "${SNAPSHOT_OUT:=./var/snapshot/snapshot.json}"


# 2. Refresh upstream inputs in parallel: sparse clone of defillama-server
#    plus the two release manifests the build step joins against. Each lands
#    in its own path so the background jobs share no state.
SERVER_DIR="$UPSTREAM_DIR/defillama-server"
SNAPSHOT_DIR="$(dirname "$SNAPSHOT_OUT")"
TVL_MODULES="$SNAPSHOT_DIR/tvlModules.json"
DIM_MODULES="$SNAPSHOT_DIR/dimensionModules.json"
mkdir -p "$SNAPSHOT_DIR"

section "Fetch upstream (parallel)"
step "sparse clone DefiLlama/defillama-server (defi/src/protocols)"
step "curl DefiLlama-Adapters/latest/tvlModules.json"
step "curl dimension-adapters/latest/dimensionModules.json"
T_FETCH=$SECONDS

pids=()

(
  if [ -d "$SERVER_DIR/.git" ]; then
    git -C "$SERVER_DIR" fetch --depth=1 -q
    git -C "$SERVER_DIR" reset --hard origin/master -q
  else
    git clone --depth=1 --filter=blob:none --sparse -q \
      "$UPSTREAM_REMOTE/defillama-server.git" "$SERVER_DIR"
  fi
  # Re-applied every run so the cone is self-healing if a prior run left the default /* + !/*/ pattern.
  git -C "$SERVER_DIR" sparse-checkout set defi/src/protocols
) &
pids+=("$!")

(
  curl -fsSL \
    "https://github.com/DefiLlama/DefiLlama-Adapters/releases/download/latest/tvlModules.json" \
    -o "$TVL_MODULES.tmp"
  mv "$TVL_MODULES.tmp" "$TVL_MODULES"
) &
pids+=("$!")

(
  curl -fsSL \
    "https://github.com/DefiLlama/dimension-adapters/releases/download/latest/dimensionModules.json" \
    -o "$DIM_MODULES.tmp"
  mv "$DIM_MODULES.tmp" "$DIM_MODULES"
) &
pids+=("$!")

# Bare `wait` only returns the last job's status. Wait on each pid so any
# failure across the three parallel fetches aborts the script.
fail=0
for pid in "${pids[@]}"; do
  if ! wait "$pid"; then fail=1; fi
done
[ "$fail" -eq 0 ] || { echo "refresh.sh: one or more upstream fetches failed" >&2; exit 1; }
SERVER_SHA=$(git -C "$SERVER_DIR" rev-parse --short HEAD 2>/dev/null || echo "?")
TVL_KIB=$(( $(wc -c < "$TVL_MODULES") / 1024 ))
DIM_KIB=$(( $(wc -c < "$DIM_MODULES") / 1024 ))
done_ "server@$SERVER_SHA  tvlModules ${TVL_KIB} KiB  dimensionModules ${DIM_KIB} KiB  ($((SECONDS - T_FETCH))s)"

# 3. Run the bun extractor over defillama-server/data*.ts to produce $PROTOCOLS_JSON.
section "Extract protocol catalog"
step "bun tools/extract-protocols.ts > $PROTOCOLS_JSON"
mkdir -p "$(dirname "$PROTOCOLS_JSON")"
T_EXTRACT=$SECONDS
bun run tools/extract-protocols.ts > "$PROTOCOLS_JSON"
PROTOCOL_COUNT=$(bun -e "const d = await Bun.file('$PROTOCOLS_JSON').json(); console.log(Object.values(d).reduce((s, a) => s + a.length, 0))" 2>/dev/null || echo "?")
done_ "$PROTOCOL_COUNT protocols extracted  ($((SECONDS - T_EXTRACT))s)"

# 4. Build snapshot.json from the catalog and the two manifests, then sync into matrix and protocol_identities.
section "Build snapshot"
step "bun tools/build-snapshot.ts > $SNAPSHOT_OUT"
T_BUILD=$SECONDS
bun run tools/build-snapshot.ts
SNAPSHOT_PROTOS=$(bun -e "const s = await Bun.file('$SNAPSHOT_OUT').json(); console.log(s.protocols.length)" 2>/dev/null || echo "?")
SNAPSHOT_CELLS=$(bun -e "const s = await Bun.file('$SNAPSHOT_OUT').json(); console.log(s.cells.length)" 2>/dev/null || echo "?")
SNAPSHOT_KIB=$(( $(wc -c < "$SNAPSHOT_OUT") / 1024 ))
done_ "$SNAPSHOT_PROTOS protocols, $SNAPSHOT_CELLS cells, ${SNAPSHOT_KIB} KiB  ($((SECONDS - T_BUILD))s)"

# Sync the snapshot into Postgres when a DATABASE_URL is configured. The CI
# refresh job leaves this unset so the workflow stays bash + bun, no Go.
if [ -n "${DATABASE_URL:-}" ]; then
  section "Sync to Postgres"
  if [ ! -x bin/sync-db ] || [ -n "$(find cmd/sync-db internal/config internal/db internal/logger internal/models -newer bin/sync-db -type f 2>/dev/null)" ]; then
    step "go build -o bin/sync-db ./cmd/sync-db"
    go build -o bin/sync-db ./cmd/sync-db
  fi
  step "bin/sync-db"
  T_SYNC=$SECONDS
  bin/sync-db
  done_ "matrix + protocol_identities reloaded  ($((SECONDS - T_SYNC))s)"
else
  section "Sync to Postgres"
  warn "DATABASE_URL unset; skipping sync-db (CI mode)"
fi

# 5. Build the frontend in a sibling staging dir, then atomic-swap web/out/.
#    Building in place would clear web/out/ early in next build's export phase
#    and the Go static handler would serve 404s until the new tree finished
#    writing. The sibling pattern keeps web/out/ untouched until the final mv.
#    rsync skips node_modules / .next / out for speed; node_modules is reused
#    via symlink so bun does not have to reinstall.
if [ -z "${SKIP_WEB_BUILD:-}" ] && [ -f web/package.json ]; then
  section "Build web (atomic swap)"
  step "stage web/ -> ./web.stage"
  WEB_STAGE="$PWD/web.stage"
  # Clear stragglers from a prior failed/killed run
  rm -rf "$WEB_STAGE" web/out.old
  rsync -aq --delete \
    --exclude=node_modules --exclude=.next --exclude=out \
    web/ "$WEB_STAGE/"
  ln -snf "$PWD/web/node_modules" "$WEB_STAGE/node_modules"
  step "bun run build"
  T_WEB=$SECONDS
  ( cd "$WEB_STAGE" && bun run build )
  [ -d web/out ] && mv web/out web/out.old
  mv "$WEB_STAGE/out" web/out
  rm -rf "$WEB_STAGE"
  rm -rf web/out.old
  WEB_KIB=$(du -sk web/out 2>/dev/null | awk '{print $1}')
  done_ "web/out/ atomic-swapped (${WEB_KIB:-?} KiB)  ($((SECONDS - T_WEB))s)"
else
  section "Build web"
  warn "SKIP_WEB_BUILD set or web/package.json missing; skipping"
fi

printf '\n%s▸ refresh complete in %ds%s\n\n' "$c_green" "$((SECONDS - T_START))" "$c_reset"
