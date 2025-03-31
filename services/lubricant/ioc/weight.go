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
	CoreGrpcAuthInterceptor
	CoreGrpcServer
	CoreSyncTaskSystem
	CoreGatewayService
	CoreGatewayAgentService

	GatewayStatusGuard
	BackendHandlerReport
	BackendHandlerErrLogs
	BackendHandlerDataUpload

	CoreWebServer
)
