# Backend Setup

## Environment

1. Copy `.env.example` to `.env`.
2. Fill required values:
   - `DB_URL`
   - `PORT`
   - `ENV`
   - `ROUTING_SERVICE_ADDR`
3. For Firebase Auth token verification, set this when needed:
   - `FIREBASE_PROJECT_ID` (optional but recommended for local ADC)

## Firebase Auth (No JSON Key File Required)

This project uses Firebase Admin SDK with Application Default Credentials (ADC).
It does not require a service-account JSON key file for normal local and cloud workflows.

Recommended local setup from latest Google/Firebase docs:

```bash
gcloud auth application-default login
```

Then set your Firebase project ID for local dev if not already present:

```bash
export FIREBASE_PROJECT_ID="your-project-id"
```

Alternative accepted by Firebase Admin SDK:

```bash
export GOOGLE_CLOUD_PROJECT="your-project-id"
```

Notes:

- In Google Cloud runtimes, attached service accounts are preferred.
- Service-account JSON keys still work, but are not the recommended default.

## Run

Run with hot reload:

```bash
air
```

## Test Firebase Auth End To End

All `/api/v1/*` endpoints are protected by auth middleware, so `/api/v1/health` is a quick auth check.

1. Ensure Email/Password sign-in is enabled in Firebase Auth.
2. Get your Firebase Web API key from Firebase project settings.
3. Mint an ID token using the helper script:

```bash
./test/get_firebase_id_token.sh "<FIREBASE_WEB_API_KEY>" "<EMAIL>" "<PASSWORD>"
```

4. Save token and call the protected endpoint:

```bash
TOKEN="$(./test/get_firebase_id_token.sh "<FIREBASE_WEB_API_KEY>" "<EMAIL>" "<PASSWORD>")"
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
