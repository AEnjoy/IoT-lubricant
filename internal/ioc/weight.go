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
	CoreGrpcAuthInterceptor
	CoreGrpcServer
	CoreGatewayService
	CoreGatewayAgentService

	CoreWebServer
)
