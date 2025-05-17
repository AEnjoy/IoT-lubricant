package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/model"
	"github.com/aenjoy/iot-lubricant/pkg/utils"
	"github.com/aenjoy/iot-lubricant/pkg/utils/file"
	"github.com/aenjoy/iot-lubricant/pkg/utils/openapi"
	"github.com/aenjoy/iot-lubricant/pkg/version"
	"github.com/aenjoy/iot-lubricant/services/agent"

	"github.com/joho/godotenv"
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

	configFile := os.Getenv(constant.ENV_CONF_FILE_ENV)
	hostname := os.Getenv(constant.ENV_HOST) //ip:port
	bindGrpc := os.Getenv(constant.ENV_BIND_GRPC)

	var config model.EdgeSystem
	if configFile == "" {
		logger.Warnln("No config file specified, using default values.")
		configFile = constant.AgentDefaultConfigFileName
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

func init() {
	version.ServiceName = "IoTEdgeAgent"
	version.PrintVersionInfo()
}
