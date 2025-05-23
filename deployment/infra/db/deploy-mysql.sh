#!/bin/bash
set -e

# add repo
helm repo add bitnami https://charts.bitnami.com/bitnami

# create secret
lubricant_MYSQL_ROOT_PASSWORD=123456
lubricant_MYSQL_CUSTOM_PASSWORD=123456

# kubectl create ns database

kubectl create secret generic -n database mysql-secret \
  --from-literal=mysql-root-password=${lubricant_MYSQL_ROOT_PASSWORD} \
  --from-literal=mysql-password=${lubricant_MYSQL_CUSTOM_PASSWORD}\
  --from-literal=lubricant-password=${lubricant_MYSQL_CUSTOM_PASSWORD} \
  --from-literal=casdoor-password=${lubricant_MYSQL_CUSTOM_PASSWORD}\
  --from-literal=mysql-username=lubricant

kubectl apply -f deployment/infra/db/mysql.yaml
# helm upgrade --install mysql bitnami/mysql --version 12.3.0 -n database -f deployment/infra/db/values.yaml

sleep 3
mysql_pod=$(kubectl get pods -n database | awk '/mysql/ {print $1}')
bash scripts/function/wait_pod.sh $mysql_pod database
kubectl exec $mysql_pod -n database -- bash -c 'mysqladmin ping -uroot -p123456'

echo "Database initialization..."
kubectl cp deployment/docker/init.d/database.sql $mysql_pod:/tmp/database.sql -n database
kubectl exec $mysql_pod -n database -- bash -c 'mysql -uroot -p123456 < /tmp/database.sql'
kubectl logs -n database -l app=mysql
