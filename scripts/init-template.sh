#!/usr/bin/env bash
# init-template.sh — rename this template's module path and references after
# instantiating from gultekinmakif/go-http-server. Run ONCE in a fresh repo.
#
# Usage:
#   ./scripts/init-template.sh <new-module-path>
#
# Example:
#   ./scripts/init-template.sh github.com/akif/my-new-thing
#
# What it does:
#   - sed replace "github.com/gultekinmakif/go-http-server" → <new-module-path>
#   - sed replace "go-http-server"                          → basename of module
#   - sed replace "http_server"                             → basename with - → _
#
# Files touched: go.mod, all *.go, Makefile, docker-compose.yml,
#                configs/.env.example, api/openapi.yaml, README.md
#
# NOT touched: LICENSE (edit copyright manually), plan.md, scripts/init-template.sh.

set -euo pipefail

cd "$(dirname "$0")/.."

NEW_MODULE="${1:-}"
if [[ -z "$NEW_MODULE" ]]; then
  printf 'usage: ./scripts/init-template.sh <new-module-path>\n' >&2
  printf 'example: ./scripts/init-template.sh github.com/akif/my-new-thing\n' >&2
  exit 1
fi

# Loose validation: at least one slash, no URL scheme, safe charset.
if [[ "$NEW_MODULE" == http* ]]; then
  printf 'error: drop the scheme — pass "github.com/owner/repo", not "https://..."\n' >&2
  exit 1
fi
if [[ ! "$NEW_MODULE" =~ ^[a-zA-Z0-9._/-]+/[a-zA-Z0-9._-]+$ ]]; then
  printf 'error: module path must look like "host.tld/owner/repo", got: %s\n' "$NEW_MODULE" >&2
  exit 1
fi

OLD_MODULE="github.com/gultekinmakif/go-http-server"
OLD_NAME="go-http-server"
OLD_DB="http_server"

NEW_NAME="${NEW_MODULE##*/}"           # basename of the path
NEW_DB="${NEW_NAME//-/_}"              # hyphens → underscores for the DB name

if [[ "$NEW_MODULE" == "$OLD_MODULE" ]]; then
  printf 'error: nothing to do — module is already %s\n' "$OLD_MODULE" >&2
  exit 1
fi

printf '  old module: %s\n' "$OLD_MODULE"
printf '  new module: %s\n' "$NEW_MODULE"
printf '  new name:   %s\n' "$NEW_NAME"
printf '  new db:     %s\n' "$NEW_DB"
printf '\n'

# Files to update (relative to repo root).
TARGETS=(
  go.mod
  Makefile
  docker-compose.yml
  configs/.env.example
  api/openapi.yaml
  README.md
)
# Plus every .go file.
while IFS= read -r f; do TARGETS+=("$f"); done < <(find . -name '*.go' \
  -not -path './tmp/*' -not -path './bin/*' -not -path './vendor/*' -not -path './.git/*')

for f in "${TARGETS[@]}"; do
  [[ -f "$f" ]] || continue
  # macOS sed needs the suffix arg with -i; create a .bak then drop it.
  sed -i.bak \
    -e "s|${OLD_MODULE}|${NEW_MODULE}|g" \
    -e "s|${OLD_NAME}|${NEW_NAME}|g" \
    -e "s|${OLD_DB}|${NEW_DB}|g" \
    "$f"
  rm -f "${f}.bak"
  printf '  updated %s\n' "$f"
done

printf '\n==> renamed to %s\n\n' "$NEW_MODULE"
printf 'Next steps:\n'
printf '  1. go mod tidy\n'
printf '  2. update LICENSE copyright line manually\n'
printf '  3. (optional) rm scripts/init-template.sh — it has done its job\n'
