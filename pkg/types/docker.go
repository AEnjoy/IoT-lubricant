package types

import (
	"encoding/json"

	"github.com/AEnjoy/IoT-lubricant/pkg/types/task"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
)

type Container struct {
	Source     Image             `json:"source"`
	ImageName  string            `json:"image_name"`
	ExposePort map[string]int    `json:"expose_port,omitempty"` // map[containerPort]hostPort
	Network    string            `json:"network"`               // host, bridge, none(default)
	Env        map[string]string `json:"env,omitempty"`         // map[key]value
	Endpoint   string            `json:"endpoint"`              // optional
}
type Image struct {
	PullWay uint8 `json:"pull_way"` // 0: from binary, 1: from url, 2: from registry

	FromBinary   []byte `json:"from_binary"`   // a `tar` ball file
	FromUrl      string `json:"from_url"`      // such as https://example.com/mysql.tar
	FromRegistry string `json:"from_registry"` // such as `docker.io`
	RegistryPath string `json:"registry_path"` // such as `library/mysql`
	Tag          string `json:"tag"`           // such as `latest`
	// FromFile   string `json:"from_file"` // at `Core` side path of the image
}

type CreateAgentRequest struct { // CreateDriverAgentRequest
	AgentInfo           Agent              `json:"agent_info"`
	AgentContainerInfo  Container          `json:"agent_container_info"`
	DriverContainerInfo Container          `json:"driver_container_info"`
	OpenApiDoc          openapi.OpenAPICli `json:"open_api_doc"`
}

func (CreateAgentRequest) TaskOperation() task.Operation {
	return task.OperationAddAgent
}
func (this CreateAgentRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(this)
}

var AgentContainer = Container{
	Source: Image{
		PullWay:      2,
		FromRegistry: "hub.iotroom.top",
		RegistryPath: "AEnjoy/lubricant-agent",
		Tag:          "latest",
	},
	ImageName: "lubricant-agent",
	Network:   "bridge",
}
