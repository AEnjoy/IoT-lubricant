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

	AddGatewayInternal(c *gin.Context)
	RemoveGatewayInternal(c *gin.Context)

	AgentPushTask(c *gin.Context)
	AddAgentInternal(c *gin.Context)
}
type IAuth interface {
	Login(c *gin.Context)
	Signin(c *gin.Context)
	SetAuthCrt(c *gin.Context)
	RefreshToken(c *gin.Context)
}
type IMonitor interface {
	// BaseInfo 返回网关个数，agent个数，离线个数信息，node个数信息
	BaseInfo(c *gin.Context)
}
