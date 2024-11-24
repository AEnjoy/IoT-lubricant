package test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/AEnjoy/IoT-lubricant/pkg/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/test"
	api "github.com/AEnjoy/IoT-lubricant/pkg/test/agent"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	pb "github.com/AEnjoy/IoT-lubricant/protobuf/agent"
	"github.com/google/uuid"
)

var _ Service = new(Mini)

type Mini struct {
}

func (Mini) App(cli pb.EdgeServiceClient, abort bool) error {
	var (
		r          *test.Result
		ctx              = context.Background()
		_1         int32 = 1
		_algorithm       = "-"
	)

	r = api.TestPing(cli)
	r.CheckResult(abort)

	test.GatewayID = uuid.NewString()

	r = api.TestRegisterGateway(cli)
	r.CheckResult(abort)

	originalApiData, err := os.ReadFile("mock_driver/clock/api.json")
	if err != nil {
		return err
	}

	var doc openapi.ApiInfo
	if err := json.Unmarshal(originalApiData, &doc); err != nil {
		return err
	}
	enableApi, err := edge.EnableApi(doc, &edge.Params{}, "/api/v1/get/time")
	if err != nil {
		return err
	}
	enableApiData, _ := json.Marshal(enableApi)

	_, err = cli.SetAgent(ctx, &pb.SetAgentRequest{
		AgentID: test.AgentID,
		AgentInfo: &pb.AgentInfo{
			AgentID:     test.AgentID,
			GatherCycle: &_1,
			Algorithm:   &_algorithm,
			DataSource: &pb.OpenapiDoc{
				OriginalFile: originalApiData,
				EnableFile:   enableApiData,
			},
		},
	})
	if err != nil {
		return err
	}

	startGatherResp, err := cli.StartGather(ctx, &pb.StartGatherRequest{})
	if err != nil {
		return err
	}
	if startGatherResp.Code != http.StatusOK {
		return errors.New(startGatherResp.Message)
	}

	data, err := cli.GetGatherData(ctx, &pb.GetDataRequest{AgentID: test.AgentID})
	if err != nil {
		return err
	}
	if data.GetInfo().GetCode() != http.StatusOK {
		return errors.New(data.GetInfo().GetMessage())
	}
	os.Exit(0)
	return nil
}
