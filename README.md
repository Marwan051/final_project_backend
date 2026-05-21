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
# requires `air` for hot reload (optional: brew install cosmtrek/tap/air)
air
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

Base URL (production): `https://api.example.com/api/v1/`

**Common headers for all endpoints:**

```
Authorization: Bearer <access_token>    (required for /api/v1/*, not for /health)
Content-Type: application/json          (required for POST bodies)
```

**Common error responses:**

- `400 Bad Request` — Invalid JSON or missing required fields
- `401 Unauthorized` — Missing or invalid `Authorization` header
- `403 Forbidden` — User lacks permission (e.g., not admin for admin-only endpoints)
- `500 Internal Server Error` — Server error or downstream service failure

---

### System Endpoints

#### `GET /health`

Health check. **No authentication required.**

**Response (200 OK):**

```json
{
  "status": "ok",
  "timestamp": "2026-05-21T18:56:54Z"
}
```

---

### Routing Endpoints

#### `POST /api/v1/route`

Calculates optimal routes between two coordinates using multi-modal transit.

**Request body:**

| Parameter        | Type    | Required | Description                                      |
| ---------------- | ------- | -------- | ------------------------------------------------ |
| `start_lat`      | float64 | yes      | Start latitude                                   |
| `start_lon`      | float64 | yes      | Start longitude                                  |
| `end_lat`        | float64 | yes      | End latitude                                     |
| `end_lon`        | float64 | yes      | End longitude                                    |
| `max_transfers`  | int32   | no       | Max transfers allowed (default: 3)               |
| `walking_cutoff` | int32   | no       | Walking cutoff in meters (default: 500)          |
| `priority`       | string  | no       | Route priority (`fastest`, `cheapest`, etc.)     |
| `top_k`          | int32   | no       | Number of result journeys to return (default: 5) |
| `filters`        | object  | no       | Optional filters for modes and streets           |
| `weights`        | object  | no       | Custom weights for routing optimization          |

**Example request:**

```bash
curl -sS -X POST http://localhost:3000/api/v1/route \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  --data-raw '{
    "start_lat": 40.7128,
    "start_lon": -74.0060,
    "end_lat": 40.7580,
    "end_lon": -73.9855,
    "max_transfers": 2,
    "walking_cutoff": 600,
    "top_k": 3
  }'
```

**Response (200 OK):**

```json
{
  "geometry_encoding": "polyline",
  "selected_priority": "fastest",
  "num_journeys": 3,
  "journeys": [
    {
      "id": 1,
      "text_summary": "Walk, then take the subway",
      "summary": {
        "total_time_minutes": 18,
        "walking_distance_meters": 300,
        "transit_distance_meters": 5200,
        "total_distance_meters": 5500,
        "transfers": 0,
        "cost": 2.75,
        "modes_en": ["walking", "subway"],
        "main_streets_en": ["Broadway", "5th Ave"]
      },
      "legs": [
        {
          "type": "walking",
          "distance_meters": 300,
          "duration_minutes": 4,
          "polyline": "yvxeFjy..."
        },
        {
          "type": "transit",
          "mode_en": "subway",
          "trip_id": "trip_123",
          "route_short_name": "A",
          "headsign": "Far Rockaway",
          "fare": 2.75,
          "distance_meters": 5200,
          "duration_minutes": 14,
          "from_stop": {
            "stop_id": "stop_101",
            "name": "42nd Street",
            "coord": [40.756, -73.9903]
          },
          "to_stop": {
            "stop_id": "stop_205",
            "name": "34th Street",
            "coord": [40.7505, -73.988]
          }
        }
      ]
    }
  ],
  "start_trips_found": 12,
  "end_trips_found": 8,
  "total_routes_found": 96,
  "total_after_dedup": 3
}
```

**Error responses:**

- `400 Bad Request` — Missing or invalid coordinates
- `500 Internal Server Error` — Routing service unreachable or error

---

#### `POST /api/v1/routing/reload-prefix-times` (Admin only)

Reloads prefix times for the routing graph. Useful after updating source data.

**Request body:** Empty or `{}`

**Response (200 OK):**

```json
{
  "status": "success",
  "message": "Prefix times reloaded successfully",
  "trips_reloaded": 1542
}
```

**Error responses:**

- `403 Forbidden` — User is not an admin
- `500 Internal Server Error` — Reload operation failed

---

#### `POST /api/v1/routing/rebuild-network` (Admin only)

