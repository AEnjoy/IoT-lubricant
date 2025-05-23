package api

import (
	"github.com/gin-gonic/gin"
)

type IUser interface {
	Create(c *gin.Context)
	GetUserInfo(c *gin.Context)
}
type IGateway interface {
	AddHost(c *gin.Context)
	ListHosts(c *gin.Context)
	DescriptionHost(c *gin.Context)
	DescriptionGateway(c *gin.Context)
	EditGateway(c *gin.Context)
	AddGateway(c *gin.Context)
	RemoveGateway(c *gin.Context)

	AgentPushTask(c *gin.Context)
	AddAgentInternal(c *gin.Context)
}
type IAuth interface {
	Login(c *gin.Context)
	Signin(c *gin.Context)
	SetAuthCrt(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetPrivateKey(c *gin.Context)
}
type IMonitor interface {
	// BaseInfo 返回网关个数，agent个数，离线个数信息，node个数信息
	BaseInfo(c *gin.Context)
}
type IAgent interface {
	// Operator : start-agent,stop-agent,start-gather,stop-agent,get-openapidoc
	Operator(*gin.Context)
	// Set : set agent gather config
	Set(*gin.Context)
	GetData(*gin.Context) //todo
	GetAgentInfo(c *gin.Context)
	List(*gin.Context)
}
type ITask interface {
	QueryTask(*gin.Context)
	GetTaskList(*gin.Context)
}
type ILog interface {
	GetLogList(*gin.Context)
}
type IProject interface {
	AddProject(*gin.Context)
	RemoveProject(*gin.Context)
	AddDataStoreEngine(*gin.Context)
	GetProjectDataStoreEngineStatus(*gin.Context)
	UpdateEngineInfo(*gin.Context)
	AgentBindProject(*gin.Context)

	AddWasher(*gin.Context)
	AgentBindWasher(*gin.Context)
}
