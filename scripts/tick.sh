#!/usr/bin/env bash
set -euo pipefail

# Dev loop: re-run refresh.sh every INTERVAL seconds.
# Ctrl-C to stop
# !! use setup-cron.sh on prod
cd "$(dirname "$0")/.."
INTERVAL="${INTERVAL:-3600}"
while true; do
  ./scripts/refresh.sh || true
  sleep "${INTERVAL}"
done
