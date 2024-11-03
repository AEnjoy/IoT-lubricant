package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/edge/data"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge/config"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	pb "github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
	"google.golang.org/grpc"
)

type agentServer struct {
	pb.UnimplementedEdgeServiceServer
}

func (a agentServer) Ping(_ context.Context, ping *meta.Ping) (*meta.Ping, error) {
	return &meta.Ping{Flag: 2}, nil
}

func (a agentServer) RegisterGateway(_ context.Context, request *pb.RegisterGatewayRequest) (*pb.RegisterGatewayResponse, error) {
	var resp pb.RegisterGatewayResponse
	if request.GetAgentID() != config.Config.ID {
		resp.AgentID = config.Config.ID
		resp.Info = &meta.CommonResponse{Code: 500, Message: "target agentID error"}
	} else {
		resp.AgentID = config.Config.ID
		config.GatewayID = request.GetGatewayID()
		resp.Info = &meta.CommonResponse{Code: 200, Message: "success"}
	}
	return &resp, nil
}

func (a agentServer) SetAgent(_ context.Context, request *pb.SetAgentRequest) (*pb.SetAgentResponse, error) {
	var resp pb.SetAgentResponse
	if request.GetAgentID() != config.Config.ID {
		resp.Info = &meta.CommonResponse{Code: 500, Message: "target agentID error"}
		return &resp, nil
	}

	resp.Info = &meta.CommonResponse{Code: 200, Message: "success"}
	if info := request.GetAgentInfo(); info != nil {
		config.Config.ID = info.AgentID
		if alg := info.Algorithm; alg != nil {
			config.Config.Algorithm = *alg
		}
		if ds := info.DataSource; ds != nil {
			if len(ds.OriginalFile) != 0 {
				var o openapi.OpenAPICli
				err := json.Unmarshal(ds.OriginalFile, &o)
				if err != nil {
					resp.Info.Message = err.Error()
					resp.Info.Code = 500
					return &resp, err
				}
				config.Config.Config = &openapi.ApiInfo{
					OpenAPICli: o,
				}
			}
			if len(ds.EnableFile) != 0 {
				var o openapi.OpenAPICli
				err := json.Unmarshal(ds.OriginalFile, &o)
				if err != nil {
					resp.Info.Message = err.Error()
					resp.Info.Code = 500
					return &resp, err
				}
				config.Config.EnableConfig = &openapi.ApiInfo{
					OpenAPICli: o,
				}
			}
		}
		if desc := info.Description; desc != nil {
			config.Config.Description = *desc
		}
		if gc := info.GatherCycle; gc != nil {
			config.Config.Cycle = int(*gc)
		}
		if gw := info.GatewayID; gw != nil {
			config.GatewayID = *gw
		}
		// todo stream
	}
	return &resp, nil
}

func (a agentServer) GetOpenapiDoc(_ context.Context, request *pb.GetOpenapiDocRequest) (*pb.OpenapiDoc, error) {
	var o, e []byte
	if request.GetAgentID() != config.Config.ID {
		return nil, errors.New("target agentID error")
	}
	switch expression := request.DocType; expression {
	case pb.OpenapiDocType_originalFile:
		var err error
		o, err = json.Marshal(config.Config.Config.(*openapi.ApiInfo))
		if err != nil {
			return nil, err
		}
	case pb.OpenapiDocType_enableFile:
		var err error
		e, err = json.Marshal(config.Config.EnableConfig.(*openapi.ApiInfo))
		if err != nil {
			return nil, err
		}
	case pb.OpenapiDocType_All:
		var err error
		o, err = json.Marshal(config.Config.Config.(*openapi.ApiInfo))
		if err != nil {
			return nil, err
		}
		e, err = json.Marshal(config.Config.EnableConfig.(*openapi.ApiInfo))
		if err != nil {
			return nil, err
		}
	}

	return &pb.OpenapiDoc{
		EnableFile:   e,
		OriginalFile: o,
	}, nil
}

