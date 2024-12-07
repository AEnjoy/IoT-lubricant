package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

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

	var doc openapi.OpenAPICli
	if err := json.Unmarshal(originalApiData, &doc); err != nil {
		return err
	}

	apiInfo, err := openapi.NewOpenApiCliEx(originalApiData, nil)
	if err != nil {
		return err
	}

	enableApi, err := openapi.EnableApi(apiInfo, &openapi.EnableParams{GetParams: map[string]string{}}, "/api/v1/get/time")
	if err != nil {
		return err
	}
	enableApiData, _ := json.Marshal(enableApi.(*openapi.ApiInfo).Enable)

	resp, err := cli.SetAgent(ctx, &pb.SetAgentRequest{
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
	if resp.GetInfo().Code != http.StatusOK {
		println("failed")
		return errors.New(resp.GetInfo().GetMessage())
	}

	_, err = cli.GetOpenapiDoc(ctx, &pb.GetOpenapiDocRequest{
		AgentID: test.AgentID,
		DocType: pb.OpenapiDocType_originalFile,
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

	time.Sleep(5 * time.Second)
	data, err := cli.GetGatherData(ctx, &pb.GetDataRequest{AgentID: test.AgentID})
	if err != nil {
		return err
	}
	if data.GetInfo().GetCode() != http.StatusOK {
		return errors.New(data.GetInfo().GetMessage())
	}
	fmt.Println(data.String())

	_, err = cli.StopGather(ctx, &pb.StopGatherRequest{})
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}
