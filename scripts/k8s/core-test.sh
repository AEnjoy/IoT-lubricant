#!/usr/bin/env bash
set -e

LOGIN_URL="http://casdoor-service.auth-core.svc.cluster.local:8000/api/login?clientId=6551a3584403d5264584&responseType=code&redirectUri=http%3A%2F%2F127.0.0.1%3A8080%2Fapi%2Fv1%2Fsignin&type=code&scope=read&state=casdoor&nonce=&code_challenge_method=&code_challenge="
CALLBACK_URL="http://lubricant-core.lubricant.svc.cluster.local:8080/api/v1/signin"
USER_INFO_URL="http://lubricant-core.lubricant.svc.cluster.local:8080/api/v1/user/info"
COOKIE_FILE="cookie.txt"

echo "Logging in..."
response=$(curl -s -X POST "$LOGIN_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "application": "application_lubricant",
    "organization": "built-in",
    "username": "admin",
    "autoSignin": true,
    "password": "123",
    "signinMethod": "Password",
    "type": "code"
  }')

code=$(echo "$response" | jq -r '.data')
echo "Received code: $code"

if [ -z "$code" ] || [ "$code" = "null" ]; then
  echo "Error: Failed to get code"
  echo "Response: $response"
  exit 1
fi

echo "Getting cookie..."
curl -s -X GET -c "$COOKIE_FILE" \
  "$CALLBACK_URL?code=$code&state=casdoor" > /dev/null

if [ ! -f "$COOKIE_FILE" ]; then
  echo "Error: Failed to get cookie"
  exit 1
fi
echo "Cookie saved to $COOKIE_FILE"

echo "Getting user info..."
curl -s -X GET -b "$COOKIE_FILE" "$USER_INFO_URL"
