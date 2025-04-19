package helper

import (
	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	"github.com/aenjoy/iot-lubricant/services/corepkg/config"
)

func LoginRequest2CasdoorLoginRequest(req *request.LoginRequest) *request.CasdoorPasswordAuthRequest {
	c := config.GetConfig()
	if req == nil {
		return nil
	}
	return &request.CasdoorPasswordAuthRequest{
		Application:  c.AuthApplicationName,
		Organization: c.AuthOrganization,
		Username:     req.UserName,
		AutoSignin:   true,
		Password:     req.Password,
		SigninMethod: "Password",
		Type:         "code",
	}
}
