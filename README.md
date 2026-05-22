# llama-watch
> Generated from [gultekinmakif/go-http-server](https://github.com/gultekinmakif/go-http-server)

DefiLlama coverage matrix. A daily refresh sparse-clones [defillama-server](https://github.com/DefiLlama/defillama-server) for the protocol catalog, fetches release manifests from [DefiLlama-Adapters](https://github.com/DefiLlama/DefiLlama-Adapters) and [dimension-adapters](https://github.com/DefiLlama/dimension-adapters), persists one row per `(protocol, metric)` cell, and serves the result as a sortable, filterable matrix.

## Configuration

| Var | Used by | Default | Notes |
|---|---|---|---|
| `UPSTREAM_REMOTE` | refresh | `https://github.com/DefiLlama` | Owner URL prefix; override to point at a fork. |
| `UPSTREAM_DIR` | refresh | `./var/upstream` | Where the sparse clone lives. |
| `PROTOCOLS_JSON` | refresh | `./var/snapshot/protocols.json` | Normalized catalog between `extract-protocols` and `build-snapshot`. |
| `SNAPSHOT_OUT` | refresh | `./var/snapshot/snapshot.json` | Atomic rename target consumed by `bin/sync-db`. |
| `DATABASE_URL` | server, sync-db | `postgres://postgres:postgres@localhost:5432/llama_watch?sslmode=disable` | Override for non-default DSN. |
| `SHUTDOWN_TIMEOUT` | server | `10s` | Graceful drain on SIGTERM. |
| `LOG_LEVEL` | server | `debug` | `debug` / `info` / `warn` / `error`. |
| `ENV` | server | `dev` | `dev` uses tint pretty logs; `prod` emits JSON. |
| `PORT` | server | `3000` | |
| `TEST_DATABASE_URL` | tests | *(optional)* | When set, integration tests run against it; otherwise the api test package skips. |

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

Runs the full pipeline once: upstream fetches, snapshot build, DB sync, and frontend stage. See [Cron Job](#cron-job) for the per-phase breakdown.

```sh
make refresh
```

### 4. Tests

```sh
make test
# Integration tests require Postgres on TEST_DATABASE_URL:
TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/llama_watch_test?sslmode=disable make test
```

Integration tests migrate the DB once and roll back each test's transaction.

> Run `make help` for the full command list.

## Cron Job

`scripts/refresh.sh` orchestrates the daily pipeline. In parallel:

- Sparse-clones (or fetches) `defillama-server` scoped to `defi/src/protocols`.
- Curls `tvlModules.json` from the `DefiLlama-Adapters` `latest` release.
- Curls `dimensionModules.json` from the `dimension-adapters` `latest` release.

Once they join, sequentially:

- `tools/extract-protocols.ts` normalizes the catalog into `protocols.json`.
- `tools/build-snapshot.ts` joins the catalog against both manifests into `snapshot.json`.
- `bin/sync-db` truncates `matrix` + `protocol_identities` and bulk-inserts in one transaction.
- If `web/package.json` exists, `bun run build` runs and atomic-swaps `web/out/`.

### Install the scheduler

```sh
chmod +x scripts/setup-cron.sh    # first time only
./scripts/setup-cron.sh
```

The script detects your OS and installs the matching scheduler: launchd on macOS, systemd on Linux (requires sudo), or a crontab line as fallback. Logs land in `var/log/refresh.{out,err}` (or `refresh.log` for crontab).

## Endpoints

| Method | Path | Returns |
|---|---|---|
| `GET` | `/health` | `{ status, db, db_ms }`. 200 on a clean Postgres ping, 503 on a 2s timeout. `db_ms` is the actual ping latency. |
| `GET` | `/api/matrix` | `{ columns, rows, total }`, paginated via `?limit` (default 200, max 1000) and `?offset`. |
| `GET` | `/api/matrix/{slug}` | `{ slug, name, category, chains, dimensions[] }`. Each `dimensions[]` entry has `kind`, `present`, and a `github_url` into dimension-adapters when an adapter exists. 404 envelope on unknown slug. |
| `GET` | `/api/chains` | `{ chains: [{ key, label, protocol_count }] }`. Distinct chains, ordered by key. |
| `GET` | `/api/dimensions` | `{ dimensions: [{ kind, display_name, coverage }] }`. Ordered by the registry; `coverage` is the per-metric `COUNT(*)` over `matrix`. |
| `GET` | `/*` | Static files from `web/out/`. Filesystem 404 on miss. |

## Stack

- **Backend.** Go 1.26 + stdlib `net/http` (1.22+ ServeMux). GORM v1.31 against Postgres. `log/slog` via `lmittmann/tint` for dev, JSON for prod.
- **Refresh pipeline.** Bash orchestrator, sparse `git` clone of `defillama-server`, two parallel `curl`s for the upstream manifests, a bun TS pipeline that builds `snapshot.json`, and a Go `bin/sync-db` that bulk-loads in one transaction.
- **Frontend.** Next.js 16 + React 19 + Tailwind 4 + `@tanstack/react-table` + `@tanstack/react-virtual`. Static export into `web/out/`; the Go server's file root serves it directly.
- **Scheduler.** launchd plist (macOS), systemd `.service` + `.timer` (Linux), or a crontab line. Templates in `scripts/launchd/` and `scripts/systemd/`.

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
