#!/usr/bin/env bash
set -e

pods=$(kubectl get pods --all-namespaces --field-selector=status.phase=Running | grep -E 'CrashLoopBackOff|Error' | awk '{print $1 " " $2}')
if [ -z "$pods" ]; then
   echo "No pods in Error or CrashLoopBackOff status found."
   exit 0
fi
while read -r namespace pod; do
    echo "Pod $pod in namespace $namespace is in CrashLoopBackOff or Error status."
    kubectl describe pod "$pod" -n "$namespace"
    kubectl logs -n $namespace $pod
done <<< "$pods"
