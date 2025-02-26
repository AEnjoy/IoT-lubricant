package v1

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/user"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
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
		_user = &user.Api{Db: ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(repo.ICoreDb)}
	}
	return _user
}
