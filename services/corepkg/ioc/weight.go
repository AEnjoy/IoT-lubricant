package ioc

const (
	_ = iota // Root Object
	Env
	Config
	CoreDB
	CacheCli
	DataStore
	_

	CoreMqService
	SvcLoggerService
	SvcDataStoreApiService
	CoreGrpcAuthInterceptor
	CoreGrpcServer
	CoreSyncTaskSystem
	CoreGatewayService
	CoreGatewayAgentService
	CoreProjectService

	GatewayStatusGuard
	BackendHandlerReport
	BackendHandlerErrLogs
	BackendHandlerDataUpload

	CoreWebServer
)
