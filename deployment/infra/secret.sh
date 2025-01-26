#!/usr/bin/env bash
set -e
kubectl create ns lubricant
kubectl create secret generic lubricant-secrets -n lubricant\
  --from-literal=DB_PASSWORD='123456' \
  --from-literal=AUTH_CLIENT_SECRET='<your-client-secret>'
