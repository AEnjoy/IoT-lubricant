package ioc

const (
	_ = iota // Root Object
	Env
	Config
	CoreDB
	CacheCli
	_
	CoreGrpcAuthInterceptor
	CoreGrpcServer
)
