package model

import (
	"bytes"
	"errors"
	"os"
	"os/exec"

	"github.com/dop251/goja"
)

type Server struct {
	Info map[string]string `json:"server_info" gorm:"column:server_info;serializer:json"`
}

func (Server) TableName() string {
	return "server"
}

type User struct {
	ID       int    `json:"id" gorm:"column:id"`
	UserId   string `json:"user_id" gorm:"column:user_id"` // uuid
	UserName string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`

	CreatedAt int64 `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at"`
}

func (req User) CheckPassword(password string) error {
	// todo: encrypt password
	if password == req.Password {
		return nil
	}
	return errors.New("password error")
}
func (User) TableName() string {
	return "user"
}

type Data struct {
	ID      int    `json:"id" gorm:"column:id"`
	AgentID string `json:"agent_id" gorm:"column:agent_id"` // equal to  DeviceID

	Content string `json:"data" gorm:"column:data;serializer:json"` // core.Data 序列化的json

	CreatedAt int64 `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at"`
}

func (Data) TableName() string {
	return "data"
}

type Gateway struct {
	GatewayID   string `json:"id" gorm:"column:id"`
	UserId      string `json:"user_id" gorm:"column:user_id"`
	Description string `json:"description" gorm:"column:description"`

	Address           string `json:"address" gorm:"column:address"`                         // SSH: ip:port or domain:port
	UserNameAndPasswd string `json:"username_and_passwd" gorm:"column:username_and_passwd"` //

	TlsConfig string `json:"tls_config" gorm:"column:tls_config;serializer:json"` // grpc tls config

	CreatedAt int64 `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at"`
}

func (Gateway) TableName() string {
	return "gateway"
}

type Clean struct {
	ID          int    `json:"id" gorm:"column:id"`
	AgentID     string `json:"agent_id" gorm:"column:agent_id"`
	Description string `json:"description" gorm:"column:description"`

	Interpreter string   `json:"interpreter" gorm:"column:interpreter"` // python,goja,node,bash or other
	Script      string   `json:"script" gorm:"column:script"`           // 脚本代码
	Command     []string `json:"command" gorm:"column:command"`         // 提供给解释器的额外参数

	CreatedAt int64 `json:"-" gorm:"column:created_at"`
	UpdatedAt int64 `json:"-" gorm:"column:updated_at"`
}

func (Clean) TableName() string {
	return "clean"
}

var rt *goja.Runtime

func (c *Clean) Run(data []byte) ([]byte, error) {
	switch c.Interpreter {
	case "":
		return data, nil
	case "goja":
		if rt == nil {
			rt = goja.New()
		}

		_, err := rt.RunString(c.Script)
		if err != nil {
			return data, err
		}

		processData, ok := goja.AssertFunction(rt.Get(c.Command[0]))
		if !ok {
			return data, errors.New("not a function")
		}

		result, err := processData(goja.Undefined(), rt.ToValue(data))
		if err != nil {
			return data, err
		}

		return []byte(result.String()), nil
	default:
		err := os.WriteFile("script", []byte(c.Script), 0666)
		if err != nil {
			return data, err
		}
		defer func() {
			_ = os.Remove("script")
		}()

		var newCommand []string
		newCommand = append(newCommand, "script")
		newCommand = append(newCommand, c.Command...)
		cmd := exec.Command(c.Interpreter, newCommand...)
		cmd.Stdin = bytes.NewReader(data)

		var out bytes.Buffer
		cmd.Stdout = &out

		err = cmd.Run()
		if err != nil {
			return data, err
		}

		result := out.Bytes()
		return result, nil
	}
}

type GatewayHost struct {
	Id          int    `json:"id" gorm:"column:id"`
	UserID      string `json:"user_id" gorm:"column:user_id"` // user.userID
	HostID      string `json:"host_id" gorm:"column:host_id"` //uuid
	Description string `json:"description" gorm:"column:description"`

	Host       string `json:"host" gorm:"column:host"` // ip:port
	UserName   string `json:"username" gorm:"column:username"`
	PassWd     string `json:"password" gorm:"column:password"`
	PrivateKey string `json:"private_key" gorm:"column:private_key"`
}
