package v1

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1/helper"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/internal/model/request"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	Db repo.CoreDbOperator
}

func (a Auth) Register(c *gin.Context) {

}
func (a Auth) Login(c *gin.Context) {
	req := request.LoginRequest{}
	if err := c.ShouldBind(&req); err != nil {
		helper.FailedByClient(err, c)
		return
	}
	if user, err := a.Db.QueryUser(c, req.UserName, ""); err != nil {
		helper.FailedByServer(err, c)
		return
	} else {
		err := user.CheckPassword(req.Password)
		if err != nil {
			helper.FailedByClient(err, c)
			return
		}
		tk := model.NewToken(&user)
		err = a.Db.SaveToken(c, tk)
		if err != nil {
			helper.FailedByServer(err, c)
			return
		}
		c.SetCookie(model.COOKIE_TOKEY_KEY,
			tk.AccessToken,
			tk.RefreshTokenExpiredAt,
			"/", "",
			false, true)
		helper.Success(user, c)
	}
}
