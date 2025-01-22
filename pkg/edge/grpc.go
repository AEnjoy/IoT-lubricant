package edge

import (
	"github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAgentCli(address string) (agent.EdgeServiceClient, error) {
	cli, _, err := NewAgentCliWithClose(address)
	return cli, err
}
func NewAgentCliWithClose(address string) (agent.EdgeServiceClient, func(), error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	close := func() {
		_ = conn.Close()
	}
	return agent.NewEdgeServiceClient(conn), close, nil
}
