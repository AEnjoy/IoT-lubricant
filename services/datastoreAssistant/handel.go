package datastoreAssistant

import (
	"context"
	"fmt"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/utils/compress"
	corepb "github.com/aenjoy/iot-lubricant/protobuf/core"
	"github.com/aenjoy/iot-lubricant/services/datastoreAssistant/driver"
	logg "github.com/aenjoy/iot-lubricant/services/logg/api"
	"google.golang.org/protobuf/proto"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/base64x"
)

func (a *app) handelProjectIDStr(projectIdStr any) {
	topic := fmt.Sprintf(constant.DATASTORE_PROJECT_DATA, projectIdStr.(string))
	subscribeCh, err := a.DataStore.V2mq.QueueSubscribe(topic)
	if err != nil {
		logg.L.Errorf("Failed to subscribe to topic %s: %s", projectIdStr.(string), err)
		return
	}
	defer a.DataStore.V2mq.QueueUnsubscribe(topic)
	a._handel(a.Ctx, projectIdStr.(string), subscribeCh)
}

func (a *app) _handel(ctx context.Context, projectID string, dataCh <-chan any) {
	logg.L.Debugf("[%s] Handel called", projectID)
	project, err := a.DataStore.GetProject(ctx, projectID)
	if err != nil {
		logg.L.Errorf("Failed to get project by projectId %s: %s", projectID, err)
		return
	}

	txn := a.DataStore.Begin()
	agents, err := a.DataStore.GetAgentsByProjectID(ctx, txn, projectID)
	a.DataStore.Commit(txn)
	if err != nil {
		logg.L.Errorf("Failed to get agents by projectId %s: %s", projectID, err)
		return
	}
	if len(agents) < 1 {
		logg.L.Errorf("No agents found for projectId %s", projectID)
		return
	}
	agent := agents[0]

	decompressor, err := compress.NewCompressor(agent.Algorithm)
	if err != nil {
		logg.L.Errorf("Failed to create decompressor: %s", err)
		return
	}

	engine, err := a.DataStore.GetEngineByProjectID(ctx, projectID)
	if err != nil {
		logg.L.Errorf("Failed to get engine by projectId %s: %s", projectID, err)
		return
	}
	var (
		dri    driver.IDriver
		closer func() error
	)

	switch engine.DataBaseType {
	case "mysql":
		dsn, _ := base64x.StdEncoding.DecodeString(engine.DSN)
		dri, closer, err = driver.NewMySQLDriver(string(dsn), engine.Table, project.UserID)
	case "TDEngine":
		dsn, _ := base64x.StdEncoding.DecodeString(engine.DSN)
		var tdmodel model.LinkerInfo
		_ = sonic.Unmarshal(dsn, &tdmodel)
		dri, closer, err = driver.NewTDEngineDriver(
			project.UserID, tdmodel.Host, tdmodel.User,
			tdmodel.Pass, tdmodel.Db, tdmodel.Port,
			&engine.Table, tdmodel.Schemaless,
		)
	default:
		logg.L.Errorf("Unsupported database type %s", engine.DataBaseType)
		return
	}
	if err != nil {
		logg.L.Errorf("Failed to create driver: %s userId:%s projectId:%s", err, project.UserID, projectID)
		return
	}
	defer closer()

	for {
		select {
		case <-ctx.Done():
			logg.L.Infof("Context done, exiting...")
			return
		case data, ok := <-dataCh:
			if !ok {
				logg.L.Infof("Data channel closed, exiting...")
				return
			}
			var pbdata corepb.Data
			err := proto.Unmarshal(data.([]byte), &pbdata)
			if err != nil {
				logg.L.Errorf("Failed to unmarshal data: %s", err)
				continue
			}

			for _, d := range pbdata.GetData() {
				deData, err := decompressor.Decompress(d)
				if err != nil {
					logg.L.Errorf("Failed to decompress data: %s", err)
					continue
				}

				_, err = dri.Write(deData)
				if err != nil {
					logg.L.Errorf("Failed to write data: %s", err)
				}
			}
		}
	}
}
