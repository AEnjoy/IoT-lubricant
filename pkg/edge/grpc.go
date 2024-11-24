package edge

import (
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"google.golang.org/grpc"
)

func NewAgentCli(address string) (agent.EdgeServiceClient, error) {
	conn, err := grpc.NewClient(address)
	if err != nil {
		return nil, err
	}
	return agent.NewEdgeServiceClient(conn), nil
}
