#!/usr/bin/env bash
set -e
kubectl create ns lubricant
kubectl create secret generic lubricant-secrets -n lubricant\
  --from-literal=DB_PASSWORD='123456' \
  --from-literal=AUTH_CLIENT_SECRET='dd9657c7b8cc10a72f77b283253b3a0a31b91175'
