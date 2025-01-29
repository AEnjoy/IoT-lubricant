package v1

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core/api/v1/helper"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/global"
	"github.com/AEnjoy/IoT-lubricant/internal/ioc"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/internal/model/request"
	def "github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
)

var (
	_auth *Auth
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
func (a Auth) Signin(c *gin.Context) {
	code, _ := c.GetQuery("code")
	state, _ := c.GetQuery("state")
	token, err := casdoorsdk.GetOAuthToken(code, state)
	if err != nil {
		logger.Errorln(err)
		helper.FailedByServer(err, c)
		return
	}
	c.SetCookie(model.COOKIE_TOKEY_KEY,
		token.AccessToken, token.Expiry.Second(), "/", "",
		false, true)
	u, err := casdoorsdk.ParseJwtToken(token.AccessToken)
	if err != nil {
		helper.FailedByServer(err, c)
		return
	}
	err = a.Db.SaveTokenOauth2(c, token, u.User.Id)
	if err != nil {
		helper.FailedByServer(err, c)
		return
	}
	helper.Success(u, c)
}

var setAuthCrtLock sync.Mutex

func (a Auth) SetAuthCrt(c *gin.Context) {
	if !setAuthCrtLock.TryLock() {
		err := errors.New("do not allow certificates to be set concurrently")
		helper.FailedByClient(err, c)
		return
	}
	defer setAuthCrtLock.Unlock()

	if !global.AllowSetAuthCrt {
		err := errors.New("do not allow certificates to be set")
		helper.FailedByClient(err, c)
		return
	}

	// 从PUT请求中读取文件数据
	file, err := c.FormFile("file")
	if err != nil {
		helper.FailedByClient(err, c)
		return
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		helper.FailedByServer(err, c)
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		helper.FailedByServer(err, c)
		return
	}

	if !isValidCertificate(fileBytes) {
		helper.FailedByClient(err, c)
		return
	}

	dst, err := os.Create(def.AuthCertFilePath)
	if err != nil {
		helper.FailedByServer(err, c)
		return
	}
	defer dst.Close()

	_, err = dst.Write(fileBytes)
	if err != nil {
		helper.FailedByServer(err, c)
		return
	}

	global.AllowSetAuthCrt = false
	helper.Success(nil, c)
}

// 校验文件是否为合法的证书文件
func isValidCertificate(fileBytes []byte) bool {
	block, _ := pem.Decode(fileBytes)
	if block == nil {
		return false
	}

	_, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false
	}
	return true
}

func NewAuth() *Auth {
	if _auth == nil {
		_auth = &Auth{Db: ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(repo.CoreDbOperator)}
	}
	return _auth
}
