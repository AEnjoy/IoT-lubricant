package test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/test"
	api "github.com/aenjoy/iot-lubricant/pkg/test/agent"
	"github.com/aenjoy/iot-lubricant/pkg/utils/openapi"
	agentpb "github.com/aenjoy/iot-lubricant/protobuf/agent"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

var _ Service = new(Mini)

type Mini struct {
}

func (Mini) App(cli agentpb.EdgeServiceClient, abort, init bool) error {
	var (
		r          *test.Result
		ctx              = context.Background()
		_1         int32 = 1
		_algorithm       = "-"
	)

	r = api.TestPing(cli)
	r.CheckResult(abort)

	test.GatewayID = uuid.NewString()
	if init {
		r = api.TestRegisterGateway(cli)
		r.CheckResult(abort)
	}

	originalApiData, err := os.ReadFile("mock_driver/clock/api.sonic")
	if err != nil {
		return err
	}

	var doc openapi.OpenAPICli
	if err := sonic.Unmarshal(originalApiData, &doc); err != nil {
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
	enableApiData, _ := sonic.Marshal(enableApi.(*openapi.ApiInfo).Enable)

	resp, err := cli.SetAgent(ctx, &agentpb.SetAgentRequest{
		AgentID: test.AgentID,
		AgentInfo: &agentpb.AgentInfo{
			AgentID:     test.AgentID,
			GatherCycle: &_1,
			Algorithm:   &_algorithm,
			DataSource: &agentpb.OpenapiDoc{
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
	if !init {
		r = api.TestRegisterGateway(cli)
		r.CheckResult(abort)
	}
	_, err = cli.GetOpenapiDoc(ctx, &agentpb.GetOpenapiDocRequest{
		AgentID: test.AgentID,
		DocType: agentpb.OpenapiDocType_originalFile,
	})
	if err != nil {
		return err
	}

	startGatherResp, err := cli.StartGather(ctx, &agentpb.StartGatherRequest{})
	if err != nil {
		return err
	}
	if startGatherResp.Code != http.StatusOK {
		return errors.New(startGatherResp.Message)
	}

	time.Sleep(5 * time.Second)
	data, err := cli.GetGatherData(ctx, &agentpb.GetDataRequest{AgentID: test.AgentID})
	if err != nil {
		return err
	}
	if data.GetInfo().GetCode() != http.StatusOK {
		return errors.New(data.GetInfo().GetMessage())
	}
	fmt.Println(data.String())

	_, err = cli.StopGather(ctx, &agentpb.StopGatherRequest{})
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}
