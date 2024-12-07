package edge

import (
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAgentCli(address string) (agent.EdgeServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return agent.NewEdgeServiceClient(conn), nil
}
