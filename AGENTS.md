# AGENTS.md

## Layout

- `client/` — React 19 + Vite 7 + TypeScript SPA (`npm --prefix client run dev`, port 8000)
- `server/` — Go 1.24 stdlib HTTP server (no framework), port 8080. Entrypoint `server/main.go`; handlers in `server/routes/v1/`
- `migrations/` — plain SQL, mounted into Postgres at `docker-entrypoint-initdb.d`. **Only runs on first DB volume creation** — schema changes to an existing volume require manual psql
- `bin/localstack/init-s3.sh` — creates the S3 bucket + CORS on LocalStack startup

## Running

```sh
docker compose up --build          # everything (server, client, db, s3)
npm --prefix client run dev        # vite dev server instead of containerized client
```

Env lives in `.env`: `DB_URL`, `DB_PASSWORD`, `LOCALSTACK_AUTH_TOKEN`.

## Gotchas

- Storage is LocalStack S3 hardcoded to `http://localhost:4566` in `server/storage/main.go`
- Server requires `DB_URL` env or exits at startup; LocalStack init also requires `LOCALSTACK_AUTH_TOKEN`.
- Client talks to `/api/v1/*` routes; CORS handled manually in `server/routes.go`.
