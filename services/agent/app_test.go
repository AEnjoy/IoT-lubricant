package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/utils/openapi"
	json "github.com/bytedance/sonic"
	"github.com/bytedance/sonic/encoder"
	"github.com/google/uuid"
)

func WriteConfig() {
	data, _ := json.Marshal(openAPIConfig.OpenAPICli)
	_ = os.Mkdir("test", 0755)
	_ = os.WriteFile("test/api.json", data, 0644)
	_ = os.WriteFile("test/api.json.enable", data, 0644)
}

var testServer *httptest.Server

func StartTestServer() string {
	if testServer != nil {
		return testServer.URL
	}
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query()
		v.Get("key1")
		resp := struct {
			Data string `json:"data"`
		}{
			fmt.Sprintf("%s-%s-%s", v.Get("key1"), v.Get("key2"), uuid.NewString()),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = encoder.NewStreamEncoder(w).Encode(resp)
	}))
	//go testServer.Start()
	<-time.After(500 * time.Millisecond)
	return testServer.URL
}

var mockID = uuid.NewString()

var openAPIConfig = &openapi.ApiInfo{
	OpenAPICli: openapi.OpenAPICli{
		Servers: []openapi.Server{
			{
				URL: StartTestServer(),
			},
		},
		Paths: map[string]openapi.PathItem{
			"/test1": {
				Get: &openapi.Operation{
					Parameters: []openapi.Parameter{
						{
							Name:     "key1",
							Required: true,
							Schema: openapi.Schema{
								Properties: map[string]openapi.Property{
									"key1": {
										Type: "string-key1",
									},
								},
							},
						},
						{
							Name:     "key2",
							Required: true,
							Schema: openapi.Schema{
								Properties: map[string]openapi.Property{
									"key2": {
										Type: "string-key2",
									},
								},
							},
						},
					},
					Responses: map[string]openapi.Response{
						"200": {
							Description: "ok",
						},
					},
				},
			},
		},
	},
}

var testConfig = &model.EdgeSystem{
	ID:          mockID,
	Cycle:       1,
	ReportCycle: 4,
	FileName:    "test/api.json",
}

func TestEdgeApp(t *testing.T) {
	t.Skip("Deprecated")
	//t.Log("This test will take about 10s to complete")
	//t.Log("Start Time:", time.Now())
	//WriteConfig()
	//assert := assert.New(t)
	//ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(8*time.Second))
	//// start nats server
	//{
	//	opts := &server.Options{
	//		Port:  4222,
	//		Debug: true,
	//	}
	//
	//	natsServer, err := server.NewServer(opts)
	//	assert.NoError(err)
	//
	//	t.Log("starting nats server")
	//	go natsServer.Start()
	//	if !natsServer.ReadyForConnections(5 * time.Second) {
	//		t.Fatal("nats server did not start")
	//	}
	//	defer natsServer.Shutdown()
	//}
	//
	//natsMq, err := mq.NewNatsMq[[]byte](nats.DefaultURL)
	//assert.NoError(err)
	//app := &app{
	//	config:  config,
	//	ctrl:    ctx,
	//	cancel:  cf,
	//	OpenApi: openAPIConfig,
	//	mq:      natsMq,
	//}
	//
	//go func() {
	//	assert.NoError(app.Run())
	//}()
	//
	//// test register
	//regCh, err := app.mq.Subscribe(types.Topic_AgentRegister + mockID)
	//assert.NoError(err)
	//
	//var reg types.Register
	//assert.NoError(json.Unmarshal(<-regCh, &reg))
	//assert.Equal(mockID, reg.ID)
	//
	//ping := types.Ping{
	//	Status: 1,
	//}
	//data, _ := json.Marshal(ping)
	//_ = app.mq.Publish(types.Topic_AgentRegisterAck+mockID, data)
	//
	//// mq test (mock-gateway)
	//var success bool
	//t.Log("Test send message to topic")
	//ch, err := app.mq.Subscribe(types.Topic_AgentDataPush + mockID)
	//assert.NoError(err)
	//for {
	//	select {
	//	case <-ctx.Done():
	//		assert.True(success, "test failed because no data send success")
	//		return
	//	case d := <-ch:
	//		var data gateway.DataMessage
	//		t.Log("receive data from topic:", time.Now())
	//		assert.NoError(json.Unmarshal(d, &data))
	//		assert.Equal(int32(2), data.Flag)
	//		assert.Equal(mockID, data.AgentId)
	//		success = true
	//	}
	//}

}
