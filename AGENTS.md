# AGENTS.md

## Layout

- `client/` — React 19 + Vite 7 + TypeScript SPA (`npm --prefix client run dev`, port 8000)
- `server/` — Go 1.24 stdlib HTTP server (no framework), port 8080. Entrypoint `server/main.go`; handlers in `server/routes/v1/`
- `migrations/` — plain SQL, mounted into Postgres at `docker-entrypoint-initdb.d`. **Only runs on first DB volume creation** — schema changes to an existing volume require manual psql
- `bin/localstack/init-awslocal.sh` — creates the S3 bucket + CORS, the `vidtube-uploads` SQS queue, and the bucket→queue notification on LocalStack startup. Config JSONs live alongside it (`cors.json`, `queue-policy.json`, `bucket-notifications.json`)

## Running

```sh
docker compose up --build          # everything (server, client, db, aws)
npm --prefix client run dev        # vite dev server instead of containerized client
```

Env lives in `.env`: `DB_URL`, `DB_PASSWORD`, `LOCALSTACK_AUTH_TOKEN`.

## Gotchas

- LocalStack (service name `aws`) serves S3 and SQS on one port, 4566. Compose sets both `S3_ENDPOINT` and `SQS_ENDPOINT` to `http://localhost:4566`; the server shares LocalStack's network namespace via `network_mode: "service:aws"`, so localhost is literal. Presigned URLs use `S3_ENDPOINT`. Both clients currently use dummy credentials for LocalStack; production requires credential and endpoint configuration changes
- The server publishes no ports itself; 8080 is mapped on the `aws` container since they share a network namespace. DB_URL must still resolve the `db` hostname (same compose network)
- Video readiness is driven by S3 events → SQS → the consumer goroutine in `server/uploadingest/main.go`; there is no client-facing "complete" endpoint
- Server requires `DB_URL` env or exits at startup; Compose sets `PORT=8080` explicitly and passes `LOCALSTACK_AUTH_TOKEN` to LocalStack. LocalStack runs `init-awslocal.sh` at the ready stage on startup; restart `aws` to rerun changed initialization configuration. Postgres migrations still only run on first DB volume creation
- Client talks to `/api/v1/*` routes; CORS handled manually in `server/routes.go`. Browser uploads to S3 rely on `bin/localstack/cors.json`
