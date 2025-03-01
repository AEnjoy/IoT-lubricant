package v1

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	def "github.com/aenjoy/iot-lubricant/pkg/default"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	"github.com/aenjoy/iot-lubricant/pkg/model/response"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/lubricant/api/v1/helper"
	"github.com/aenjoy/iot-lubricant/services/lubricant/config"
	"github.com/aenjoy/iot-lubricant/services/lubricant/global"
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
)

var (
	_auth *Auth
)

type Auth struct {
	Db repo.ICoreDb
}

func (a Auth) Register(c *gin.Context) {

}
func (a Auth) Login(c *gin.Context) {
	conf := config.GetConfig()
	req := helper.LoginRequest2CasdoorLoginRequest(helper.RequestBind[request.LoginRequest](c))
	if req == nil {
		return
	}

	client := http.Client{}
	u, _ := url.Parse(fmt.Sprintf("%s/api/login", conf.AuthEndpoint))
	params := u.Query()
	params.Add("clientId", conf.AuthClientID)
	params.Add("responseType", "code")
	params.Add("redirectUri", fmt.Sprintf("http://%s/api/v1/signin"))
	params.Add("type", "code")
	params.Add("scope", "read")
	params.Add("state", "casdoor")
	u.RawQuery = params.Encode()
	marshal, err := sonic.Marshal(req)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError, exception.ErrNewException(err, exceptionCode.ErrorEncodeJSON), c)
		return
	}
	resp, err := client.Post(u.String(), "application/json", bytes.NewReader(marshal))
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError, exception.ErrNewException(err, exceptionCode.ErrorCommunicationWithAuthServer), c)
		return
	}
	var loginResult response.CasdoorLoginResponse
	err = decoder.NewStreamDecoder(resp.Body).Decode(&loginResult)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError, exception.ErrNewException(err, exceptionCode.ErrorDecodeJSON), c)
		return
	}

	if loginResult.Status != "ok" {
		helper.FailedWithJson(http.StatusUnauthorized, exception.ErrNewException(err, exceptionCode.ErrorCommunicationWithAuthServer), c)
		return
	}
	signin(loginResult.Data, "casdoor", c, a.Db)
}
func (a Auth) Signin(c *gin.Context) {
	code, _ := c.GetQuery("code")
	state, _ := c.GetQuery("state")
	signin(code, state, c, a.Db)
}

func signin(code, state string, c *gin.Context, db repo.ICoreDb) {
	token, err := casdoorsdk.GetOAuthToken(code, state)
	if err != nil {
		logger.Errorln(err.Error())
		helper.FailedWithJson(http.StatusUnauthorized, exception.ErrNewException(
			err, exceptionCode.ErrorInvalidAuthKey), c)
		return
	}
	c.SetCookie(model.COOKIE_TOKEY_KEY,
		token.AccessToken, int(token.Expiry.Unix()-time.Now().Unix()), "/", "",
		//token.AccessToken, token.Expiry.Second(), "/", "",
		//token.AccessToken, int(24*time.Hour), "/", "",
		false, true)
	u, err := casdoorsdk.ParseJwtToken(token.AccessToken)
	if err != nil {
		logger.Errorf("parse token error: %v", err)
		helper.FailedByServer(err, c)
		return
	}
	err = db.SaveTokenOauth2(c, token, u.User.Id)
	if err != nil {
		logger.Errorf("save token error: %v", err)
		helper.FailedByServer(err, c)
		return
	}
	helper.SuccessJson(u, c)
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
	helper.SuccessJson(nil, c)
}

// 校验文件是否为合法的证书文件
func isValidCertificate(fileBytes []byte) bool {
	block, _ := pem.Decode(fileBytes)
	if block == nil {
		return false
	}

	_, err := x509.ParseCertificate(block.Bytes)
	return err == nil
}

func NewAuth() *Auth {
	if _auth == nil {
		_auth = &Auth{Db: ioc.Controller.Get(ioc.APP_NAME_CORE_DATABASE).(repo.ICoreDb)}
	}
	return _auth
}
