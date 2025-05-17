package main

import (
	"github.com/aenjoy/iot-lubricant/cmd/reporter/internal"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/version"
	"github.com/aenjoy/iot-lubricant/services/reporterHandler"

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

	app := reporterHandler.NewApp(
		reporterHandler.UseTaskQueue(internal.GetSyncTaskQueue()),
		reporterHandler.UseDataStore(internal.GetDataStore()))
	panic(app.Run())
}

func init() {
	version.ServiceName = "IoTLubricantCore-ReportHandler"
	version.PrintVersionInfo()
}
