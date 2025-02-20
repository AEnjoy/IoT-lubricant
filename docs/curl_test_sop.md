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


