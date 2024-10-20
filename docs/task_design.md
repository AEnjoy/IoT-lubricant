# Lubricant Task Design:

This article is used to describe and specify the `task list` for communication between `Core` and `Gateway`

**Note**:

1. All codes in the design document are considered as **pseudocodes** or codes that **cannot be used directly**
2. In the following, if there is no special reference, `Core` is considered as `Lubricant Core`

[TOC]

## Design goals:

1. Need a **Definition** of **Task data structure**
2. Need to define the `operation type` as a `constant`
3. Need to understand the network structure of the deployed target device, containerized runtime (Docker), and daemon process information
4. Reserve space for future improvement in containerization capabilities
5. Need to be able to transfer API document structures and configure parameters required for containerized deployment of Docker

## Design none-goals:

None

## Data structure:

- Task data structure
  - 
  - TaskID
  - Operator (Role:User/Core/Gateway?,ID:UUID?)
  - Task Executor (Role:User/Core/Gateway?)
  - Executor ID   (ID:UUID)
  -
  - Operation type (const enum  it will be defined)
  - Operation command ([]byte a json object)
  - Rollback capability (bool with rollback command)
  -
  - Operation time (timestamp write to DB)

- Operator: 
  - 
  - Role (Visitor/User/Core/Gateway/Agent/Schedule) enum type
  - Model:pkg/model/role.go
  - int const from 0 to 5
  - Visitor: no permission
  - Prefix: ROLE_

- Executor: Role (User/Core/Gateway/Agent) enum type
  -
  - Model:pkg/model/role.go
  - int const from 1 to 4
  - Prefix: ROLE_

- Operation type: 
  - 
  - Operation enum type
  - Visitor: 0 
  - User: Login:10; Logout:11; ChangePassword(-); CreateTask:13; QueryTask:14; ViewTaskResult:15;
  - Core: AddGateway:20; RemoveGateway:21; AddAgent:22; RemoveAgent:23; AddSchedule:24; RemoveSchedule:25;
  - Gateway: AddDriverContainer:30; RemoveDriverContainer:31; AddAgentContainer:32; RemoveAgentContainer:33;
  - Agent: EnableOpenAPI:40; DisableOpenAPI:41;SendRequest:42;GetOpenAPIDoc:43;GetEnableOpenAPI:44;

- AddGateway
  
  ```json
  {
    "way":0, 
    "config": hostInfo/preConfig,
  }
  ```
  
  Config:
  
  hostInfo:
  
  ```json
  {
    "host": "example.client.com",
    "port": 22,
    "user": "root",
    "pwd": "123456"
  }
  ```
  
  AddWay: 1 (IP(HOST):PORT+SSH User:PWD)   2 (via pre config)
  
  preConfig:
  ```json
  {
    "id": "gatewayID",
    "host": "example.server.com",
    "port": 9090,
    "tls": true,
    "tlc_config": {
      "enable": true,
      "skip_verify": false,
      "from_file": false,
      "key": "",
      "cert": "",
      "ca": ""
    }
  }
  ```

- RemoveGateway

  ```json
  {
    "id": "gatewayID",
    "removeAgent": true
  }
  ```
- AddAgent

  Because the `Agent` is a `calling agent` for Driver OpenAPI, when deploying the Agent, the Agent needs to know `the IP 
  address of the Driver-Container`, `the OpenAPIDoc of the service`, and `the complete OpenAPI parameters that are enabled`
  
  This makes the task of AddAgent divided into several parts: `Driver-Container-DeployInfo`,`Driver-Container-Info`,`OpenAPI-Doc`,
  and `OpenAPI-Enable`.
  
  Driver-Container-DeployInfo: Docker Container Deploy Info(config)
  
  Driver-Container-Info: Docker Container Info(id,ip)
  
  OpenAPI-Doc: an OpenAPI document object
  
  OpenAPI-Enable(Optional because it can enable after the Agent is deployed): an array of OpenAPI parameters that are enabled
  
  ```json
  {
    "container":Driver-Container-DeployInfo,
    "containerInfo":Driver-Container-Info,
    "openapiDoc":OpenAPI-Doc,
    "openapiEnable":OpenAPI-Enable
  }
  ```

- RemoveAgent:
  ```json
  {
    "id":"agentID",
    "removeDriver":true
  }
  ```

- AddSchedule

  Add a scheduled task: send requests to an API regularly and collect results
  
  ```json
  {
    "id":"agentID",
    "interval":60,
    "request":{},
    "api": "api_name"
  }
  ```
  
  Return an error if the API does not exist,or return a scheduleID.

- RemoveSchedule:
  ```json
  {
    "id":"scheduleID"
  }
  ```

- AddDriverContainer,RemoveDriverContainer,AddAgentContainer,and RemoveAgentContainer:

  Its implementation has been described in `AddAgent`

- EnableOpenAPI:

  Enable a configuration culture for data collection
  
  ```json
  {
    "id":"agentID",
    "openapi_enable": openapi_doc_object
  }
  ```
  
- DisableOpenAPI:

  ```json
  {
    "id":"agentID",
    "openapi_disable": [] // paths
  }
  ```

- SendRequest:

  ```json
  {
    "id":"agentID",
    "method": "GET",
    "path": "",
    "request":{}
  }
  ```

- GetOpenAPIDoc:

  ```json
  {
    "id":"agentID"
  }
  ```
  
- GetEnableOpenAPI:

  ```json
  {
    "id":"agentID"
  }
  ```
  

