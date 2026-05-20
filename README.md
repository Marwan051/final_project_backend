# Backend Setup

## Environment

1. Copy `.env.example` to `.env`.
2. Fill required values:
   - `DB_URL`
   - `PORT`
   - `ENV`
   - `ROUTING_SERVICE_ADDR`
   - `GRPC_REQUEST_TIMEOUT` to override the default 10s gRPC deadline
3. Configure Supabase auth:
   - `SUPABASE_URL`
   - `SUPABASE_SECRET_KEY` or `SUPABASE_SERVICE_ROLE_KEY`
   - `SUPABASE_PUBLISHABLE_KEY` if you also use the same project from a browser/client app

## Supabase Auth

This backend verifies Supabase access tokens server-side using the Supabase project URL and a server key.

Supabase's current guidance is to use the publishable key for client-side usage and the secret key for server-side usage. Legacy `anon` and `service_role` keys still work, but the new keys are preferred.

Use the project URL and one of the server keys in your environment:

```bash
SUPABASE_URL="https://your-project-id.supabase.co"
SUPABASE_SECRET_KEY="sb_secret_..."
```

If your project still uses legacy keys, `SUPABASE_SERVICE_ROLE_KEY` can be used as a fallback.

## Run

Run with hot reload:

```bash
air
```

## API Routes

`GET /health` is public. Everything under `/api/v1/*` requires a valid Supabase access token, and the two routing maintenance endpoints below also require an admin claim.

| Method | Path | Purpose | Access |
| --- | --- | --- | --- |
| `GET` | `/health` | Returns the service health status and current timestamp. | Public |
| `POST` | `/api/v1/route` | Calculates the best route between two coordinate pairs. | Authenticated |
| `POST` | `/api/v1/agent/query` | Sends a message to the agent and returns its reply. | Authenticated |
| `POST` | `/api/v1/geocode` | Converts a text address into latitude and longitude coordinates. | Authenticated |
| `POST` | `/api/v1/nearby-trips` | Finds trips near a point within a radius. | Authenticated |
| `POST` | `/api/v1/traffic/trigger` | Starts a traffic update job manually. | Authenticated |
| `GET` | `/api/v1/traffic/status` | Returns the current traffic update status. | Authenticated |
| `POST` | `/api/v1/traffic/update-trip` | Submits a trip to affect traffic modelling. | Authenticated |
| `POST` | `/api/v1/traffic/street` | Returns the calculated traffic load for one street. | Authenticated |
| `GET` | `/api/v1/traffic/streets` | Lists all tracked streets and their traffic levels. | Authenticated |
| `POST` | `/api/v1/routing/reload-prefix-times` | Reloads prefix times for the routing graph. | Admin only |
| `POST` | `/api/v1/routing/rebuild-network` | Rebuilds the routing network from source data. | Admin only |

## Test Supabase Auth End To End

All `/api/v1/*` endpoints require auth, while `/health` is a quick unauthenticated liveness check.

1. Make sure password sign-in is enabled in your Supabase project.
2. Get your project URL and publishable key from the Supabase dashboard.
3. Mint an access token with the Supabase Auth password grant:

```bash
curl -sS -X POST "${SUPABASE_URL}/auth/v1/token?grant_type=password" \
   -H "apikey: ${SUPABASE_PUBLISHABLE_KEY}" \
   -H "Authorization: Bearer ${SUPABASE_PUBLISHABLE_KEY}" \
   -H "Content-Type: application/json" \
   --data-binary '{"email":"<EMAIL>","password":"<PASSWORD>"}'
```

4. Copy the `access_token` from the response and call one of the protected endpoints:

```bash
TOKEN="<ACCESS_TOKEN>"
curl -i -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/api/v1/route
```

Expected results for a protected route such as `/api/v1/route`:

- Missing token: `401 Unauthorized`
- Invalid token: `401 Unauthorized`
- Valid token: `200 OK` with JSON body

Quick negative checks:

```bash
curl -i http://localhost:3000/api/v1/route
curl -i -H "Authorization: Bearer not-a-real-token" http://localhost:3000/api/v1/route
```