Rebuilds the entire routing network from source data. **Use with caution** — can take time.

**Request body:** Empty or `{}`

**Response (200 OK):**

```json
{
  "status": "success",
  "message": "Network rebuilt successfully",
  "trips_reloaded": 5420
}
```

**Error responses:**

- `403 Forbidden` — User is not an admin
- `500 Internal Server Error` — Rebuild operation failed

---

### Agent Endpoints

#### `POST /api/v1/agent/query`

Sends a message to the agent and receives a reply. Useful for natural language queries about routes, trips, etc.

**Request body:**

| Parameter    | Type   | Required | Description                            |
| ------------ | ------ | -------- | -------------------------------------- |
| `user_query` | string | yes      | The user message or query              |
| `session_id` | string | no       | Session ID for conversation continuity |

**Example request:**

```bash
curl -sS -X POST http://localhost:3000/api/v1/agent/query \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  --data-raw '{
    "user_query": "What is the fastest way to get from Central Park to Times Square?",
    "session_id": "sess_abc123"
  }'
```

**Response (200 OK):**

```json
{
  "answer": "The fastest route from Central Park to Times Square takes about 12 minutes via the subway on the A line or walking along Broadway.",
  "session_id": "sess_abc123"
}
```

**Error responses:**

- `400 Bad Request` — Missing or empty `user_query`
- `503 Service Unavailable` — Agent service is down
- `502 Bad Gateway` — Agent service error

---

### Geocoding Endpoints

#### `POST /api/v1/geocode`

Converts an address string into geographic coordinates (latitude, longitude).

**Request body:**

| Parameter  | Type    | Required | Description                                  |
| ---------- | ------- | -------- | -------------------------------------------- |
| `address`  | string  | yes      | Address or place name                        |
| `language` | string  | no       | Language code (e.g., `en`, `ar`)             |
| `bias`     | bool    | no       | Whether to bias results toward user location |
| `user_lat` | float64 | no       | User's current latitude (for bias)           |
| `user_lng` | float64 | no       | User's current longitude (for bias)          |

**Example request:**

```bash
curl -sS -X POST http://localhost:3000/api/v1/geocode \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  --data-raw '{
    "address": "1600 Amphitheatre Parkway, Mountain View, CA",
    "language": "en",
    "bias": true,
    "user_lat": 40.7128,
    "user_lng": -74.0060
  }'
```

**Response (200 OK):**

```json
{
  "success": true,
  "query": "1600 Amphitheatre Parkway, Mountain View, CA",
  "language": "en",
  "bias": true,
  "count": 1,
  "results": [
    {
      "latitude": 37.4224764,
      "longitude": -122.0842499,
      "formatted_address": "1600 Amphitheatre Parkway, Mountain View, CA 94043, USA"
    }
  ],
  "error": ""
}
```

**Error responses:**

- `400 Bad Request` — Missing address
- `500 Internal Server Error` — Geocoding service error

---

### DB Tools Endpoints

#### `POST /api/v1/nearby-trips`

Finds all transit trips within a specified radius of a given point.

**Request body:**

| Parameter  | Type    | Required | Description                                          |
| ---------- | ------- | -------- | ---------------------------------------------------- |
| `lat`      | float64 | yes      | Center latitude                                      |
| `lon`      | float64 | yes      | Center longitude                                     |
| `radius_m` | float64 | yes      | Search radius in meters                              |
| `starts`   | bool    | no       | If true, query trip start stops; if false, all stops |
| `epsg`     | int32   | no       | EPSG code for spatial query (default: 4326)          |

**Example request:**

```bash
curl -sS -X POST http://localhost:3000/api/v1/nearby-trips \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  --data-raw '{
    "lat": 40.7128,
    "lon": -74.0060,
    "radius_m": 500,
    "starts": true
  }'
```

**Response (200 OK):**

```json
{
  "lat": 40.7128,
  "lon": -74.006,
  "radius_m": 500,
  "starts": true,
  "count": 5,
  "trips": [
    {
      "trip_id": "trip_001",
      "route_id": "route_A",
      "trip_headsign": "Far Rockaway",
      "trip_headsign_ar": "بعيد روكاواي",
      "direction_id": 0,
      "route_short_name": "A",
      "route_short_name_ar": "أ",
      "route_name": "Eighth Avenue Line",
      "route_name_ar": "خط الجادة الثامنة",
      "distance_m": 150.5,
      "closest_stop_id": "stop_101",
      "closest_stop_name": "42nd Street",
      "closest_stop_name_ar": "شارع 42",
      "closest_stop_lat": 40.756,
      "closest_stop_lon": -73.9903,
      "closest_stop_sequence": 1
    }
  ]
}
```

