package test

import pb "github.com/AEnjoy/IoT-lubricant/protobuf/agent"

type Service interface {
	App(cli pb.EdgeServiceClient, abort bool) error
}
