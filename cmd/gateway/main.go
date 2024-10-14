package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/AEnjoy/IoT-lubricant/pkg/gateway"
	"github.com/AEnjoy/IoT-lubricant/pkg/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/joho/godotenv"
)

const (
	GATEWAY_ID_STR            = "GATEWAY_ID"
	MQ_LISTEN_PORT_STR        = "GATEWAY_MQ_PORT"
	CORE_HOST_STR             = "CORE_HOST"
	CORE_GRPC_LISTEN_PORT_STR = "CORE_PORT"
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

	port := os.Getenv(MQ_LISTEN_PORT_STR)
	id := os.Getenv(GATEWAY_ID_STR)

	app := gateway.NewApp(
		gateway.SetGatewayId(id),
		gateway.SetPort(port),
		gateway.UseDB(model.NewGatewayDb(nil)),
	)
	panic(app.Run())
}
