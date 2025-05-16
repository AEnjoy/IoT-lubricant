package datastoreAssistant

import (
	"context"

	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	"github.com/aenjoy/iot-lubricant/services/datastoreAssistant/internal"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *app) CheckLinker(_ context.Context, req *svcpb.CheckLinkerRequest) (*svcpb.CheckLinkerResponse, error) {
	switch r := req.GetRequest().(type) {
	case *svcpb.CheckLinkerRequest_Mysql:
		return &svcpb.CheckLinkerResponse{
			Result: svcpb.CheckLinkerResult(svcpb.CheckLinkerResult_value[internal.DsnTest("mysql", r.Mysql.GetDsn(), req.GetUserID())]),
		}, nil
	case *svcpb.CheckLinkerRequest_Tde:
		return &svcpb.CheckLinkerResponse{
			Result: svcpb.CheckLinkerResult(svcpb.CheckLinkerResult_value[internal.DsnTest("TDEngine", r.Tde.GetDsn(), req.GetUserID())]),
		}, nil
	}
	return nil, status.Errorf(codes.InvalidArgument, "invalid request type")
}
func (a *app) StoreData(context.Context, *svcpb.StoreDataRequest) (*svcpb.StoreDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreData not implemented")
}
func (a *app) Ping(_ context.Context, req *metapb.Ping) (*metapb.Ping, error) {
	if req.GetFlag() == 0 {
		return &metapb.Ping{Flag: 1}, nil
	}
	return &metapb.Ping{Flag: 2}, nil
}
