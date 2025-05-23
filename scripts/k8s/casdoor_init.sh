#!/usr/bin/env bash
set -e

# there is not jq command at nginx pod
apt-get update && apt-get install -y jq

CASDOOR_URL='http://casdoor-service.auth-core.svc.cluster.local:8000'

# Create App
curl -f -X POST "$CASDOOR_URL/api/add-application?username=built-in/admin&password=123" \
 -H "Content-Type: application/json" -d '@create_app.json'

#  Set Cert
time=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
jq ".createdTime = \"$time\"" update_cert.json > update_cert.json.new
mv update_cert.json.new update_cert.json
curl -f -X POST "$CASDOOR_URL/api/update-cert?id=admin/cert-built-in&username=built-in/admin&password=123" \
 -H "Content-Type: application/json" -d '@update_cert.json'

# Get Cert
#curl -s "$CASDOOR_URL/api/get-cert?id=admin/cert-built-in&username=built-in/admin&password=123" \
#  | jq -r '.data.certificate'| echo -e "$(cat)" > ./crt.pem
#openssl x509 -in crt.pem -text -noout
