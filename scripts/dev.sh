#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

GOBIN=$(go env GOBIN)
[ -z "$GOBIN" ] && GOBIN=$(go env GOPATH)/bin
export PATH="$GOBIN:$PATH"

if ! command -v air >/dev/null; then
  echo "installing air..."
  go install github.com/air-verse/air@latest
fi

exec air
