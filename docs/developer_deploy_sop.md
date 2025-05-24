# how to prepare develop environment:

## using docker to deploy:

you need to deploy mysql, redis , nats(or other message queue), and casdoor for development. 

1. deploy mysql:
```bash
docker run -d --name=mysql-server -p 3306:3306  -e MYSQL_ROOT_PASSWORD=123456 mysql 
docker cp deployment/docker/init.d/database.sql mysql-server:/tmp/database.sql
docker exec mysql-server bash -c 'mysql -uroot -p123456 < /tmp/database.sql'
```

2. deploy redis:
```bash
docker run -d --name myredis -p 6379:6379 redis --requirepass "123456"
```

3. deploy nats:
```bash
docker run --name nats --rm -d -p 4222:4222 -p 8222:8222 nats --http_port 8222
```

4. deploy casdoor:
```bash
docker run --name=casdoor -d -p 8000:8000 casbin/casdoor-all-in-one 
curl "127.0.0.1:8000/api/add-application?username=built-in/admin&password=123" \
 -H "Content-Type: application/json" -d '@scripts/k8s/create_app.json'
 
curl -s "127.0.0.1:8000/api/get-cert?id=admin/cert-built-in&username=built-in/admin&password=123" | jq -r '.data.certificate' > ./crt.pem
```

5. start Core - server

copy `cmd/core/core.env.template` to `cmd/core/core.env`, and edit the configs.

running core.sh


## deploy to kubernetes:

Notes: **It is difficult to develop in kubernetes. It is recommended to use this record only during testing**

1. set up kind or the other kubernetes cluster
2. running

    a. deployment/infra/nsinit.sh

    b. deployment/infra/secret.sh

    c. deployment/infra/db/redis.sh
    
    d. deployment/infra/db/deploy-mysql.sh
    
    e. deployment/infra/db/tdengine.sh
3. deploy lubricant by yamls:`kubectl apply -f  deployment/infra`

## deploy to docker-compose:

Execute those commands at project root dir: (need root privilege)

```bash
make build-all
make copy-files
FAST_BUILD=1 make docker-build -j

cd deployment/docker
docker compose up -d
```
