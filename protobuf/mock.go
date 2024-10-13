package protobuf

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// this file is used by mockery

type BidiStreamingServer[Req any, Res any] interface {
	Recv() (*Req, error)
	Send(*Res) error

	Header() (metadata.MD, error)
	Trailer() metadata.MD
	CloseSend() error
	Context() context.Context
	SendMsg(m any) error
	RecvMsg(m any) error
}
