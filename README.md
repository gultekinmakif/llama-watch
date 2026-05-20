# llama-watch
> Generated from [gultekinmakif/go-http-server](https://github.com/gultekinmakif/go-http-server)
DefiLlama coverage matrix.

A daily refresh sparse-clones [DefiLlama/defillama-server](https://github.com/DefiLlama/defillama-server) for the protocol catalog, fetches the [DefiLlama-Adapters](https://github.com/DefiLlama/DefiLlama-Adapters) and [dimension-adapters](https://github.com/DefiLlama/dimension-adapters) release-asset manifests, persists one row per (protocol, metric) coverage cell, and serves the result as a sortable, filterable matrix.

## Configuration

| Var | Used by | Default | Notes |
|---|---|---|---|
| `UPSTREAM_REMOTE` | refresh | `https://github.com/DefiLlama` | Owner URL prefix for the defillama-server sparse clone; override to point at a fork. |
| `UPSTREAM_DIR` | refresh | `./var/upstream` | Where the defillama-server sparse clone lives. |
| `PROTOCOLS_JSON` | refresh | `./var/snapshot/protocols.json` | Normalized catalog emitted by `tools/extract-protocols.ts` and consumed by `tools/build-snapshot.ts`. |
| `SNAPSHOT_OUT` | refresh | `./var/snapshot/snapshot.json` | Atomic rename target produced by `tools/build-snapshot.ts` and consumed by `bin/sync-db`. |
| `DATABASE_URL` | server, sync-db | `postgres://postgres:postgres@localhost:5432/llama_watch?sslmode=disable` | Override for non-default DSN. |
| `SHUTDOWN_TIMEOUT` | server | `10s` | Graceful drain on SIGTERM. |
| `LOG_LEVEL` | server | `debug` | `debug` / `info` / `warn` / `error`. |
| `ENV` | server | `dev` | `dev` enables tint pretty logs; `prod` emits JSON. |
| `PORT` | server | `3000` | |
| `TEST_DATABASE_URL` | tests | *(optional)* | When set, `internal/api` tests run against this Postgres; when unset, the api test package skips cleanly. |

## Quick start

Prereqs: `bun`, `go` 1.26+, `postgres` (or Docker), `rsync` (refresh.sh uses it for the atomic web swap; present on macOS and stock Linux).

### 1. Postgres

```sh
docker run -d --name llama-watch-db -p 5432:5432 \
  -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=llama_watch \
  postgres:17
```

### 2. Server

> Use `make dev` for hot reload with air.

```sh
make run
```

### 3. Refresh

Sparse-clones `defillama-server`, fetches the two upstream release manifests, runs the bun pipeline (`tools/extract-protocols.ts` then `tools/build-snapshot.ts`), and bulk-loads the result via `bin/sync-db`.

```sh
make refresh
```

### 4. Tests

```sh
make test
# Integration tests require Postgres on TEST_DATABASE_URL:
TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/llama_watch_test?sslmode=disable make test
```

The api package opens the singleton via `postgres.New`, runs `Migrate`, and rolls back each test's transaction at the end.

> Run `make help` for the full command list.

## Cron Job

`scripts/refresh.sh` is the orchestrator, scheduled daily. Each tick, in parallel:

- Sparse-clones (or fetches) `defillama-server` into `var/upstream/defillama-server/`, scoped to `defi/src/protocols`.
- Curls `tvlModules.json` from the `DefiLlama-Adapters` `latest` release into `var/snapshot/tvlModules.json`.
- Curls `dimensionModules.json` from the `dimension-adapters` `latest` release into `var/snapshot/dimensionModules.json`.

After the upstream jobs join, sequentially:

- Runs `tools/extract-protocols.ts` to normalize the catalog into `var/snapshot/protocols.json`.
- Runs `tools/build-snapshot.ts` to join the catalog against both manifests into `var/snapshot/snapshot.json`.
- Builds (if stale) and runs `bin/sync-db`, which truncates `matrix` + `protocol_identities` and bulk-inserts from the snapshot in one transaction.
- If `web/package.json` exists, runs `bun run build` and atomic-swaps `web/out/`. Skipped silently otherwise.

### Install the scheduler

```sh
chmod +x scripts/setup-cron.sh    # first time only
./scripts/setup-cron.sh
```

The script detects your OS and installs the matching scheduler: launchd on macOS, systemd on Linux (requires sudo), or a crontab line as fallback. Logs land in `var/log/refresh.{out,err}` (or `refresh.log` for crontab).

## Endpoints

| Method | Path | Returns |
|---|---|---|
| `GET` | `/health` | `{"status":"ok","db":"ok","db_ms":12}` on 200, or `{"status":"down","db":"down","db_ms":2000}` on 503 if the DB ping fails (2s timeout). `db_ms` is the actual ping latency, useful for plotting degradation before it becomes downtime. |
| `GET` | `/api/matrix` | `{ columns, rows, total }` paginated via `?limit` (default 200, max 1000) and `?offset`. |
| `GET` | `/api/matrix/{slug}` | Detail view for one matrix row: `{ slug, name, category, chains, dimensions[] }`. Each `dimensions[]` entry carries `kind`, `present`, and a `github_url` pointing into the dimension-adapters repo when the protocol has an adapter for that metric. 404 envelope if the slug is unknown. |
| `GET` | `/api/chains` | `{ chains: [{ key, label, protocol_count }] }`. Distinct chains across all protocols, ordered by key; label is titlecase of key. |
| `GET` | `/api/dimensions` | `{ dimensions: [{ kind, display_name, coverage }] }`. Ordered by `internal/registry/columns.go`; `coverage` is the per-metric `COUNT(*)` over the matrix table. |
| `GET` | `/*` | Static files from `web/out/` (the prerendered Next export). 404 from the file system on miss. |

## Stack

- **Backend.** Go 1.26 + stdlib `net/http` (1.22+ ServeMux). GORM v1.31 against Postgres. `log/slog` via `lmittmann/tint` for dev colored logs, JSON for prod.
- **Refresh pipeline:**
  - Bash orchestrator (`scripts/refresh.sh`)
  - Sparse `git` checkout of `defillama-server` (`defi/src/protocols`) plus two parallel `curl`s for `tvlModules.json` and `dimensionModules.json`.
  - bun TS pipeline: `tools/extract-protocols.ts` normalizes the catalog, `tools/build-snapshot.ts` joins it against both manifests into `snapshot.json`.
  - Go `bin/sync-db`: opens one transaction, truncates `matrix` + `protocol_identities`, bulk-inserts in batches of 500.
- **Frontend.** Next.js 16 + React 19 + Tailwind 4 + `@tanstack/react-table` + `@tanstack/react-virtual`.
  - **Static export into `web/out/`**; the Go server's file root serves it directly.
- **Scheduler.** launchd plist (macOS), systemd `.service` + `.timer` (Linux), or a crontab line.
  - Templates in `scripts/launchd/` and `scripts/systemd/`.

## Frontend

The Next.js 16 static export lives under [`web/`](web/). See [`web/README.md`](web/README.md) for the frontend quick start.


## Development

CI runs `go build`, `go vet`, `go test`, `golangci-lint`, and `shellcheck --severity=warning scripts/*.sh` on every PR.

Run the same locally before pushing; install `shellcheck` via `brew install shellcheck` (or your package manager) if you don't have it.

### Releases

1. Cut a branch `release/vX.Y.Z` off `main`.
2. Bump `web/version.json` to `{ "version": "vX.Y.Z" }` in that branch.
3. Open and merge the PR to `main`.

Only merged PRs from `release/v*` branches trigger a release.
Direct pushes to `main` (including the daily snapshot refresh bot) **do not**.
<!-- On merge, `.github/workflows/release.yml` validates that `web/version.json` matches the branch name, then creates the matching git tag and a GitHub release with auto-generated notes. Vercel sees the push to `main` and rebuilds; the build inlines `web/version.json`, so the footer chip and its release-page link both render the new version automatically. -->

## License

See [LICENSE](LICENSE).
