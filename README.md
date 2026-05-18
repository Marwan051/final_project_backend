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

## Test Supabase Auth End To End

All `/api/v1/*` endpoints are protected by auth middleware, so `/api/v1/health` is a quick auth check.

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

4. Copy the `access_token` from the response and call the protected endpoint:

```bash
TOKEN="<ACCESS_TOKEN>"
curl -i -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/api/v1/health
```

Expected results:

- Missing token: `401 Unauthorized`
- Invalid token: `401 Unauthorized`
- Valid token: `200 OK` with JSON body

Quick negative checks:

```bash
curl -i http://localhost:3000/api/v1/health
curl -i -H "Authorization: Bearer not-a-real-token" http://localhost:3000/api/v1/health
```
