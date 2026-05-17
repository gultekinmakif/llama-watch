#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# 1. Load env (DATABASE_URL, INTERVAL, UPSTREAM_DIR, PROTOCOLS_JSON, SNAPSHOT_OUT).
[ -f .env ] && source .env

: "${UPSTREAM_REMOTE:=https://github.com/DefiLlama}"
: "${INTERVAL:=3300}"
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

(
  if [ -d "$SERVER_DIR/.git" ]; then
    git -C "$SERVER_DIR" fetch --depth=1
    git -C "$SERVER_DIR" reset --hard origin/master
  else
    git clone --depth=1 --filter=blob:none --sparse \
      "$UPSTREAM_REMOTE/defillama-server.git" "$SERVER_DIR"
    git -C "$SERVER_DIR" sparse-checkout set \
      defi/src/protocols defi/src/constants.ts
  fi
) &

(
  curl -fsSL \
    "https://github.com/DefiLlama/DefiLlama-Adapters/releases/download/latest/tvlModules.json" \
    -o "$TVL_MODULES.tmp"
  mv "$TVL_MODULES.tmp" "$TVL_MODULES"
) &

(
  curl -fsSL \
    "https://github.com/DefiLlama/dimension-adapters/releases/download/latest/dimensionModules.json" \
    -o "$DIM_MODULES.tmp"
  mv "$DIM_MODULES.tmp" "$DIM_MODULES"
) &

wait

# 3. Run the bun extractor over defillama-server/data*.ts to produce $PROTOCOLS_JSON.
mkdir -p "$(dirname "$PROTOCOLS_JSON")"
bun run tools/extract-protocols.ts > "$PROTOCOLS_JSON"

# 4. Build snapshot.json from the catalog and the two manifests, then sync into matrix and protocol_identities.
bun run tools/build-snapshot.ts

if [ ! -x bin/sync-db ] || [ -n "$(find cmd/sync-db internal/db internal/models -newer bin/sync-db -type f 2>/dev/null)" ]; then
  go build -o bin/sync-db ./cmd/sync-db
fi
bin/sync-db

# 5. Ensure bin/refresh exists and is newer than its sources, then rebuild if stale.
if [ ! -x bin/refresh ] || [ -n "$(find cmd/refresh internal/dimensions internal/snapshot internal/models -newer bin/refresh -type f 2>/dev/null)" ]; then
  go build -o bin/refresh ./cmd/refresh
fi

# 6. Run the Go refresh binary; the --interval gate inside skips redundant runs.
bin/refresh \
  --interval="$INTERVAL" \
  --upstream-dir="$UPSTREAM_DIR" \
  --protocols-json="$PROTOCOLS_JSON" \
  --snapshot-out="$SNAPSHOT_OUT"

# 7. Build the frontend and atomic-swap web/out/. Skip when web/ is unscaffolded.
if [ -f web/package.json ]; then
  ( cd web && bun run build )
  # Phase 2 may emit directly to web/out/; revisit this swap when the real build path lands.
  rm -rf web/out.old
  [ -d web/out ] && mv web/out web/out.old
  mv web/out.new web/out
else
  echo "web/package.json missing; skipping frontend build"
fi
