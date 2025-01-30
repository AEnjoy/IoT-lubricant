package main

import (
	"flag"
	"os"

	"github.com/AEnjoy/IoT-lubricant/internal/app/core"
	appinit "github.com/AEnjoy/IoT-lubricant/internal/app/core/init"
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/router"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
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

	err := appinit.AppInit()
	if err != nil {
		panic(err)
	}

	listenPort := os.Getenv(HTTP_LISTEN_PORT_STR)
	hostName := os.Getenv(LUBRICANT_HOSTNAME_STR)
	app := core.NewApp(
		core.SetHostName(hostName),
		core.SetPort(listenPort),
		core.UseGinEngine(router.CoreRouter()),
		core.UseDB(repo.DefaultCoreClient()),
		core.UseServerKey(),
		core.UseCasdoor(),
	)
	panic(app.Run())
}
