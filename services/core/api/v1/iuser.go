package v1

import (
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/services/core/api/v1/user"
	"github.com/AEnjoy/IoT-lubricant/services/core/models"
	"github.com/gin-gonic/gin"
)

var (
	_     IUser = (*user.Api)(nil)
	_user IUser
)

type IUser interface {
	Create(c *gin.Context)
	GetUserInfo(c *gin.Context)
}

func NewUser() IUser {
	if _user == nil {
		_user = &user.Api{Db: ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(models.ICoreDb)}
	}
	return _user
}
