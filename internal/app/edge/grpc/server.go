package grpc

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/AEnjoy/IoT-lubricant/internal/app/edge/data"
	"github.com/AEnjoy/IoT-lubricant/internal/pkg/grpc/middleware"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/edge/config"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/types/exception/code"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	pb "github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
	json "github.com/bytedance/sonic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type agentServer struct {
	pb.UnimplementedEdgeServiceServer
}

func (*agentServer) Ping(context.Context, *meta.Ping) (*meta.Ping, error) {
	return &meta.Ping{Flag: 2}, nil
}

func (*agentServer) RegisterGateway(_ context.Context, request *pb.RegisterGatewayRequest) (*pb.RegisterGatewayResponse, error) {
	if config.Config.ID == "" {
		return nil, status.Error(codes.InvalidArgument, code.ErrorAgentNeedInit.GetMsg())
	}
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

func (a *agentServer) SetAgent(ctx context.Context, request *pb.SetAgentRequest) (*pb.SetAgentResponse, error) {
	var resp pb.SetAgentResponse
	if config.Config.ID != "" {
		if request.GetAgentID() != config.Config.ID {
			resp.Info = &meta.CommonResponse{Code: 500, Message: "target agentID error"}
			return &resp, nil
		}
	}

	resp.Info = &meta.CommonResponse{Code: 200, Message: "success"}
	if info := request.GetAgentInfo(); info != nil {
		config.Config.ID = info.AgentID
		if alg := info.Algorithm; alg != nil {
			config.Config.Algorithm = *alg
		}
		if ds := info.DataSource; ds != nil {
			if len(ds.OriginalFile) != 0 {
				logger.Info("Check origin config loaded")
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
				logger.Info("Check enable config loaded")
				var o openapi.Enable
				err := json.Unmarshal(ds.EnableFile, &o)
				if err != nil {
					resp.Info.Message = err.Error()
					resp.Info.Code = 500
					return &resp, err
				}

				if config.Config == nil {
					resp.Info.Message = "You should initialize the OriginalFile before setting the EnableFile"
					resp.Info.Code = 500
					return &resp, err
				}
				logger.Debugln(string(ds.EnableFile))
				c := config.Config.Config.(*openapi.ApiInfo)
				c.Enable = o
				config.Config.Config = c

				if !edge.CheckConfigInvalid(config.Config.Config) {
					return &pb.SetAgentResponse{Info: &meta.CommonResponse{
						Code:    http.StatusBadRequest,
						Message: "config invalid"}}, nil
				}
				_ = config.RefreshSlot()
				logger.Info("Enable config loaded----valid")
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
	} else {
		return &pb.SetAgentResponse{Info: &meta.CommonResponse{Code: http.StatusBadRequest, Message: "request body is empty"}}, nil
	}
	logger.Info("Saving config...")
	err := config.SaveConfig(config.SaveType_ALL)
	if err != nil {
		return &pb.SetAgentResponse{Info: &meta.CommonResponse{Code: 500, Message: "Can't save config"}}, err
	}

	if request.GetStop() != nil {
		_, err := a.StopGather(ctx, nil)
		if err != nil {
			return &pb.SetAgentResponse{Info: &meta.CommonResponse{Code: 500, Message: err.Error()}}, err
		}
	}
	if request.GetStart() != nil {
		_, err := a.StartGather(ctx, nil)
		if err != nil {
			return &pb.SetAgentResponse{Info: &meta.CommonResponse{Code: 500, Message: err.Error()}}, err
		}
	}
	return &resp, nil
}

func (*agentServer) GetOpenapiDoc(_ context.Context, request *pb.GetOpenapiDocRequest) (*pb.OpenapiDoc, error) {
	if config.Config.ID == "" {
		return nil, status.Error(codes.InvalidArgument, code.ErrorAgentNeedInit.GetMsg())
	}
	var o, e []byte
	if request.GetAgentID() != config.Config.ID {
		return nil, errors.New("target agentID error")
	}
	switch expression := request.DocType; expression {
	case pb.OpenapiDocType_originalFile:
		var err error
		if config.Config.Config != nil {
			o, err = json.Marshal(config.Config.Config.(*openapi.ApiInfo))
			if err != nil {
				return nil, err
			}
		}
	case pb.OpenapiDocType_enableFile:
		var err error
		if config.Config.Config.GetEnable() != nil {
			e, err = json.Marshal(config.Config.Config.GetEnable())
			if err != nil {
				return nil, err
			}
		}
	case pb.OpenapiDocType_All:
		var err error
		if config.Config.Config != nil {
			o, err = json.Marshal(config.Config.Config.(*openapi.ApiInfo))
			if err != nil {
				return nil, err
			}
		}
		if config.Config.Config.GetEnable() != nil {
			e, err = json.Marshal(config.Config.Config.GetEnable())
			if err != nil {
				return nil, err
			}
		}
	}

	return &pb.OpenapiDoc{
		EnableFile:   e,
		OriginalFile: o,
	}, nil
}

func (a *agentServer) GetAgentInfo(ctx context.Context, request *pb.GetAgentInfoRequest) (*pb.GetAgentInfoResponse, error) {
	if config.Config.ID == "" {
		return nil, status.Error(codes.InvalidArgument, code.ErrorAgentNeedInit.GetMsg())
	}
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

func (*agentServer) GetGatherData(_ context.Context, request *pb.GetDataRequest) (*pb.DataMessage, error) {
	if config.Config.ID == "" {
		return nil, status.Error(codes.InvalidArgument, code.ErrorAgentNeedInit.GetMsg())
	}
	if request.GetAgentID() != config.Config.ID {
		return nil, errors.New("target agentID error")
	}
	var resp pb.DataMessage
	l := data.Collector.GetDataLen(int(request.GetSlot()))
	if l > 0 {
		resp.DataLen = int32(l)
		_data := data.Collector.GetData(int(request.GetSlot()))
		resp.DataGatherStartTime = _data[0].Timestamp.Format("2006-01-02 15:04:05")
		resp.SplitTime = int32(config.Config.Cycle)
		for _, packet := range _data {
			resp.Data = append(resp.Data, packet.Data)
		}
		resp.Info = &meta.CommonResponse{Code: http.StatusOK, Message: "success"}
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

func (*agentServer) SendHttpMethod(_ context.Context, request *pb.SendHttpMethodRequest) (*pb.SendHttpMethodResponse, error) {
	if config.Config.ID == "" {
		return nil, status.Error(codes.InvalidArgument, code.ErrorAgentNeedInit.GetMsg())
	}
	var resp pb.SendHttpMethodResponse
	resp.Data = &pb.DataMessage{}
	if request.GetAgentID() != config.Config.ID {
		resp.Info = &meta.CommonResponse{Code: 500, Message: "target agentID error"}
		return &resp, nil
	}
	if !edge.CheckConfigInvalidGet(config.Config.Config) && config.Config.Config == nil {
		resp.Info = &meta.CommonResponse{Code: 500, Message: "Invalid internal configuration"}
		return &resp, nil
	}
	_, ok := config.Config.Config.GetPaths()[request.Path]
	if !ok {
		resp.Info = &meta.CommonResponse{Code: 500, Message: "Invalid path"}
		return &resp, nil
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
func (*agentServer) StartGather(ctx context.Context, _ *pb.StartGatherRequest) (*meta.CommonResponse, error) {
	if config.Config.ID == "" {
		return nil, status.Error(codes.InvalidArgument, code.ErrorAgentNeedInit.GetMsg())
	}
	ctx, cancel := utils.CreateTimeOutContext(ctx, utils.DefaultTimeout_Oper)
	defer cancel()
	if config.IsGathering() {
		return &meta.CommonResponse{Code: http.StatusInternalServerError, Message: "Gather is working now"}, nil
	}
	if !edge.CheckConfigInvalidGet(config.Config.Config) {
		return &meta.CommonResponse{Code: http.StatusInternalServerError, Message: "Invalid internal configuration"}, nil
	}
	select {
	case <-ctx.Done():
		return &meta.CommonResponse{Code: http.StatusInternalServerError, Message: "StartGather timeout"}, errors.New("timeout")
	case config.GatherSignal <- context.Background():
		return &meta.CommonResponse{Code: http.StatusOK, Message: "success"}, nil
	}
}
func (*agentServer) StopGather(ctx context.Context, _ *pb.StopGatherRequest) (*meta.CommonResponse, error) {
	if config.Config.ID == "" {
		return nil, status.Error(codes.InvalidArgument, code.ErrorAgentNeedInit.GetMsg())
	}
	ctx, cancel := utils.CreateTimeOutContext(ctx, utils.DefaultTimeout_Oper)
	defer cancel()

	if !config.IsGathering() {
		return &meta.CommonResponse{Code: http.StatusInternalServerError, Message: "Gather is not working"}, nil
	}
	select {
	case <-ctx.Done():
		return &meta.CommonResponse{Code: http.StatusInternalServerError, Message: "StopGather timeout"}, errors.New("timeout")
	case config.StopSignal <- context.Background():
		return &meta.CommonResponse{Code: http.StatusOK, Message: "success"}, nil
	}
}
func NewServer(bind string) {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.GetLoggerInterceptor(),
			middleware.GetRecovery(middleware.GetRegistry(middleware.GetSrvMetrics()))))
	pb.RegisterEdgeServiceServer(grpcServer, &agentServer{})
	logger.Infoln("agent grpc-server start at: ", bind)
	panic(grpcServer.Serve(lis))
}
