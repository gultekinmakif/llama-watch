#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

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

pids=()

(
  if [ -d "$SERVER_DIR/.git" ]; then
    git -C "$SERVER_DIR" fetch --depth=1
    git -C "$SERVER_DIR" reset --hard origin/master
  else
    git clone --depth=1 --filter=blob:none --sparse \
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

# 3. Run the bun extractor over defillama-server/data*.ts to produce $PROTOCOLS_JSON.
mkdir -p "$(dirname "$PROTOCOLS_JSON")"
bun run tools/extract-protocols.ts > "$PROTOCOLS_JSON"

# 4. Build snapshot.json from the catalog and the two manifests, then sync into matrix and protocol_identities.
bun run tools/build-snapshot.ts

if [ ! -x bin/sync-db ] || [ -n "$(find cmd/sync-db internal/db internal/models -newer bin/sync-db -type f 2>/dev/null)" ]; then
  go build -o bin/sync-db ./cmd/sync-db
fi
bin/sync-db

# 5. Build the frontend in a sibling staging dir, then atomic-swap web/out/.
#    Building in place would clear web/out/ early in next build's export phase
#    and the Go static handler would serve 404s until the new tree finished
#    writing. The sibling pattern keeps web/out/ untouched until the final mv.
#    rsync skips node_modules / .next / out for speed; node_modules is reused
#    via symlink so bun does not have to reinstall.
if [ -f web/package.json ]; then
  WEB_STAGE="$PWD/web.stage"
  # Clear stragglers from a prior failed/killed run
  rm -rf "$WEB_STAGE" web/out.old
  rsync -a --delete \
    --exclude=node_modules --exclude=.next --exclude=out \
    web/ "$WEB_STAGE/"
  ln -snf "$PWD/web/node_modules" "$WEB_STAGE/node_modules"
  ( cd "$WEB_STAGE" && bun run build )
  [ -d web/out ] && mv web/out web/out.old
  mv "$WEB_STAGE/out" web/out
  rm -rf "$WEB_STAGE"
  rm -rf web/out.old
else
  echo "web/package.json missing; skipping frontend build"
fi
