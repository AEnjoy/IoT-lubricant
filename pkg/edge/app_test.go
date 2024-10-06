package edge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	"github.com/google/uuid"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func WriteConfig() {
	data, _ := json.Marshal(openAPIConfig.OpenAPICli)
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
			fmt.Sprintf("%s%s%s", v.Get("key1"), v.Get("key2"), uuid.NewString()),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
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

var config = &model.EdgeSystem{
	ID:          mockID,
	Cycle:       1,
	ReportCycle: 4,
	FileName:    "test/api.json",
}

func TestEdgeApp(t *testing.T) {
	t.Log("not all implement yet")
	WriteConfig()
	assert := assert.New(t)
	ctx, cf := context.WithDeadline(context.Background(), time.Now().Add(8*time.Second))
	// start nats server
	{
		opts := &server.Options{
			Port: 4222,
		}

		natsServer, err := server.NewServer(opts)
		assert.NoError(err)

		t.Log("starting nats server")
		go natsServer.Start()
		if !natsServer.ReadyForConnections(10 * time.Second) {
			t.Fatal("nats server did not start")
		}
		defer natsServer.Shutdown()
	}

	natsMq, err := mq.NewNatsMq[[]byte](nats.DefaultURL)
	assert.NoError(err)
	app := &app{
		config:  config,
		ctrl:    ctx,
		cancel:  cf,
		OpenApi: openAPIConfig,
		mq:      natsMq,
	}

	go assert.NoError(app.Run())
}
