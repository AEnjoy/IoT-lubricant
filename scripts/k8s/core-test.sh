#!/usr/bin/env bash

pod_name=$1

LOGIN_URL="http://casdoor-service.auth-core.svc.cluster.local:8000/api/login?clientId=6551a3584403d5264584&responseType=code&redirectUri=http%3A%2F%2Flubricant-core.lubricant.svc.cluster.local%3A8080%2Fapi%2Fv1%2Fsignin&type=code&scope=read&state=casdoor&nonce=&code_challenge_method=&code_challenge="
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
response=$(curl -s -X GET -c "$COOKIE_FILE" \
  "$CALLBACK_URL?code=$code&state=casdoor")
if [ $? -ne 0 ]; then
  echo "Error: Failed to get cookie curl request failed"
  echo "Response: $response"
  kubectl logs "$pod_name" -n lubricant
  exit 1
fi

msg=$(echo "$response" | jq -r '.msg')

if [ "$msg" != "success" ]; then
  echo "Error: Login failed, msg=$msg"
  echo "Response: $response"
  kubectl logs "$pod_name" -n lubricant
  exit 1
fi

if [ ! -f "$COOKIE_FILE" ]; then
  echo "Error: Failed to get cookie"
  echo "File $COOKIE_FILE does not exist"
  kubectl logs "$pod_name" -n lubricant
  exit 1
fi
echo "Cookie saved to $COOKIE_FILE"

cat $COOKIE_FILE

echo "Getting user info..."
curl -s -X GET -b "$COOKIE_FILE" "$USER_INFO_URL"
if [ $? -ne 0 ]; then
  echo "Error: Failed to get user info"
  kubectl logs "$pod_name" -n lubricant
  exit 1
fi
