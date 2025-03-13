#!/usr/bin/env bash
set -e

kubectl create secret generic -n=database tsdatabase-secret \
  --from-literal=root-password=123456 \

kubectl apply -f deployment/infra/db/tdengine.yaml
