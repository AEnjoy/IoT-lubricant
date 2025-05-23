# Lubricant Core Backend Design
Lubricant Core Backend Component Function Design

**Note**: 

1. All codes in the design document are considered as **pseudocodes** or codes that **cannot be used directly**
2. In the following, if there is no special reference, `Core` is considered as `Lubricant Core`

[TOC]

## General Design:

### High Availability Requirements:

1. Ensure that the `Core service` is guarded by `daemons` or `services`, or `container pools`, and restarted when it crashes to achieve high availability
   1. It can be achieved by using `systemd`, `supervisor` or `Core-deamon-self` to restart the service when it crashes
   2. Multiple pods created by `Kubernetes' deployment` `Docker-compose` or `Kubernetes' StatefulSet`, container crash automatic restart policy to achieve stable operation of multiple replicas
2. After initialization, use `defer recover()` to ensure the stability of services panic
3. Code quality:
   1. All `non placeholder errors` should **be handled**
   2. Unit tests should try their best to achieve **path coverage**
   3. All `pointers` and `arrays` should be **judged whether they can be accessed**
   4. All commits should pass UnitTest and E2ETest
   5. By pass golang-lint to ensure the code quality
4. `Logger` should be used to record all errors logs, and report to the `log collector` and developers

### Performance Requirements:

1. In addition to dynamic data type objects, other types of objects should be **initialized when the program starts**.
2. In addition to the basic data type and configuration, other composite data types are **passed through pointers or references**
3. In addition to the `any or interface{}` type parameter specified in the immutable function signature definition, we should minimize the use of `any` and **use the precompiled data `type-generic`**
4. The `Core` should reduce the direct modification and access to data, and **reduce the `serialization/deserialization`** of data
5. **Non-blocking** ,**Multi-threading** and **Asynchronous** processing of each task
   1. When the task is obtained, the response ID(NextID) of the task is **returned directly**
   2. Add tasks to the **task pool** for unified processing
   3. Use `goroutine`, `select` and `context` to process tasks in parallel
