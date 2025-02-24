#!/usr/bin/env bash
set -e

POD_NAME='nginx' # $(kubectl get pods -l app=nginx -o jsonpath="{.items[0].metadata.name}")

bash scripts/function/wait_pod.sh "$POD_NAME" "default"
#if [ -z "$POD_NAME" ]; then
#  echo "No running Nginx Pod found."
#  exit 1
#fi

echo "Found Nginx Pod: $POD_NAME"

kubectl exec "$POD_NAME" -- mkdir -p /root/k8s
kubectl cp scripts/k8s "$POD_NAME:/root"

if [ $? -eq 0 ]; then
  echo "Files successfully copied to $POD_NAME:/root/k8s"
else
  echo "Failed to copy files to $POD_NAME:/root/k8s"
  exit 1
fi
