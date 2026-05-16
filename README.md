# llama-watch

DefiLlama coverage matrix. 

A Go cron parses [DefiLlama/DefiLlama-Adapters](https://github.com/DefiLlama/DefiLlama-Adapters), [DefiLlama/dimension-adapters](https://github.com/DefiLlama/dimension-adapters), and [DefiLlama/defillama-server](https://github.com/DefiLlama/defillama-server) once an hour, persists which adapter files cover which protocol × dimension pairs, and serves a sortable, filterable matrix.

## Stack

- **Backend.** Go 1.26 + stdlib `net/http` (1.22+ ServeMux). GORM v1.31 against Postgres. `log/slog` via `lmittmann/tint` for dev colored logs, JSON for prod.
- **Refresh pipeline:** 
  - Bash orchestrator (`scripts/refresh.sh`)
  - bun extractor (`tools/extract-protocols.ts` reads `defillama-server`'s `data{1..6}.ts`) 
  - Go `bin/refresh`:
    - walks the dir trees of both adapter clones
    - upserts protocols + adapter_files
    - writes a JSON snapshot
- **Frontend.** Next.js 16 + React 19 + Tailwind 4 + `@tanstack/react-table` + `@tanstack/react-virtual`.
  - **Static export into `web/out/`**; the Go server's file root serves it directly.
- **Scheduler.** launchd plist (macOS), systemd `.service` + `.timer` (Linux), or a crontab line.
  - Templates in `scripts/launchd/` and `scripts/systemd/`.

## Quick start

```sh
# 1. Postgres (any path)
docker run -d --name llama-watch-db -p 5432:5432 \
  -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=llama_watch \
  postgres:17
```

### List all the commands:
```sh

make help 

```

### Server (auto-migrates on boot)
```sh
make run                   
```
### or for hot reload via air:
```sh
make dev
```

### Refresh (clones the three upstream repos, runs bun + bin/refresh)
```sh
make refresh
```

## Endpoints

| Method | Path | Returns |
|---|---|---|
| `GET` | `/health` | `{"status":"ok","db":"ok","db_ms":12}` on 200, or `{"status":"down","db":"down","db_ms":2000}` on 503 if the DB ping fails (2s timeout). `db_ms` is the actual ping latency, useful for plotting degradation before it becomes downtime. |
| `GET` | `/api/matrix` | `{ columns, rows, total }` paginated via `?limit` (default 200, max 1000) and `?offset`. |
| `GET` | `/api/matrix/{slug}` | Detail view for one matrix row: `{ slug, name, category, chains, methodology, dimensions[] }`. Each `dimensions[]` entry covers one of the 9 pinned kinds with `present` + `file_path` + `repo` + `last_commit`. 404 envelope if the slug is unknown. |
| `GET` | `/*` | Static files from `web/out/` (the prerendered Next export). 404 from the file system on miss.|

See `specs/ENDPOINTS.md` for full response shapes. `methodology` and `last_commit` are intentionally empty / null in v1 — the refresh pipeline doesn't surface those yet.

## Refresh / cron

`scripts/refresh.sh` is the orchestrator. Each tick:

- `git pull`s the three upstream repos into `var/upstream/` (clones on first run).
- Runs the bun extractor: `tools/extract-protocols.ts` → `var/extracted/protocols.json`.
- Runs `bin/refresh`, rebuilds if its sources changed. The `--interval` flag (default 3300s) skips when the previous run finished too recently.
- If `web/package.json` exists, runs `pnpm build` and atomic-swaps `web/out/`. Skipped silently otherwise.

### Pick one scheduler

1. **launchd (macOS).** Copy `scripts/launchd/com.gultekinmakif.llama-watch.refresh.plist` into `~/Library/LaunchAgents/`, replace the `/ABSOLUTE/PATH/TO/llama-watch` placeholders, then:
   ```sh
   launchctl load -w ~/Library/LaunchAgents/com.gultekinmakif.llama-watch.refresh.plist
   ```
2. **systemd (Linux).** Copy `scripts/systemd/llama-watch-refresh.{service,timer}` into `/etc/systemd/system/`, replace the placeholders, then:
   ```sh
   systemctl enable --now llama-watch-refresh.timer
   ```
3. **crontab (any Unix).** Append:
   ```
   0 * * * * cd /path/to/llama-watch && ./scripts/refresh.sh >> var/log/refresh.log 2>&1
   ```
   Set `MAILTO=...` at the top of the crontab to catch non-zero exits.

Logs land in `var/log/refresh.{out,err}` (or `refresh.log` for the crontab path).

## Configuration

Server (`bin/server`) reads env vars at boot. Refresh script (`scripts/refresh.sh`) reads `.env` if present, else falls back to the defaults below.

| Var | Used by | Default | Notes |
|---|---|---|---|
| `DATABASE_URL` | server, refresh | `postgres://postgres:postgres@localhost:5432/llama_watch?sslmode=disable` | Override for non-default DSN. |
| `PORT` | server | `3000` | |
| `ENV` | server | `dev` | `dev` → tint pretty logs, `prod` → JSON. |
| `LOG_LEVEL` | server | `debug` | `debug` / `info` / `warn` / `error`. |
| `SHUTDOWN_TIMEOUT` | server | `10s` | Graceful drain on SIGTERM. |
| `INTERVAL` | refresh | `3300` | Skip the run if `now - last_finished < INTERVAL` seconds. |
| `UPSTREAM_DIR` | refresh | `./var/upstream` | Where the three repo clones live. |
| `PROTOCOLS_JSON` | refresh | `./var/extracted/protocols.json` | Bun extractor output the Go binary reads. |
| `SNAPSHOT_OUT` | refresh | `./var/snapshot/snapshot.json` | Atomic rename target. |
| `REPOS` | refresh | `DefiLlama-Adapters dimension-adapters defillama-server` | Whitespace-separated; cloned from `https://github.com/DefiLlama/<name>.git`. |
| `TEST_DATABASE_URL` | tests | *(optional)* | When set, `internal/api` tests run against this Postgres; when unset, the api test package skips cleanly. |

## Layout

```
cmd/
  server/    long-running HTTP (cmd/server/main.go)
  refresh/   one-shot cron orchestrator (cmd/refresh/main.go)
internal/
  api/       /api handlers + DTOs + SQL helpers
  config/    env loader
  db/postgres/   GORM handle + AutoMigrate + dimensions seed
  dimensions/    walker + protocols loader + slug-join builder + sub-metric regex
  handlers/  /health
  logger/    slog factory (tint dev, JSON prod)
  middleware/    recoverer, request_id, request logger
  models/    GORM table definitions
  server/    HTTP server bootstrap + route mount
  snapshot/  atomic JSON writer primitive
specs/       SCHEMA.md, PARSER.md, ENDPOINTS.md, REFRESH.md (the contract for the rebuild)
scripts/     refresh.sh + dev.sh + launchd/systemd templates
tools/       bun extractor (not committed yet)
var/         gitignored: upstream clones, extracted JSON, snapshot, logs
web/         frontend (Phase 2)
```

## Tests

```sh
make test                              # unit tests only; api/* tests self-skip without TEST_DATABASE_URL
TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/llama_watch_test?sslmode=disable make test
```

The api package's tests open the singleton via `postgres.New`, run `Migrate`, and isolate each test inside a transaction with `Rollback` at the end.

## License

See [LICENSE](LICENSE).
