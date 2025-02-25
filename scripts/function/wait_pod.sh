#!/usr/bin/env bash
set -e

pod_name=$1
namespace=$2
timeout=120
attempts=3

echo "Waiting for pod '$pod_name' to be ready..."

for ((i = 1; i <= attempts; i++)); do
    echo "Waiting for pod '$pod_name' in namespace $namespace to be ready... (attempt $i of $attempts)"
    if kubectl wait --for=condition=ready pod "$pod_name" -n "$namespace" --timeout="${timeout}s"; then
        echo "Pod '$pod_name' in namespace $namespace is ready."
        exit 0
    else
        if [[ $i -lt $attempts ]]; then
            echo "Pod '$pod_name' in namespace $namespace is not ready, retrying in 1 second..."
            sleep 1
        else
            echo "Pod '$pod_name' in namespace $namespace failed to be ready after $attempts attempts."
            echo "Printing pod description and logs:"
            kubectl describe pod -n $namespace $pod_name
            kubectl logs -n $namespace $pod_name
            exit 1
        fi
    fi
done
