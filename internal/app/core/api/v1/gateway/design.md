# Gateways Api Design

**Notice:**

The current design document is still working and will be moved to `/doc` when it is completed

## Apis

root route: /api/v1/gateway

Gateway manage:

1. add gateway host: `POST /add-host`
2. edit gateway host info: `POST /edit-host`
3. remove gateway host: `DELETE /remove-host` -> 8.
4. get gateway host info: `GET /get-host?hostid`
5. list gateway hosts info: `GET /list-hosts?userid`
6. deploy gateway instance: `POST /deploy-instance`
7. update gateway instance: `POST /update-instance`
8. uninstall gateway instance: `POST /uninstall-instance`

Gateway functions:

1. get gateway status: `GET /status?gatewayid`
2. get error logs: `GET /error-logs?gatewayid`
3. set gateway instance config: `POST /set-config?gatewayid`
4. get gateway instance config: `GET /get-config?gatewayid`

Gateway agent instance manage:

root route: /api/v1/gateway/:gatewayId/agent

1. add: `POST /add`
2. create: `POST /create`
3. start: `POST /start`
4. stop: `POST /stop`
5. get info: `GET /`
6. edit info: `POST /edit`
7. update instance: `POST /update-instance`
8. update config: `POST /update-config`
(internal)
9. push task: `POST /push-task`
10. get task status: `GET /task-status`
