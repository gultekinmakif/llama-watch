# llama-watch

[![Go Report Card](https://goreportcard.com/badge/github.com/gultekinmakif/llama-watch)](https://goreportcard.com/report/github.com/gultekinmakif/llama-watch)

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

See docs at [gultekinmakif.github.io/llama-watch](https://gultekinmakif.github.io/llama-watch/).


## Configuration

| Var | Default | Notes |
|---|---|---|
| `PORT` | `3000` | |
| `ENV` | `dev` | `dev` → tint pretty logs, `prod` → JSON |
| `LOG_LEVEL` | `debug` | `debug` / `info` / `warn` / `error` |
| `SHUTDOWN_TIMEOUT` | `10s` | Graceful drain timeout |
| `DATABASE_URL` | *(required)* | Postgres DSN, e.g. `postgres://postgres:postgres@localhost:5432/llama_watch?sslmode=disable` |

## License

see [LICENSE](LICENSE).
