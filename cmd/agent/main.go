package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	edge2 "github.com/AEnjoy/IoT-lubricant/internal/app/edge"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils"
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
	fmt.Printf("Runing Platform Info: %s/%s\n", runtime.GOOS, runtime.GOARCH)
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

	var config model.EdgeSystem
	_ = yaml.Unmarshal(f, &config)

	app := edge2.NewApp(
		edge2.UseCtrl(context.Background()),
		edge2.UseConfig(&config),
		edge2.UseGRPC(bindGrpc),
		edge2.UseHostAddress(hostname),
		edge2.UseOpenApi(openapi.NewOpenApiCli(config.FileName)),
		edge2.UseSignalHandler(utils.HandelExitSignal(nil, edge2.SaveConfig, nil, 30*time.Second)),
	)
	panic(app.Run())
}
