package main

import (
	"flag"
	"os"
	"time"

	"github.com/AEnjoy/IoT-lubricant/cmd/core/internal"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core"
	coreConfig "github.com/AEnjoy/IoT-lubricant/internal/app/core/config"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils"
)

const (
	HTTP_LISTEN_PORT_STR   = "HTTP_LISTEN_PORT"
	LUBRICANT_HOSTNAME_STR = "HOSTNAME"
)

func main() {
	var envFilePath string
	flag.StringVar(&envFilePath, "env", "", "Path to .env file")
	flag.Parse()
	printBuildInfo()

	err := internal.AppInit()
	if err != nil {
		panic(err)
	}

	listenPort := os.Getenv(HTTP_LISTEN_PORT_STR)
	hostName := os.Getenv(LUBRICANT_HOSTNAME_STR)
	app := core.NewApp(
		core.SetHostName(hostName),
		core.SetPort(listenPort),
		core.UseGinEngine(),
		//core.UseDB(repo.DefaultCoreClient()),
		core.UseServerKey(),
		core.UseCasdoor(),
		core.UseSignalHandler(utils.HandelExitSignal(nil, coreConfig.SaveConfig, nil, 30*time.Second)),
	)
	panic(app.Run())
}
