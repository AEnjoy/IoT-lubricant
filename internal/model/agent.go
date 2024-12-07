package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/container"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
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
	AgentInfo           Agent               `json:"agent_info"`
	AgentContainerInfo  container.Container `json:"agent_container_info"`
	DriverContainerInfo container.Container `json:"driver_container_info"`
	OpenApiDoc          openapi.OpenAPICli  `json:"open_api_doc"`
}

func (CreateAgentRequest) TaskOperation() task.Operation {
	return task.OperationAddAgent
}
func (this CreateAgentRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(this)
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
