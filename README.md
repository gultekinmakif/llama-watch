# go-http-server

[![Go Report Card](https://goreportcard.com/badge/github.com/gultekinmakif/go-http-server)](https://goreportcard.com/report/github.com/gultekinmakif/go-http-server)
[![Template](https://img.shields.io/badge/template-use%20this-2ea44f)](https://github.com/gultekinmakif/go-http-server/generate)

> **This is a GitHub template.** Click **"Use this template"** at the top of the repo (or the badge above) to spin up a fresh service. See [Use this template](#use-this-template) below for the one-command rename.

Minimal Go HTTP server template. See docs at [gultekinmakif.github.io/go-http-server](https://gultekinmakif.github.io/go-http-server/).

Standard library `net/http` with:

- Go 1.22+ method-prefix routing
- middleware: request ID, request logger, panic recoverer
- DB: Postgres + GORM
- Docker: multi-stage, distroless
- env: typed config, validated at startup
- logs: `slog` — tint in dev, JSON in prod
- graceful shutdown on SIGINT/SIGTERM

## Use this template

Create a new repo from this template (or click "Use this template" in the UI):
```sh
gh repo create my-go-server --template gultekinmakif/go-http-server --private --clone
cd my-go-server

# 2. Rename the module path + DB name + Docker tag in one shot
./scripts/init-template.sh github.com/<you>/my-go-server
go mod tidy
```

After that, `make docker-up` and you're running.

## Quick start

```sh
docker compose up --build        # app + postgres on :3000
curl localhost:3000/health
```

Or against your own Postgres:

```sh
cp configs/.env.example .env     
make run                         
make dev                         # hot-reload 
make help
```

## Endpoints

| Method | Path | Notes |
|---|---|---|
| `GET` | `/health` | Liveness -> `{"status":"ok"}` |
| `POST` | `/posts` | Create. Body: `{title, body, slug}`. Returns 201 + record. |
| `GET` | `/posts` | List all posts |
| `GET` | `/posts/{slug}` | Read by slug. 404 on miss. |

See docs at [gultekinmakif.github.io/go-http-server](https://gultekinmakif.github.io/go-http-server/).


## Configuration

| Var | Default | Notes |
|---|---|---|
| `PORT` | `3000` | |
| `ENV` | `dev` | `dev` → tint pretty logs, `prod` → JSON |
| `LOG_LEVEL` | `debug` | `debug` / `info` / `warn` / `error` |
| `SHUTDOWN_TIMEOUT` | `10s` | Graceful drain timeout |
| `DATABASE_URL` | *(required)* | Postgres DSN, e.g. `postgres://postgres:postgres@localhost:5432/http_server?sslmode=disable` |

## License

see [LICENSE](LICENSE).
