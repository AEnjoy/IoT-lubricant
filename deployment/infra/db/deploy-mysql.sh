#!/bin/bash
set -e

# create secret
lubricant_MYSQL_ROOT_PASSWORD=123456
lubricant_MYSQL_CUSTOM_PASSWORD=123456

# kubectl create ns database

kubectl create secret generic -n database mysql-secret \
  --from-literal=mysql-root-password=${lubricant_MYSQL_ROOT_PASSWORD} \
  --from-literal=lubricant-password=${lubricant_MYSQL_CUSTOM_PASSWORD} \
  --from-literal=casdoor-password=${lubricant_MYSQL_CUSTOM_PASSWORD}

kubectl apply -f deployment/infra/db/mysql.yaml

sleep 3
mysql_pod=$(kubectl get pods -n database | awk '/mysql/ {print $1}')
bash scripts/function/wait_pod.sh $mysql_pod database

echo "Database initialization..."
kubectl cp deployment/infra/database.sql $mysql_pod:/tmp/database.sql -n database
kubectl exec -it $mysql_pod -n database -- bash -c 'mysql -uroot -p123456 < /tmp/database.sql'