**Error responses:**

- `400 Bad Request` — Missing or invalid lat/lon/radius
- `500 Internal Server Error` — Database query failed

---

### Traffic Endpoints

#### `POST /api/v1/traffic/trigger`

Manually triggers a background traffic update job.

**Request body:**

| Parameter            | Type     | Required | Description                                    |
| -------------------- | -------- | -------- | ---------------------------------------------- |
| `trip_ids`           | string[] | no       | Specific trip IDs to update; empty = all trips |
| `notify_routing_api` | bool     | no       | Whether to notify routing service after update |

**Example request:**

```bash
curl -sS -X POST http://localhost:3000/api/v1/traffic/trigger \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  --data-raw '{
    "trip_ids": ["trip_001", "trip_002"],
    "notify_routing_api": true
  }'
```

**Response (200 OK):**

```json
{
  "status": "started",
  "trips_updated": 2,
  "trips_failed": 0,
  "message": "Traffic update triggered for 2 trips"
}
```

**Error responses:**

- `400 Bad Request` — Invalid request format
- `500 Internal Server Error` — Update job failed to start

---

#### `GET /api/v1/traffic/status`

Gets the current traffic update status.

**Request body:** None (GET request)

**Example request:**

```bash
curl -H "Authorization: Bearer ${TOKEN}" \
  http://localhost:3000/api/v1/traffic/status
```

**Response (200 OK):**

```json
{
  "status": "idle",
  "last_update": "2026-05-21T18:45:30Z",
  "trips_in_data": 1542,
  "is_running": false
}
```

---

#### `POST /api/v1/traffic/update-trip`

Submits a specific trip for traffic modelling update.

**Request body:**

| Parameter            | Type   | Required | Description                         |
| -------------------- | ------ | -------- | ----------------------------------- |
| `trip_id`            | string | yes      | Trip ID to update                   |
| `notify_routing_api` | bool   | no       | Notify routing service after update |

**Example request:**

```bash
curl -sS -X POST http://localhost:3000/api/v1/traffic/update-trip \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  --data-raw '{
    "trip_id": "trip_001",
    "notify_routing_api": true
  }'
```

**Response (200 OK):**

```json
{
  "status": "success",
  "trips_updated": 1,
  "trips_failed": 0,
  "message": "Trip updated successfully"
}
```

---

#### `POST /api/v1/traffic/street`

Fetches traffic load and route information for a specific street.

**Request body:**

| Parameter       | Type   | Required | Description                       |
| --------------- | ------ | -------- | --------------------------------- |
| `name`          | string | yes      | Street name                       |
| `language`      | string | no       | Language code (`en`, `ar`)        |
| `max_waypoints` | int32  | no       | Max waypoints to include in route |

**Example request:**

```bash
curl -sS -X POST http://localhost:3000/api/v1/traffic/street \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  --data-raw '{
    "name": "Broadway",
    "language": "en",
    "max_waypoints": 10
  }'
```

**Response (200 OK):**

```json
{
  "street": "Broadway",
  "street_ar": "برودواي",
  "segments": 5,
  "waypoints_used": 8,
  "total_distance_km": 3.45,
  "total_duration_min": 12.5,
  "legs": [
    {
      "distance_m": 700,
      "distance_text": "700 m",
      "duration_seconds": 180,
      "duration_text": "3 min"
    }
  ],
  "routes": [
    {
      "label": "Current traffic",
      "distance_m": 3450,
      "distance_text": "3.45 km",
      "duration_seconds": 750,
      "duration_text": "12.5 min"
    }
  ],
  "error": ""
}
```

---

#### `GET /api/v1/traffic/streets`

Lists all streets currently tracked in the traffic system.

**Request body:** None (GET request)

**Example request:**

```bash
curl -H "Authorization: Bearer ${TOKEN}" \
  http://localhost:3000/api/v1/traffic/streets
```

**Response (200 OK):**

```json
{
  "count": 45,
  "streets": [
    {
      "name": "Broadway",
      "name_ar": "برودواي",
      "aliases": ["Bwy"],
      "segments": 8,
      "total_length_km": 21.5
    },
    {
      "name": "5th Avenue",
      "name_ar": "جادة الخامسة",
      "aliases": [],
      "segments": 12,
      "total_length_km": 18.3
    }
  ]
}
```