6. Use **high-performance libraries** such as [gin](https://github.com/gin-gonic/gin) and [sonic](https://github.com/bytedance/sonic)
7. Reduce output to `stdout`, and only output **necessary information** and **error information**

### High Scalability Requirements:

1. Allow for **dynamic expansion** of `Core` and `Gateway` services through `scripts`
2. The `Core` should be able to **load and unload** `scripts` dynamically
3. Use `etcd` or `consul` as the `service discovery` and `service registration` center
4. (Not safe)Use library `github.com/dop251/goja`to implement the `JavaScript` execution engine
5. (Not safe)Allows the user to run `SQL` scripts on their own data
6. (Not safe)Allows the user to manually configure the script executor to execute custom code

### Safe and Robustness Requirements:

1. All connection requests need to be authenticated
2. Tls should be used to ensure data transmission security

### Containerization and Cloud Native Affinity Requirements:

1. Use `Docker` to package the `Core` services
2. Support use `Docker` or `Kubernetes` to manage the `Core` services

### Highly automated Requirements:

1. Use `Makefile` to automate the compilation and packaging process of the `Core` service
2. Ensure that the unit test works properly
3. Use `GitHub Actions` to automate the testing and deployment process of the `Core` service

### API Design:

1. All internally implemented apis should not be exposed to external callers
2. The general api should be encapsulated and stored in `api/<version>/` for users to call
3. Internal services and logics should not call public apis, which violates the circular dependency principle
4. The path api/ should be used as a separate gomod
5. In addition to the `struct` necessary for `data definition`, others should be encapsulated as `interfaces`

## Web Server And APIs

- Lubricant-Core will use `go-gin` as web server and `gin-swagger` for API documentation.
- The Web Server Component named `ioc.WebServer` will be registered and managed by `IoC`.
- APIs will be divided into two parts: `Publicly available operational API interfaces` and `Publicly available generic interfaces`
  - The former will be used for `direct operations` on the Core, while the latter will be used to provide `static resources`
- Middleware will be added to the web server to handle common tasks such as request logging, authentication, cache, etc.


### Module Design:
```mermaid
graph TD;
    A[Developers complete the development of the frontend-page] --> B[yarn build command to edge the frontend-page];
    B --> C[Embed static page into the application through go:embed];

    D[Web Server Components Call Order] --> E[Lubricant-Core];
    E --> F[APPInit];
    F --> G[Initialize IoC object];
    G --> H[Register Web Server Components into IoC];
    H --> I[Call IoC Init to initialize all objects];
    I --> J[Initialize the Root Object and start the web server];

    K[Simultaneously Registered Objects] --> L[Middleware: RequestLogging, Authentication, Cache];
    K --> M[gin.Engine Root Object];
    K --> N[Static resources];
    K --> O[Main routes];
    K --> P[Routing groups];
    K --> Q[Sub routes];

    R[Initialization Order] --> S[Root Object];
    S --> T[Static resources];
    T --> U[Main routes];
    U --> V[Middleware];
    V --> W[Routing groups];
    W --> X[Sub routes];
```

### Library Dependencies:

- `go-gin` as web server
```go
import "github.com/gin-gonic/gin"
```
- `gin-swagger` for API documentation

### IoC hosting:

The Web Server Component will be registered and managed by `IoC`.

Including `root objects`, `static resources`, `main routes`, `routing groups`, `sub routes`, etc... All should be hosted by `IoC`

- Name: `ioc.WebServer`  
  ```go
  const APP_NAME_CORE_WEB_SERVER = "lubricant-core-web-server"
  
  const CoreWebServer = iota
  ```
- Weight: Any but higher than `ioc.DataStore`
  - `root objects`: Equivalent to ioc.WebServer.Weight
  - `static resources`: Equivalent or greater than ioc.WebServer.Weight
  - `main routes`: Any weight greater than ioc.WebServer.Weight
  - `routing groups`: Any weight greater than `main routes`
  - `sub routes`: Any weight greater than `routing groups`
- Naming rules:

For const string: `const APP_NAME_CORE_WEB_SERVER_<COMPONENT_NAME> = "lubricant-core-web-server_<component_name>"`

For Weight: `const CoreWebServer<component_name>`

### Middleware:

Only `Publicly available operational API interfaces` will be handled by middleware.

- Request Logging: Use `gin.Logger()` middleware to log all requests.
- Authentication: Use `pkg/middleware/auth` middleware to authenticate requests.
- Cache: Use `pkg/middleware/cache` middleware to cache responses.

### Embedded resources:
Use Golang `embed` feature to embed `static resources` into the binary.

Static resources' path: `frontend/dist`

Build frontend:
```shell
yarn install
yarn run build
```

Embed static resources to the binary:

```go
package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist
var f embed.FS

func main() {
  r := gin.Default()
  st, _ := fs.Sub(f, "frontend/dist")
  r.StaticFS("/", http.FS(st))

  err := r.Run("0.0.0.0:5999")
  if err != nil {
    fmt.Println(err)
  }
}
```

### Test:

Ut and e2e test will be added to ensure the correctness of the API.

### API Documentation:

Todo: Need to be designed

- Unit test: 
  - Use `go test` to test the API.
- End-to-end test:
1. Get The Test authorization token
2. Send requests to the API

## gRPC Server

- gRPC server will be used to provide `remote procedure calls` between `Core` , `Gateway Client` and `User Client`.
- The gRPC server Component named `ioc.GRPCServer` will be registered and managed by `IoC`.
- The gRPC Server will provide `Gateway oriented` and `User oriented` services
- TLS(Optional) will be used for encryption, and mutual authentication is required
- For detailed design, please refer to [gRPC_service_design](gRPC_service_design.md)

### Module Design:

The gRPC Server is used to process requests from one or more gRPC clients

When the gRPC client attempts to establish a connection with the server, it will be required to perform Tls connection authentication. After the TLS authentication is passed, the client ID will be verified, and remote calls will be allowed only after the all authentication is passed.

The client will establish a service with the server, including 3 `bidirectional flow` remote call APIs and a one-way flow call API

They are:

bidirectional flow:
ping:Used to realize connectivity detection
getTask:It is used to enable the client to obtain the task from the server, and the server pushes the task to the client
pushData:It is used to exchange data between client and server

one-way flow:
pushMessageId:It is used to push messages from the client to the server, and get the message content back

When the client calls the `pushData` and `pushMessageId` API, the server will directly create an external thread for data processing according to the specific data/message content type, and immediately return the data processing receipt.
After receiving the receipt, the client can call `getTask` or `pushMessageId` by carrying the receipt, and parse the running result of the task according to the specific fields in the return

```mermaid
graph TD;
    A[gRPC Server] -->|Processes requests from| B[gRPC Clients]
    B -->|Establishes connection| C[TLS Authentication]
    C -->|Verifies client ID| D[Authentication Passed]
    D -->|Allows remote calls| E[Service Established]
    
    E --> F[Bidirectional Flow APIs]
    E --> G[One-Way Flow API]

    F --> H[ping]
    H -->|Connectivity detection| I[Realized]

    F --> J[getTask]
    J -->|Obtain task from server| K[Task pushed to client]

    F --> L[pushData]
    L -->|Exchange data| M[Between client and server]

    G --> N[pushMessageId]
    N -->|Push messages| O[Client to server]
    N -->|Get message content back| P[Content received]

    L -->|Creates external thread| Q[Data Processing]
    Q -->|Returns receipt| R[Data Processing Receipt]
    R -->|Client calls with receipt| S[getTask or pushMessageId]
    S -->|Parse running result| T[Based on return fields]
```

### Library Dependencies:
- `protobuf` for message and service definition
    ```shell
      go install google.golang.org/protobuf/cmd/protoc-gen-go
      go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
      go get google.golang.org/protobuf
    ```
- `grpc` for gRPC server
    ```shell
      go get google.golang.org/grpc
    ```

Build the protobuf files:

```shell
protoc -I=. --go_out=.   --go-grpc_out=. --go-grpc_opt=module= protobuf/core/data.proto
```

### IoC hosting:

The gRPC Server Component will be registered and managed by `IoC`.

Due to the design of gRPC, the complete logic of gRPC will be treated as a large object

The small objects under it are only `authorized` objects

- Name: `ioc.GRPCServer`
  ```go
  const APP_NAME_CORE_GRPC_SERVER = "lubricant-core-grpc-server"
  
  const CoreGrpcServer = iota
  ```
- Weight: Any but higher than `ioc.DataStore`
  - `Authorized objects`: Less than ioc.WebServer.Weight
- Naming rules:
  - `const APP_NAME_CORE_GRPC_SERVER_<SERVICE_NAME> = "lubricant-core-grpc-server_<service_name>"`
  - `const CoreGrpcServer<service_name>`

### Data Model:

**Gateway oriented:**
```protobuf
syntax = "proto3";

service coreService {
  rpc ping(stream Ping) returns (stream Ping) {}
  rpc getTask(stream Task) returns (stream Task) {}
  rpc pushMessageId(MessageIdInfo) returns (MessageIdInfo) {}
  rpc pushData(stream Data) returns (stream Data) {}
}

message Ping {
  int32 flag = 1;
}
message Task {
  string taskId = 1;
  string messageId = 2;
  bytes content = 3;
}

message MessageIdInfo {
  string messageId = 1;
  string gatewayId = 2;
  string time = 3;
  string emergency = 4;
  bytes content = 5;
}

message Data {
  string messageId = 1;
  string gatewayId = 2;
  string agentID = 3;
  repeated bytes data = 4;
  int32 dataLen = 5;
  string time = 6;
  int32 cycle = 7;
}
```

**User oriented:**

```protobuf
syntax = "proto3";

service userService {

}

message messageName{
  
}
```

### Test:

It is a bit difficult to test gRPC, but we should try our best to make the software robust

In terms of unit testing, use `mockery` to generate mock objects for gRPC server

In terms of end-to-end testing, we can use the `grpcurl` or other tool to send requests to the gRPC server and check the results.

## DataStore

DataStore is a component that stores data in a persistent way. 

As we all know, in the case of massive data, Disk I/O is an important factor restricting performance. Therefore, we should not frequently access Disks, but speed up our requests through `Caching Technology`

- The DataStore consists of two parts: a relational database for persistent storage(MySQL) and a cache for accelerating data access(Redis Optional)
- **Data cleaning**: Since the data collected from the original sensor is only compressed and encoded, we need to decompress and decode the processed data
  - Data compression is **optional** and depends on the configuration. Support alogs: `gzip`, `lz4`, `zstd` and `not compressed`
  - **Manual script support**: allows secondary processing of uploaded data by manually writing **JavaScript**, **SQL** and **other scripts** by the user(Not Safe)
- Support forwarding: You can customize the external data forwarding interface by configuring external APIs, cloud function APIs, etc.(Safe); Or directly forward to other applications according to the specific protocol(Not Safe)
- Data storage: store the cleaned(optional) data in the database and cache(optional)

### Database:

Now, we use MySQL as the default database, and will support more databases in the future:

- Transaction support should be necessary
- We should encapsulate all request SQL into a method in the `struct` and implement it as an instance of the `interface`
  - It will be easier to maintain, extend and **mock for testing**
- All model structs should be defined in the `pkg/model` package
- Use `gorm` to implement the database operation

Database model:

```go
var _ ioc.Object = (*CoreDb)(nil) // IoC
var _ CoreDbCli = (*CoreDb)(nil) // CoreDbCli interface

type CoreDbCli interface {
 // apis
}
type CoreDb struct {
	db *gorm.DB
}
// apis
```

#### Interface design:

We will use the **Object-Oriented Design** to design the functional interface of the database client

The `CoreDbCli` will support the following features:
- 
- Transaction support: `BeginTx`, `Commit`, `Rollback`
  - In function signatures, the txn parameter is optional. If the txn parameter is not passed, the database client will not use transaction.
- Context support: 
  - All command will be executed with a context, and the context will be passed to the database client by `WithContext`
- Mock for testing: Support mocking the database client through `mockery`

We will support the following methods:
- 
- Transaction: `BeginTx`, `Commit`, `Rollback`
- For Gateway:
  - 
  - `IsGatewayIdExists` id -> bool
  - `GetGatewayInfo` GatewayID-> Gateway and error
  - `AddGateway` userID,gateway -> error
  - `UpdateGateway` gateway -> error // Set
  - `DeleteGateway` GatewayID -> error
  - `GetGatewayList` -> []Gateway and error
  - `GetGatewayTask` GatewayID -> []Task and error
  - `AddGatewayTask` GatewayID,task -> error 
- For Agent:
  - 
  - `IsAgentIdExists` id -> bool
  - `AddAgent` GatewayID,agent -> error
  - `UpdateAgent` AgentID,agent -> error  // Set
  - `DeleteAgent` AgentID -> error
  - `GetAgentList` GatewayID -> []Agent and error
  - `AddAgentTask` AgentID,task -> error
  - `GetAgentInfo` AgentID -> Agent and error
  - `SetAgentGather` AgentID,config -> error // API Path, Algorithm
- For Data:
  - 
  - `GetDataCleaner` AgentID -> Cleaner and error
  - `StoreAgentGatherData` AgentID,data -> error
  - `DeleteAgentGatherData` AgentID -> error
- For User:
  - 
  - `AddUser` user -> error
  - `UpdateUser` user -> error
  - `DeleteUser` userID -> error
  - `VerifyUser` userID,password -> bool

### Cache (Optional):

- We will use Redis as the default cache, and will support more cache systems in the future
- Cache client will be encapsulated as an interface to facilitate calling, component replacement and unit testing
- The `cache client interface` will be **forced to be enabled as a component**, but whether to use `NilCache`, `Redis` or other memory databases is up to the user

Cache client interface:

All implementations need to meet the following interface, including the **empty cache type**

```go
type CacheCli[T any] interface {
	SetEx(ctx context.Context, key string, value T, duration time.Duration) error
	Set(ctx context.Context, key string, value T) error
	HSet(ctx context.Context, key string, field string, value T) error
	HGet(ctx context.Context, key string, field string) (T, error)
	Get(ctx context.Context, key string) (T, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	Delete(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, duration time.Duration) error
	Close(ctx context.Context) error
}
```

If its instantiation implementation is an empty cache type, you need to ensure that the error returned by each method is `cache.ErrNullCache`

```go
var ErrNullCache = errors.New("cache client is nil")
```

In addition, the `cache client instantiation` should also implement `ioc.Object` interface

### IoC Hosting:

The `DataStore` component will be registered and managed by `IoC`.

- Name: `ioc.DataStore` 
    ```go
        const APP_NAME_CORE_DATABASE = "lubricant-core-database-core" // mysql
        const APP_NAME_CORE_CACHE = "lubricant-core-cache-core" // redis(cache)
  
        const APP_NAME_CORE_DATABASE_STORE = "lubricant-core-datastore"
    ```
  
- Weight: Any but less than `ioc.GRPCServer` and other services components
  - APP_NAME_CORE_DATABASE_STORE: Any but higher than `APP_NAME_CORE_DATABASE` and `APP_NAME_CORE_CACHE`



## Configuration:

- As the first important component besides the command line parameters, the configuration component should be loaded at startup and be the earliest initialized component in IoC
- It should be initialized by loading configuration files or environment variables
```go
var _ ioc.Object = (*Config)(nil)

type Config struct {
	// app
	AppVersion string
	TlsEnable  bool      `yaml:"tls" env:"TLS_ENABLE" envDefault:"false"`
	HTTPTls    bool      `yaml:"tls_http" env:"HTTP_TLS_ENABLE" envDefault:"false"`
	GRPCTls    bool      `yaml:"tls_grpc" env:"GRPC_TLS_ENABLE" envDefault:"false"`
	Tls        model.Tls `yaml:"tls_config" env:"TLS_CONFIG" envPrefix:"TLS_"`

	// grpc
	GrpcPort int `yaml:"port" env:"GRPC_LISTEN_PORT" envDefault:"9090"`

	// web
	Host    string `yaml:"host" env:"HTTP_LISTEN_HOST" envDefault:"0.0.0.0"`
	WebPort int    `yaml:"port" env:"HTTP_LISTEN_PORT" envDefault:"8080"`
	Domain  string `yaml:"domain" env:"HOSTNAME" envDefault:"localhost"`

	// mysql
	MySQLHost     string `yaml:"host" env:"DATASOURCE_HOST,required"`
	MySQLPort     int    `yaml:"port" env:"DATASOURCE_PORT,required"`
	MySQLDB       string `yaml:"database" env:"DATASOURCE_DB,required"`
	MySQLUsername string `yaml:"username" env:"DATASOURCE_USERNAME,required"`
	MySQLPassword string `yaml:"password" env:"DATASOURCE_PASSWORD,required"`
	MySQLDebug    bool   `yaml:"debug" env:"DATASOURCE_DEBUG" envDefault:"false"`

	// redis
	RedisEnable   bool   `yaml:"enable" env:"REDIS_ENABLE" envDefault:"false"`
	RedisHost     string `yaml:"host" env:"REDIS_HOST"`
	RedisPort     int    `yaml:"port" env:"REDIS_PORT"`
	RedisPassword string `yaml:"password" env:"REDIS_PASSWORD"`
	RedisDB       int    `yaml:"db" env:"REDIS_DB"`
}
```

### IoC hosting:

- Name: `ioc.Config`
  ```go
    const APP_NAME_CONFIG = "lubricant-core-config"
  ```
  
- Weight: Any but less than other all services components

## Initialization:

The core has two modes of operation. 

- One is as a daemon. At this time, it does not need to initialize other components except configuration
- The other is the service mode. At this time, it needs to call `IoC.Init()` to initialize all components

