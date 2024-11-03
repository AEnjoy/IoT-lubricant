package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/AEnjoy/IoT-lubricant/internal/edge"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/file"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/mq"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	CONF_FILE_ENV = "CONFIG"
	HOST_ENV      = "HOST"
	BIND_GRPC_ENV = "BIND_GRPC"
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

	configFile := os.Getenv(CONF_FILE_ENV)
	hostname := os.Getenv(HOST_ENV) //ip:port
	bindGrpc := os.Getenv(BIND_GRPC_ENV)

	f, err := os.ReadFile(configFile)
	if err != nil {
		logger.Warnln("Failed to read config file:", err)
	}

	var config types.EdgeSystem
	_ = yaml.Unmarshal(f, &config)

	if file.IsFileExists(config.FileName + ".enable") {
		config.FileName = config.FileName + ".enable"
	}

	app := edge.NewApp(
		edge.UseCtrl(context.Background()),
		edge.UseConfig(&config),
		edge.UseMq(mq.NewNatsMq[[]byte](fmt.Sprintf("nats://%s", hostname))),
		edge.UseGRPC(bindGrpc),
		edge.UseHostAddress(hostname),
		edge.UseOpenApi(openapi.NewOpenApiCli(config.FileName)),
	)
	panic(app.Run())
}
