package main

import (
	"context"
	"flag"
	"os"
	"time"

	def "github.com/aenjoy/iot-lubricant/pkg/default"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/utils"
	"github.com/aenjoy/iot-lubricant/pkg/utils/file"
	"github.com/aenjoy/iot-lubricant/pkg/utils/openapi"
	"github.com/aenjoy/iot-lubricant/pkg/version"
	"github.com/aenjoy/iot-lubricant/services/agent"
	"github.com/joho/godotenv"
)

const (
	CONF_FILE_ENV = "CONFIG"
	HOST_ENV      = "HOST"
	BIND_GRPC_ENV = "BIND_GRPC"
)

func main() {
	var envFilePath string
	flag.StringVar(&envFilePath, "env", "", "Path to .env file")
	flag.Parse()

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

	var config model.EdgeSystem
	if configFile == "" {
		logger.Warnln("No config file specified, using default values.")
		configFile = def.AgentDefaultConfigFileName
	}
	if err := file.ReadYamlFile(configFile, &config); err != nil {
		logger.Warnln("Failed to read config file:", err)
	}

	app := agent.NewApp(
		agent.UseCtrl(context.Background()),
		agent.UseConfig(&config),
		agent.UseGRPC(bindGrpc),
		agent.UseHostAddress(hostname),
		agent.UseOpenApi(openapi.NewOpenApiCli(config.FileName)),
		agent.UseSignalHandler(utils.HandelExitSignal(nil, agent.SaveConfig, nil, 30*time.Second)),
	)
	panic(app.Run())
}

var (
	ServiceName       = "IoTEdgeAgent"
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