func (a agentServer) GetAgentInfo(ctx context.Context, request *pb.GetAgentInfoRequest) (*pb.GetAgentInfoResponse, error) {
	c := int32(config.Config.Cycle)
	ds, err := a.GetOpenapiDoc(ctx, &pb.GetOpenapiDocRequest{AgentID: request.GetAgentID(), DocType: pb.OpenapiDocType_All})
	if err != nil {
		return nil, err
	}
	return &pb.GetAgentInfoResponse{
		AgentInfo: &pb.AgentInfo{
			AgentID:     config.Config.ID,
			Algorithm:   &config.Config.Algorithm,
			Description: &config.Config.Description,
			GatherCycle: &c,
			GatewayID:   &config.GatewayID,
			DataSource:  ds,
		},
		Info: &meta.CommonResponse{Code: 200, Message: "success"},
	}, nil
}

func (a agentServer) Data(_ context.Context, request *pb.GetDataRequest) (*pb.DataMessage, error) {
	if request.GetAgentID() != config.Config.ID {
		return nil, errors.New("target agentID error")
	}
	var resp pb.DataMessage
	data.DCL.Lock()
	defer data.DCL.Unlock()
	if len(data.DataCollect) != 0 {
		resp.DataLen = int32(len(data.DataCollect))
		resp.DataGatherStartTime = data.DataCollect[0].Timestamp.Format("2006-01-02 15:04:05")
		resp.SplitTime = int32(config.Config.Cycle)
		for _, packet := range data.DataCollect {
			resp.Data = append(resp.Data, packet.Data)
		}
		resp.Info = &meta.CommonResponse{Code: http.StatusOK, Message: "success"}
		data.DataCollect = make([]*edge.DataPacket, 0)
		return &resp, nil
	} else {
		resp.Info = &meta.CommonResponse{Code: http.StatusTooEarly, Message: "data is not ready"}
		return &resp, nil
	}
}

//func (a agentServer) GetDataStream(request *pb.GetDataRequest, g grpc.ServerStreamingServer[pb.DataChunk]) error {
//	//TODO implement me
//	panic("implement me")
//}

func (a agentServer) SendHttpMethod(_ context.Context, request *pb.SendHttpMethodRequest) (*pb.SendHttpMethodResponse, error) {
	var resp pb.SendHttpMethodResponse
	resp.Data = &pb.DataMessage{}
	if request.GetAgentID() != config.Config.ID {
		resp.Info = &meta.CommonResponse{Code: 500, Message: "target agentID error"}
		return &resp, errors.New("target agentID error")
	}
	switch request.Method {
	case http.MethodGet:
		var par []openapi.Parameter
		kvMap := request.GetParams().(*pb.SendHttpMethodRequest_Kv).Kv.GetKv()
		for k, v := range kvMap {
			parameter := openapi.Parameter{}
			parameter.Set(k, v)
			par = append(par, parameter)
		}
		respData, err := config.Config.Config.SendGETMethod(request.GetPath(), par)
		if err != nil {
			return nil, err
		}
		resp.Info = &meta.CommonResponse{Code: 200, Message: "success"}
		resp.Data.Data = [][]byte{respData}
		resp.Data.DataLen = 1
		resp.Data.SplitTime = 1
		resp.Data.DataGatherStartTime = time.Now().Format("2006-01-02 15:04:05")
		return &resp, nil
	case http.MethodPost:
		body := request.GetParams().(*pb.SendHttpMethodRequest_Body).Body.GetBody()
		// todo: get content type from request
		dataResp, err := config.Config.Config.SendPOSTMethodEx(request.GetPath(), "application/json", body)
		if err != nil {
			return nil, err
		}
		resp.Info = &meta.CommonResponse{Code: 200, Message: "success"}
		resp.Data.Data = [][]byte{dataResp}
		resp.Data.DataLen = 1
		resp.Data.SplitTime = 1
		resp.Data.DataGatherStartTime = time.Now().Format("2006-01-02 15:04:05")
		return &resp, nil
	}
	return nil, errors.New("method not support")
}

func NewServer(bind string) {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterEdgeServiceServer(grpcServer, &agentServer{})
	panic(grpcServer.Serve(lis))
}
