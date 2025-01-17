package model

import (
	"fmt"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/container"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	"github.com/AEnjoy/IoT-lubricant/protobuf/proxy"
	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types/network"
)

type Device struct {
	Id     string `json:"id" gorm:"column:id;primary_key"`
	UserId string `json:"user_id" gorm:"column:user_id"`

	DeviceBasicInfo

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (Device) TableName() string {
	return "devices"
}

type DeviceBasicInfo struct {
	Name            string `json:"name" gorm:"column:name" validate:"required,max=64"`
	Type            string `json:"type" gorm:"column:type" validate:"required,max=64"`
	OperationSystem string `json:"os" gorm:"column:os" validate:"required,max=64"`
	Manufacturer    string `json:"manufacturer" gorm:"column:manufacturer" validate:"required,max=64"`
	Model           string `json:"model" gorm:"column:model" validate:"required,max=64"`
	Protocol        string `json:"protocol" gorm:"column:protocol" validate:"required,max=64"`
	Language        string `json:"language" gorm:"column:language" validate:"required,max=64"`
}

func (DeviceBasicInfo) TableName() string {
	return "device_basic_info"
}

type DeviceAPI struct {
	Method      string `json:"method" validate:"required,oneof=GET POST DELETE PUT"`
	Path        string `json:"path" validate:"path"`
	Description string `json:"description" validate:"required,max=1024"`
}

func (DeviceAPI) TableName() string {
	return "device_api"
}

type CreateAgentRequest struct { // CreateDriverAgentRequest
	AgentInfo Agent `json:"agent_info"`
	*CreateAgentConf
}
type CreateAgentConf struct {
	AgentContainerInfo  *container.Container `json:"agent_container_info,omitempty"`
	DriverContainerInfo *container.Container `json:"driver_container_info,omitempty"`
	OpenApiDoc          *openapi.OpenAPICli  `json:"open_api_doc,omitempty"`
}

func ProxypbCreateAgentRequest2CreateAgentRequest(pbreq *proxy.CreateAgentRequest) *CreateAgentRequest {
	var retVal = new(CreateAgentRequest)
	retVal.AgentInfo.GatewayId = pbreq.GetInfo().GetGatewayID()
	retVal.AgentInfo.AgentId = pbreq.GetInfo().GetAgentID()
	retVal.AgentInfo.Description = pbreq.GetInfo().GetDescription()
	retVal.AgentInfo.Algorithm = pbreq.GetInfo().GetAlgorithm()
	retVal.AgentInfo.GatherCycle = int(pbreq.GetInfo().GetGatherCycle())
	// retVal.AgentInfo.Cycle = int(pbreq.GetInfo().GetReportCycle())
	// retVal.AgentInfo.Address = int(pbreq.GetInfo().GetAddress())
	var conf CreateAgentConf
	if c := pbreq.GetConf(); len(c) > 0 {
		err := sonic.Unmarshal(c, &conf)
		if err != nil {
			logger.Errorf("can't unmarshal the content when conf has been set:%v\n", err)
		}
	}

	conf.OpenApiDoc = func() *openapi.OpenAPICli {
		ds := pbreq.GetInfo().GetDataSource()
		if ds == nil || len(ds.GetOriginalFile()) == 0 {
			return nil
		}
		var doc openapi.OpenAPICli
		err := sonic.Unmarshal(ds.GetOriginalFile(), &doc)
		if err != nil {
			logger.Errorf("can't unmarshal the content when doc has been set:%v\n", err)
			return nil
		}
		return &doc
	}()

	retVal.CreateAgentConf = &conf
	return retVal
}

func (CreateAgentRequest) TaskOperation() task.Operation {
	return task.OperationAddAgent
}

type CreateAgentResponse struct {
}

var AgentContainer = container.Container{
	Source: container.Image{
		PullWay:      2,
		FromRegistry: "hub.iotroom.top",
		RegistryPath: "AEnjoy/lubricant-agent",
		Tag:          "latest",
	},
	Name:    "lubricant-agent",
	Network: network.NetworkBridge,
	ExposePort: map[string]int{
		fmt.Sprintf("%d", _default.AgentGrpcPort): 0,
	},
}

// AgentInstance 记录agent 如何启动的信息

type AgentInstance struct {
	AgentId string `gorm:"column:agent_id"`

	CreateConf  string `gorm:"column:conf,type:json"`
	ContainerID string `gorm:"column:container_id"`
	IP          string `gorm:"column:ip"`
	Local       bool   `gorm:"column:local"` // 是否与 agent-proxy 在同一台机器上
	Online      bool   `gorm:"column:online"`
}

func (AgentInstance) TableName() string { return "agent_instance" }
