#!/usr/bin/env bash
set -euo pipefail

if [ "$#" -ne 3 ]; then
  echo "Usage: $0 <FIREBASE_WEB_API_KEY> <EMAIL> <PASSWORD>" >&2
  echo "Example: $0 \"AIza...\" \"user@example.com\" \"secret\"" >&2
  exit 1
fi

api_key="$1"
email="$2"
password="$3"

response="$(curl -sS -X POST \
  "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=${api_key}" \
  -H "Content-Type: application/json" \
  --data-binary @- <<EOF
{
  "email": "${email}",
  "password": "${password}",
  "returnSecureToken": true
}
EOF
)"

if printf "%s" "$response" | grep -q '"error"'; then
  echo "Firebase sign-in failed:" >&2
  printf "%s\n" "$response" >&2
  exit 1
fi

token="$(printf "%s" "$response" | grep -o '"idToken"[[:space:]]*:[[:space:]]*"[^"]*"' | head -n1 | cut -d '"' -f4)"

if [ -z "$token" ]; then
  echo "Could not parse idToken from response:" >&2
  printf "%s\n" "$response" >&2
  exit 1
fi

printf "%s\n" "$token"
