package main

import (
	"os"
	"time"

	"github.com/aenjoy/iot-lubricant/cmd/apiserver/internal"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/utils"
	"github.com/aenjoy/iot-lubricant/pkg/version"
	"github.com/aenjoy/iot-lubricant/services/apiserver"
	"github.com/aenjoy/iot-lubricant/services/corepkg/config"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

const (
	HTTP_LISTEN_PORT_STR   = "HTTP_LISTEN_PORT"
	LUBRICANT_HOSTNAME_STR = "HOSTNAME"
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

	listenPort := os.Getenv(HTTP_LISTEN_PORT_STR)
	hostName := os.Getenv(LUBRICANT_HOSTNAME_STR)
	app := apiserver.NewApp(
		apiserver.SetHostName(hostName),
		apiserver.SetPort(listenPort),
		apiserver.UseGinEngine(),
		//core.UseDB(repo.DefaultCoreClient()),
		apiserver.UseServerKey(),
		apiserver.UseCasdoor(),
		apiserver.UseSignalHandler(utils.HandelExitSignal(nil, config.SaveConfig, nil, 30*time.Second)),
	)
	panic(app.Run())
}

var (
	ServiceName       = "IoTLubricantCore-ApiServer"
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
