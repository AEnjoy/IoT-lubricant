package main

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"strings"

	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"google.golang.org/genproto/googleapis/rpc/status"
)

var agentsID []string
var _agentIDPoolSize int

func initAgentIDPools(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	agentsID = strings.Split(string(file), "\n")
	for i := range agentsID {
		agentsID[i] = strings.TrimSpace(agentsID[i])
	}
	_agentIDPoolSize = len(agentsID)
	if _agentIDPoolSize == 0 {
		return errors.New("agentID pool is empty")
	}
	return nil
}
func randGetAgentID() string {
	return agentsID[rand.Intn(_agentIDPoolSize)]
}

func regAgentOnline(cli corepb.CoreServiceClient, ctx context.Context) {
	for _, id := range agentsID {
		cli.Report(ctx, &corepb.ReportRequest{
			GatewayId: gatewayID,
			AgentId:   id,
			Req: &corepb.ReportRequest_AgentStatus{
				AgentStatus: &corepb.AgentStatusRequest{
					Req: &status.Status{Message: "online"},
				},
			},
		})
	}
}
func regAgentOffline(cli corepb.CoreServiceClient, ctx context.Context) {
	for _, id := range agentsID {
		cli.Report(ctx, &corepb.ReportRequest{
			GatewayId: gatewayID,
			AgentId:   id,
			Req: &corepb.ReportRequest_AgentStatus{
				AgentStatus: &corepb.AgentStatusRequest{
					Req: &status.Status{Message: "offline"},
				},
			},
		})
	}
}
