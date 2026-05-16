#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# 1. Load env (DATABASE_URL, INTERVAL, UPSTREAM_DIR, PROTOCOLS_JSON, SNAPSHOT_OUT, REPOS).
[ -f .env ] && source .env

: "${REPOS:=DefiLlama-Adapters dimension-adapters defillama-server}"
: "${INTERVAL:=3300}"
: "${UPSTREAM_DIR:=./var/upstream}"
: "${PROTOCOLS_JSON:=./var/extracted/protocols.json}"
: "${SNAPSHOT_OUT:=./var/snapshot/snapshot.json}"

# 2. For each repo in REPOS, git pull (clone if missing) into $UPSTREAM_DIR/.
mkdir -p "$UPSTREAM_DIR"
for repo in $REPOS; do
  if [ -d "$UPSTREAM_DIR/$repo/.git" ]; then
    git -C "$UPSTREAM_DIR/$repo" pull --ff-only
  else
    git clone "https://github.com/DefiLlama/$repo.git" "$UPSTREAM_DIR/$repo"
  fi
done

# 3. Run the bun extractor over defillama-server/data*.ts to produce $PROTOCOLS_JSON.
mkdir -p "$(dirname "$PROTOCOLS_JSON")"
bun run tools/extract-protocols.ts > "$PROTOCOLS_JSON"

# 4. Ensure bin/refresh exists and is newer than its sources, then rebuild if stale.
if [ ! -x bin/refresh ] || [ -n "$(find cmd/refresh internal/dimensions internal/snapshot internal/models -newer bin/refresh -type f 2>/dev/null)" ]; then
  go build -o bin/refresh ./cmd/refresh
fi

# 5. Run the Go refresh binary; the --interval gate inside skips redundant runs.
bin/refresh \
  --interval="$INTERVAL" \
  --upstream-dir="$UPSTREAM_DIR" \
  --protocols-json="$PROTOCOLS_JSON" \
  --snapshot-out="$SNAPSHOT_OUT"

# 6. Build the frontend and atomic-swap web/out/. Skip when web/ is unscaffolded.
if [ -f web/package.json ]; then
  ( cd web && bun run build )
  # Phase 2 may emit directly to web/out/; revisit this swap when the real build path lands.
  rm -rf web/out.old
  [ -d web/out ] && mv web/out web/out.old
  mv web/out.new web/out
else
  echo "web/package.json missing; skipping frontend build"
fi
