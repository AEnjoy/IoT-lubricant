package protobuf

// this file is used by mockery

import "google.golang.org/grpc"

type BidiStreamingServer[Req any, Res any] interface {
	Recv() (*Req, error)
	Send(*Res) error
	grpc.ServerStream
}
