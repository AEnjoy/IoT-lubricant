# Gateways Api Design

**Notice:**

The current design document is still working and will be moved to `/doc` when it is completed

## Apis

root route: /api/v1/gateway

Gateway manage:

- [x]  add gateway host: `POST /add-host`
- [x]  edit gateway host info: `POST /edit-host`
- [ ]  remove gateway host: `DELETE /remove-host` -> uninstall.
- [x]  get gateway host info: `GET /get-host?hostid`
- [x]  list gateway hosts info: `GET /list-hosts?userid`
- [x]  deploy gateway instance: `POST /deploy-instance`
- [ ]  update gateway instance: `POST /update-instance`
- [ ]  uninstall gateway instance: `POST /uninstall-instance`
(internal)
- [ ]  add gateway: `POST /add-gateway`
- [ ]  remove gateway: `POST /remove-gateway`

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
- [ ]  push task: `POST /push-task`
- [ ]  get task status: `GET /task-status`
- [ ]  add: `POST /add`
