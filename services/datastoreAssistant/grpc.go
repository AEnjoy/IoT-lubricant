package datastoreAssistant

import (
	"context"
	"io"

	metapb "github.com/aenjoy/iot-lubricant/protobuf/meta"
	svcpb "github.com/aenjoy/iot-lubricant/protobuf/svc"
	"github.com/aenjoy/iot-lubricant/services/datastoreAssistant/internal"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"google.golang.org/grpc"

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
func (a *app) StoreData(req grpc.ClientStreamingServer[svcpb.StoreDataRequest, svcpb.StoreDataResponse]) error {
	dataCh := make(chan any, a.internalThreadNumber)
	defer close(dataCh)

	recv, err := req.Recv()
	if err != nil {
		logg.L.Errorf("failed to recv: %v", err)
		return err
	}
	projectid := recv.GetProjectID()
	dataCh <- recv.GetData()

	go a._handel(req.Context(), projectid, dataCh)

	for {
		select {
		case <-req.Context().Done():
			return nil
		case <-a.Ctx.Done():
			err := req.SendAndClose(&svcpb.StoreDataResponse{
				ProjectID: projectid,
			})
			if err != nil {
				return err
			}
		default:
		}
		recv, err := req.Recv()
		if err != nil {
			if err == grpc.ErrClientConnClosing || err == io.EOF {
				return nil
			}
			return err
		}
		if req.Context().Err() == context.Canceled || req.Context().Err() == context.DeadlineExceeded {
			logg.L.Errorf("context canceled: %v", req.Context().Err())
			return nil
		}
		dataCh <- recv.GetData()
	}
}
func (a *app) Ping(_ context.Context, req *metapb.Ping) (*metapb.Ping, error) {
	if req.GetFlag() == 0 {
		return &metapb.Ping{Flag: 1}, nil
	}
	return &metapb.Ping{Flag: 2}, nil
}
