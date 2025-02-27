package main

import (
	"flag"
	"os"
	"time"

	"github.com/aenjoy/iot-lubricant/cmd/lubricant/internal"
	"github.com/aenjoy/iot-lubricant/pkg/utils"
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
	printBuildInfo()

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