---

### Admin Endpoints Summary

The following endpoints require an **admin claim** in the user's Supabase token:

| Endpoint                              | Method | Purpose                     |
| ------------------------------------- | ------ | --------------------------- |
| `/api/v1/routing/reload-prefix-times` | POST   | Reload routing prefix times |
| `/api/v1/routing/rebuild-network`     | POST   | Rebuild routing network     |

To verify admin access, check your Supabase user's custom claims or `app_metadata.role` field.

## Swagger / OpenAPI Documentation

In development (`ENV=dev`) the server serves a Swagger UI at `/docs/`:

```bash
curl http://localhost:3000/docs/
```

The Swagger spec is auto-generated from handler comments and proto definitions. Check [docs/swagger.yaml](docs/swagger.yaml) for the complete spec.

## Authentication & Testing

### Obtaining a Supabase Access Token

1. **Ensure password sign-in is enabled** in your Supabase project dashboard.

2. **Get your Supabase credentials:**
   - Project URL from Supabase dashboard (e.g., `https://xxxx.supabase.co`)
   - Publishable key (anon key)

3. **Mint an access token via password grant:**

```bash
export SUPABASE_URL="https://your-project-id.supabase.co"
export SUPABASE_PUBLISHABLE_KEY="eyJhbGc..."

curl -sS -X POST "${SUPABASE_URL}/auth/v1/token?grant_type=password" \
  -H "apikey: ${SUPABASE_PUBLISHABLE_KEY}" \
  -H "Authorization: Bearer ${SUPABASE_PUBLISHABLE_KEY}" \
  -H "Content-Type: application/json" \
  --data-raw '{"email":"user@example.com","password":"your-password"}' \
  | jq '.access_token' -r
```

4. **Use the token for requests:**

```bash
export TOKEN="eyJhbGc..." # from step 3
curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/api/v1/traffic/status
```

### Testing Endpoints

#### Test unauthenticated access (should succeed):

```bash
curl -i http://localhost:3000/health
```

Expected: `200 OK` with `{"status":"ok","timestamp":"..."}`

#### Test authenticated access without token (should fail):

```bash
curl -i http://localhost:3000/api/v1/traffic/status
```

Expected: `401 Unauthorized`

#### Test with invalid token (should fail):

```bash
curl -i -H "Authorization: Bearer invalid-token" http://localhost:3000/api/v1/traffic/status
```

Expected: `401 Unauthorized`

#### Test with valid token (should succeed):

```bash
export TOKEN="your-valid-token"
curl -i -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/api/v1/traffic/status
```

Expected: `200 OK` with status response

#### Full integration test (route finding):

```bash
export TOKEN="your-valid-token"
curl -sS -X POST http://localhost:3000/api/v1/route \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  --data-raw '{
    "start_lat": 40.7128,
    "start_lon": -74.0060,
    "end_lat": 40.7580,
    "end_lon": -73.9855,
    "top_k": 2
  }' | jq .
```

### Local Development (Auth Disabled)

If you need to skip authentication locally:

```bash
export ENV=dev
export DISABLE_AUTH=true
```

Then all `/api/v1/*` endpoints can be called without tokens. However, admin-only endpoints will still require the admin claim check (see `requireAdmin` middleware).

## Development & Testing

### Run tests:

```bash
go test ./...
```

### Run with coverage:

```bash
go test -cover ./...
```

### Lint code:

```bash
go vet ./...
```

### Format code:

```bash
go fmt ./...
```

### Build Docker image:

```bash
docker build -t final-project-backend-server:latest .
```

### Push to registry:

```bash
docker tag final-project-backend-server:latest marwan051/final-project-backend-server:latest
docker push marwan051/final-project-backend-server:latest
```

## Contributing

- Run `go vet` and `go test ./...` before opening a PR.
- Update this README if you add new endpoints or change request/response shapes.
- Follow Go conventions: `CamelCase` for exported symbols, `snake_case` in JSON tags.
- Add Swagger comments to new handlers (see existing handler files for examples).
- Ensure proto changes are compiled: `cd internal/service/<service>/proto && protoc ...`

---

For implementation details, check handlers in [internal/api/v1/handlers](internal/api/v1/handlers) and route wiring in [internal/api/v1/routes.go](internal/api/v1/routes.go).
