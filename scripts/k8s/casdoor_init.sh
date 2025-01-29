#!/usr/bin/env bash
set -e

CASDOOR_URL='http://casdoor.auth-core.svc.cluster.local:8000'

ACCESS_TOKEN=$(curl -s -X POST "$CASDOOR_URL"/api/login \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'username=admin&password=123' | jq -r '.access_token')

echo "Access Token: $ACCESS_TOKEN"

# 定义用户参数
USER_NAME="testuser"
USER_PASSWORD="testpass123"

curl -X POST "$CASDOOR_URL"/api/add-user \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "owner": "built-in",
    "name": "'"$USER_NAME"'",
    "password": "'"$USER_PASSWORD"'",
    "displayName": "CI Test User"
  }'

APP_NAME="lubricant"

curl -X POST "$CASDOOR_URL"/api/add-application \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "owner": "built-in",
    "name": "'"$APP_NAME"'",
    "displayName": "CI Test Application",
    "redirectUris": ["'"$CASDOOR_URL"'/callback"]
  }'

# 查询应用详情并提取证书
RESPONSE=$(curl -s -X GET "$CASDOOR_URL/api/get-application?name=$APP_NAME&owner=built-in" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

CLIENT_ID=$(echo "$RESPONSE" | jq -r '.clientId')
CLIENT_SECRET=$(echo "$RESPONSE" | jq -r '.clientSecret')

echo "Client ID: $CLIENT_ID"
echo "Client Secret: $CLIENT_SECRET"

