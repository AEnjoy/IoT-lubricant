#!/usr/bin/env bash
set -e

POD_NAME='nginx'
bash scripts/function/wait_pod.sh "$POD_NAME" "default"

kubectl exec "$POD_NAME" -- mkdir -p /root/k8s
kubectl cp scripts/k8s "$POD_NAME:/root"

if [ $? -eq 0 ]; then
  echo "Files successfully copied to $POD_NAME:/root/k8s"
else
  echo "Failed to copy files to $POD_NAME:/root/k8s"
  exit 1
fi
