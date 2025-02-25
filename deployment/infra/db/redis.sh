#!/bin/bash
set -e

kubectl create secret generic -n=database redis-secret \
  --from-literal=redis-password=123456 \

kubectl apply -f deployment/infra/db/redis.yaml
