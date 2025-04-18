# curl测试接口流程

1. 登录： 获取回调
```shell
curl -X POST 'http://127.0.0.1:8000/api/login?clientId=6551a3584403d5264584&responseType=code&redirectUri=http%3A%2F%2F127.0.0.1%3A8080%2Fapi%2Fv1%2Fsignin&type=code&scope=read&state=casdoor&nonce=&code_challenge_method=&code_challenge=' \
-H 'Content-Type: application/json' \
-d '{
    "application": "application_lubricant",
    "organization": "built-in",
    "username": "admin",
    "autoSignin": true,
    "password": "123",
    "signinMethod": "Password",
    "type": "code"
}'
```

输出:
```json
{
  "status": "ok",
  "msg": "",
  "sub": "",
  "name": "",
  "data": "61bedefc366afc0d8a53",
  "data2": false
}
```
其中，data为code

2. 使用curl访问回调接口，并获取cookie
```shell
curl  'http://127.0.0.1:8080/api/v1/signin?code=61bedefc366afc0d8a53&state=casdoor' \
  -X GET  -v -i -c cookie.txt
```

3. 后续请求携带cookie:
```shell
curl 'http://127.0.0.1:8080/api/v1/user/info' -X GET  -v -i -b cookie.txt
```


## 接口

添加网关：
```shell
curl 'http://127.0.0.1:8080/api/v1/gateway/internal/gateway' -X POST -b cookie.txt -H 'Content-Type: application/json'  -d @scripts/test/request/create_gateway_internal.json
```

添加agent
```shell
curl 'http://127.0.0.1:8080/api/v1/gateway/2988e18e-6861-4d7d-8be1-5c539faad0f1/agent/internal/add' -X POST -b cookie.txt -H 'Content-Type: application/json' -d @scripts/test/request/add_agent_internal_request.json
```

对Agent的操作：

operator:

- start-gather
- stop-gather
- start-agent
- stop-agent
- get-openapidoc

```shell
curl 'http://127.0.0.1:8080/api/v1/agent/operator?agent-id=c9c603ff-5a9e-4362-abbc-284045aa2cf3&gateway-id=2988e18e-6861-4d7d-8be1-5c539faad0f1&operator=start-gather' -X GET  -v -i -b cookie.txt
```

设置agent

```shell
curl 'http://127.0.0.1:8080/api/v1/agent/set?gateway-id=gateway-id-123' -X POST -b cookie.txt -H 'Content-Type: application/json' \
 -d @scripts/test/request/set_agent_data.json
```

在设置agent前，请先编辑set_agent_data.json，将agentID设置为我们要设置的目标agentID。


异步任务查询：

```shell
curl 'http://127.0.0.1:8080/api/v1/task/query?taskId=' -X GET -b cookie.txt
```
