package v1

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1/user"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/gin-gonic/gin"
)

var (
	_     User = (*user.User)(nil)
	_user User
)

type User interface {
	Create(c *gin.Context)
}

func NewUser() User {
	if _user == nil {
		_user = &user.User{Db: ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(repo.CoreDbOperator)}
	}
	return _user
}
