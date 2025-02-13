package ioc

const (
	_ = iota // Root Object
	Env
	Config
	CoreDB
	CacheCli
	DataStore
	_

	CoreGrpcAuthInterceptor
	CoreGrpcServer
	CoreGatewayService
	CoreGatewayAgentService

	CoreWebServer
)
