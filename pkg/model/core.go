package model

import (
	"bytes"
	"database/sql"
	"errors"
	"os"
	"os/exec"
	"time"

	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
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

	CreatedAt time.Time    `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeleteAt  sql.NullTime `json:"deleteAt" gorm:"column:deleted_at;type:datetime"`
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
	ID      int    `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	AgentID string `json:"agent_id" gorm:"column:agent_id"` // equal to  DeviceID

	Content string `json:"data" gorm:"column:data;serializer:json"` // core.Data 序列化的json

	CreatedAt time.Time    `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeleteAt  sql.NullTime `json:"deleteAt" gorm:"column:deleted_at;type:datetime"`
}

func (Data) TableName() string {
	return "data"
}

type Clean struct {
	ID          int    `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	AgentID     string `json:"agent_id" gorm:"column:agent_id"`
	Description string `json:"description" gorm:"column:description"`

	Interpreter string `json:"interpreter" gorm:"column:interpreter"` // python,goja,node,bash or other
	Script      string `json:"script" gorm:"column:script"`           // 脚本代码
	Command     string `json:"command" gorm:"column:command"`         // 提供给解释器的额外参数

	CreatedAt time.Time    `json:"-" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time    `json:"-" gorm:"column:updated_at;type:datetime"`
	DeleteAt  sql.NullTime `json:"-" gorm:"column:deleted_at;type:datetime"`
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

		processData, ok := goja.AssertFunction(rt.Get(c.Command))
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
		newCommand = append(newCommand, c.Command)
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
	Id          int    `json:"-" gorm:"column:id;primary_key;autoIncrement"`
	UserID      string `json:"-" gorm:"column:user_id"`       // user.userID
	HostID      string `json:"host_id" gorm:"column:host_id"` //uuid
	Description string `json:"description" gorm:"column:description"`

	Host       string `json:"host" gorm:"column:host"` // ip:port
	UserName   string `json:"username" gorm:"column:username"`
	PassWd     string `json:"-" gorm:"column:password"`
	PrivateKey string `json:"-" gorm:"column:private_key"`
	PublicKey  string `json:"-" gorm:"column:public_key"`

	CreatedAt time.Time    `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeleteAt  sql.NullTime `json:"deleteAt" gorm:"column:deleted_at;type:datetime"`
}
type ErrorLogs struct {
	ID        int    `json:"-" gorm:"column:id;primary_key;autoIncrement"`
	ErrID     string `json:"err_id" gorm:"column:err_id"`
	Component string `json:"component" gorm:"column:component"` // one of core,agent,gateway

	Type    int32  `json:"type" gorm:"column:type"`
	Code    int32  `json:"code" gorm:"column:code"`
	Message string `json:"message" gorm:"column:message"`
	Module  string `json:"module" gorm:"column:module"`
	Stack   string `json:"stack" gorm:"column:stack"`

	CreatedAt time.Time    `json:"happened" gorm:"column:created_at;type:datetime"`
	DeleteAt  sql.NullTime `json:"deleteAt" gorm:"column:deleted_at;type:datetime"`
}

func (ErrorLogs) TableName() string {
	return "error_logs"
}
func PbErrorMessage2ModelErrorLogs(message *metapb.ErrorMessage) *ErrorLogs {
	return &ErrorLogs{
		Code: func() int32 {
			if status := message.GetCode(); status != nil {
				return status.GetCode()
			}
			return 0
		}(),
		Message: func() string {
			if status := message.GetCode(); status != nil {
				return status.GetMessage()
			}
			return ""
		}(),
		Module: message.GetModule(),
		Stack:  message.GetStack(),
		Type:   message.GetErrorType(),
		CreatedAt: func() time.Time {
			if timestamp := message.GetTime(); timestamp != nil {
				return timestamp.AsTime()
			}
			return time.Now()
		}(),
	}
}

type AsyncJob struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	UserID     string `gorm:"column:user_id" json:"-"`
	Name       string `gorm:"type:varchar(255);not null;column:name" json:"name"` // task name
	RequestID  string `gorm:"type:varchar(255);not null;unique;column:request_id" json:"requestId"`
	Status     string `gorm:"column:status;type:enum('completed', 'failed', 'pending', 'retried', 'retrying', 'started');not null" json:"status"`
	Data       string `gorm:"column:data;type:json;not null" json:"data"`
	ResultData string `gorm:"column:result_data;type:text" json:"resultData"`

	ExpiredAt time.Time    `gorm:"type:datetime;not null" json:"expiredAt"`
	CreatedAt time.Time    `gorm:"type:datetime;not null" json:"createdAt"`
	UpdatedAt time.Time    `gorm:"type:datetime;not null" json:"updatedAt"`
	DeleteAt  sql.NullTime `gorm:"column:deleted_at;type:datetime" json:"deleteAt"`
	//Meta      string    `gorm:"column:meta;type:json;not null" json:"meta"`
}

func (AsyncJob) TableName() string {
	return "async_job"
}

type GatherNodeConfig struct {
	ID           int    `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	ConfigID     string `gorm:"type:varchar(255);not null;column:config_id;unique" json:"config_id"`
	RootConfigID string `gorm:"type:varchar(255);column:root_config_id" json:"root_config_id"` //  root GatherNodeConfig.ConfigID
	Name         string `gorm:"type:varchar(255);not null;column:name;unique" json:"name"`
	Description  string `gorm:"type:text;column:description" json:"description"`

	Config       string `gorm:"type:text;column:config" json:"config"`               // json or yaml
	EnableConfig string `gorm:"type:text;column:enable_config" json:"enable_config"` // json or yaml

	CreatedAt time.Time    `gorm:"type:datetime;not null" json:"createdAt"`
	UpdatedAt time.Time    `gorm:"type:datetime" json:"updatedAt"`
	DeleteAt  sql.NullTime `gorm:"column:deleted_at;type:datetime" json:"deleteAt"`
}

func (GatherNodeConfig) TableName() string {
	return "gather_node_config"
}
