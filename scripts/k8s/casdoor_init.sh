#!/usr/bin/env bash
set -e

apt-get update && apt-get install -y jq

CASDOOR_URL='http://casdoor-service.auth-core.svc.cluster.local:8000'

# 创建应用
curl -X POST "$CASDOOR_URL/api/add-application?username=built-in/admin&password=123" \
 -H "Content-Type: application/json" -d '@create_app.json'

# 查询应用详情并提取证书
curl -s "$CASDOOR_URL/api/get-cert?id=admin/cert-built-in&username=built-in/admin&password=123" | jq -r '.data.certificate' > ./crt.pem

# cp ./crt.pem /etc/casdoor/public.pem
