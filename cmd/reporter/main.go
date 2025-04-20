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

var (
	ServiceName       = "IoTLubricantCore-ReportHandler"
	Version           string
	BuildTime         string
	GoVersion         string
	GitCommit         string
	Features          string
	BuildHostPlatform string
	PlatformVersion   string
)

func init() {
	version.ServiceName = ServiceName
	version.Version = Version
	version.BuildTime = BuildTime
	version.GoVersion = GoVersion
	version.GitCommit = GitCommit
	version.Features = Features
	version.BuildHostPlatform = BuildHostPlatform
	version.PlatformVersion = PlatformVersion
	version.PrintVersionInfo()
}
