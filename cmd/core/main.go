package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/AEnjoy/IoT-lubricant/pkg/core"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/router"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/joho/godotenv"
)

const (
	HTTP_LISTEN_PORT_STR   = "HTTP_LISTEN_PORT"
	LUBRICANT_HOSTNAME_STR = "HOSTNAME"
)

var (
	Version         string
	BuildTime       string
	GoVersion       string
	GitTag          string
	Features        string
	Platform        string
	PlatformVersion string
)

func printBuildInfo() {
	fmt.Printf("IoT-lubricant-Version: %s\n", Version)
	fmt.Printf("Build-Time: %s\n", BuildTime)
	fmt.Printf("Go-Version: %s\n", GoVersion)
	fmt.Printf("Git-Tag: %s\n", GitTag)
	fmt.Printf("Features: %s\n", Features)
	fmt.Printf("Platform: %s\n", Platform)
	fmt.Printf("Platform-Version: %s\n", PlatformVersion)
	fmt.Printf("Runing Platform Info: %s/%s", runtime.GOOS, runtime.GOARCH)
}

func main() {
	var envFilePath string
	flag.StringVar(&envFilePath, "env", "", "Path to .env file")
	flag.Parse()
	printBuildInfo()

	if envFilePath != "" {
		logger.Info("load env")
		err := godotenv.Load(envFilePath)
		if err != nil {
			logger.Info("Failed to load .env file, using system ones.")
		} else {
			logger.Infof("Loaded .env file from %s", envFilePath)
		}
	}

	listenPort := os.Getenv(HTTP_LISTEN_PORT_STR)
	hostName := os.Getenv(LUBRICANT_HOSTNAME_STR)
	app := core.NewApp(
		core.SetHostName(hostName),
		core.SetPort(listenPort),
		core.UseGinEngine(router.CoreRouter()),
		core.UseDB(model.DefaultCoreClient()),
	)
	panic(app.Run())
}
