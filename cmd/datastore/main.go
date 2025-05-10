package main

import (
	"context"
	"time"

	"github.com/aenjoy/iot-lubricant/cmd/datastore/internal"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/utils"
	"github.com/aenjoy/iot-lubricant/pkg/version"
	"github.com/aenjoy/iot-lubricant/services/datastoreAssistant"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

func main() {
	var envFilePath string
	pflag.StringVar(&envFilePath, "env", "", "Path to .env file")
	pflag.Parse()
	if envFilePath != "" {
		logger.Info("load env")
		err := godotenv.Load(envFilePath)
		if err != nil {
			logger.Info("Failed to load .env file, using system ones.")
		} else {
			logger.Infof("Loaded .env file from: %s", envFilePath)
		}
	}

	err := internal.AppInit()
	if err != nil {
		panic(err)
	}
	app := datastoreAssistant.NewApp(
		datastoreAssistant.WithContext(context.Background()),
		datastoreAssistant.NewEtcdClient(internal.GetEtcdEndpoints()),
		datastoreAssistant.SetDataStore(internal.GetDataStore()),
		datastoreAssistant.SetThreadNumber(internal.GetInternalWorkThreadNumber()),
		datastoreAssistant.UseSignalHandler(utils.HandelExitSignal(nil, datastoreAssistant.ExitHandel, nil, 30*time.Second)),
	)
	panic(app.Run())
}

func init() {
	version.ServiceName = "IoTLubricantCore-DatastoreAssistant"
	version.PrintVersionInfo()
}
