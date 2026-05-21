# final_project backend

This repository contains the Go HTTP server that exposes routing, geocoding, agent, and traffic-related APIs used by the final_project system. The server is mounted under `/api/v1/` and uses Supabase access tokens for authentication.

## Quick Start

- Copy the example environment file and fill values:

```bash
cp .env.example .env
# edit .env and set DB_URL, PORT, SUPABASE_*, ROUTING_SERVICE_ADDR, etc.
```

- Run locally (recommended for development):

```bash
# requires `air` for hot reload (optional)

```

- Run with Docker:

```bash
docker build -t final-project-backend-server:latest .
docker run -p 3000:3000 --env-file .env final-project-backend-server:latest
```

## Environment Variables

At minimum set these in your `.env` or environment:

- `DB_URL` — Postgres connection string
- `PORT` — Port the HTTP server listens on (default 3000)
- `ENV` — `dev` or `prod`
- `ROUTING_SERVICE_ADDR` — address of the routing gRPC service
- `GRPC_REQUEST_TIMEOUT` — gRPC deadline (e.g. `10s`)
- Supabase related:
  - `SUPABASE_URL`
  - `SUPABASE_SECRET_KEY` or `SUPABASE_SERVICE_ROLE_KEY`
  - `SUPABASE_PUBLISHABLE_KEY` (used in client flows)

See `internal/utils/load_env.go` for full loading behavior and defaults.

## Authentication

- `/health` is public.
- All endpoints under `/api/v1/` require a valid Supabase `Authorization: Bearer <token>` header unless `ENV=dev` and `DISABLE_AUTH=true` are set (see `internal/utils/load_env.go`).
- Some maintenance endpoints require an admin claim and are protected by `requireAdmin`.

## API Reference (v1)

Base URL (local): `http://localhost:3000/api/v1/`

Common headers:

- `Authorization: Bearer <access_token>` — required for protected endpoints
- `Content-Type: application/json` — for JSON POST bodies

Endpoints (high level):

- `GET /health` — Health check (public)

- Routing:
  - `POST /api/v1/route` — Calculate best route between points.
  - `POST /api/v1/routing/reload-prefix-times` — Reload prefix times (admin only).
  - `POST /api/v1/routing/rebuild-network` — Rebuild routing network (admin only).

- Agent:
  - `POST /api/v1/agent/query` — Send a message to the agent service and return a reply.

- Geocoding:
  - `POST /api/v1/geocode` — Convert an address string into lat/lon coordinates.

- DB tools:
  - `POST /api/v1/nearby-trips` — Query trips near a point within a radius.

- Traffic:
  - `POST /api/v1/traffic/trigger` — Trigger a traffic update job.
  - `GET /api/v1/traffic/status` — Get current traffic update status.
  - `POST /api/v1/traffic/update-trip` — Submit a trip for traffic modelling.
  - `POST /api/v1/traffic/street` — Compute traffic load for a particular street.
  - `GET /api/v1/traffic/streets` — List tracked streets and traffic levels.

Note: The server mounts the v1 router under `/api/v1/` (see [internal/server/server.go](internal/server/server.go)).

### Example Requests

- Health check (public):

```bash
curl -i http://localhost:3000/health
```

- Route calculation (example payload):

```bash
TOKEN="<ACCESS_TOKEN>"
curl -sS -X POST http://localhost:3000/api/v1/route \
   -H "Authorization: Bearer ${TOKEN}" \
   -H "Content-Type: application/json" \
   --data-raw '{"from":{"lat":40.7128,"lon":-74.0060},"to":{"lat":40.73061,"lon":-73.935242}}'
```

- Agent query:

```bash
curl -sS -X POST http://localhost:3000/api/v1/agent/query \
   -H "Authorization: Bearer ${TOKEN}" \
   -H "Content-Type: application/json" \
   --data-raw '{"message":"Find me the fastest route to the airport"}'
```

- Geocode example:

```bash
curl -sS -X POST http://localhost:3000/api/v1/geocode \
   -H "Authorization: Bearer ${TOKEN}" \
   -H "Content-Type: application/json" \
   --data-raw '{"address":"1600 Amphitheatre Parkway, Mountain View, CA"}'
```

- Traffic status (read-only):

```bash
curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/api/v1/traffic/status
```

### Admin endpoints

Requires an admin claim in the user token. Use these carefully — they affect the routing graph and background jobs.

- `POST /api/v1/routing/reload-prefix-times`
- `POST /api/v1/routing/rebuild-network`

## Swagger / API docs

In development (`ENV=dev`) the server serves a Swagger UI at `/docs/` (see `internal/server/server.go`).

## Testing & Development tips

- To exercise authenticated endpoints locally, obtain a Supabase access token (password grant) and set `Authorization: Bearer <token>`.
- If you need to disable auth for local dev, set `ENV=dev` and `DISABLE_AUTH=true` (see `internal/utils/load_env.go`).

## Contributing

- Run `go vet` and `go test ./...` before opening a PR.
- Update this README if you add new endpoints or change existing request/response shapes.

---

For implementation details, check handlers in [internal/api/v1/handlers](internal/api/v1/handlers) and route wiring in [internal/api/v1/routes.go](internal/api/v1/routes.go).
