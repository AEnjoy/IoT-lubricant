package main

import (
	"flag"
	"os"
	"time"

	"github.com/aenjoy/iot-lubricant/cmd/lubricant/internal"
	"github.com/aenjoy/iot-lubricant/pkg/utils"
	"github.com/aenjoy/iot-lubricant/pkg/version"
	"github.com/aenjoy/iot-lubricant/services/lubricant"
	"github.com/aenjoy/iot-lubricant/services/lubricant/config"
)

const (
	HTTP_LISTEN_PORT_STR   = "HTTP_LISTEN_PORT"
	LUBRICANT_HOSTNAME_STR = "HOSTNAME"
)

func main() {
	var envFilePath string
	flag.StringVar(&envFilePath, "env", "", "Path to .env file")
	flag.Parse()

	err := internal.AppInit()
	if err != nil {
		panic(err)
	}

	listenPort := os.Getenv(HTTP_LISTEN_PORT_STR)
	hostName := os.Getenv(LUBRICANT_HOSTNAME_STR)
	app := lubricant.NewApp(
		lubricant.SetHostName(hostName),
		lubricant.SetPort(listenPort),
		lubricant.UseGinEngine(),
		//core.UseDB(repo.DefaultCoreClient()),
		lubricant.UseServerKey(),
		lubricant.UseCasdoor(),
		lubricant.UseSignalHandler(utils.HandelExitSignal(nil, config.SaveConfig, nil, 30*time.Second)),
	)
	panic(app.Run())
}

var (
	ServiceName       = "IoTLubricantCore"
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
