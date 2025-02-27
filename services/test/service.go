package test

import agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"

type Service interface {
	App(cli agentpb.EdgeServiceClient, abort, init bool) error
}
