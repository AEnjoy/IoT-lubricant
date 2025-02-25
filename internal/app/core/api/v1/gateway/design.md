# Gateways Api Design

**Notice:**

The current design document is still working and will be moved to `/doc` when it is completed

## Apis

root route: /api/v1/gateway

Gateway manage:

- [x]  add gateway host: `POST /add-host` and `POST /host`
- [x]  edit gateway host info: `POST /edit-host` and `PUT /host`
- [ ]  remove gateway host: `DELETE /remove-host` and `DELETE /host` -> uninstall.
- [x]  get gateway host info: `GET /get-host?hostid` and `GET /host?hostid`
- [x]  list gateway hosts info: `GET /list-hosts?userid` and `GET /hosts?userid`
- [x]  deploy gateway instance: `POST /deploy-instance`
- [ ]  update gateway instance: `POST /update-instance`
- [ ]  uninstall gateway instance: `POST /uninstall-instance`
(internal)
- [x]  add gateway: `POST /internal/add-gateway` and `POST /internal/gateway`
- [x]  remove gateway: `POST /internal/remove-gateway` and `DELETE /internal/gateway`
- [ ]  remove gateway host: `POST /internal/remove-gateway-host`

Gateway functions:

- [x]  get gateway status: `GET /status?gatewayid`
- [x]  get error logs: `GET /error-logs?gatewayid`
- [x]  set gateway instance config: `POST /set-config?gatewayid`
- [x]  get gateway instance config: `GET /get-config?gatewayid`

Gateway agent instance manage:

root route: /api/v1/gateway/:gatewayId/agent

- [ ]  create: `POST /create`
- [ ]  remove: `DELETE /remove`
- [ ]  start: `POST /start`
- [ ]  stop: `POST /stop`
- [ ]  get info: `GET /`
- [ ]  edit info: `POST /edit`
- [ ]  update instance: `POST /update-instance`
- [ ]  update config: `POST /update-config`
(internal)
- [x]  push task: `POST /internal/push-task` and `POST /internal/task`
- [ ]  get task status: `GET /internal/task-status`
- [x]  add: `POST /internal/add`
