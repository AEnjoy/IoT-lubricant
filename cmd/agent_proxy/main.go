package main

import (
	"flag"
	"os"

	"github.com/AEnjoy/IoT-lubricant/internal/app/gateway"
	"github.com/AEnjoy/IoT-lubricant/internal/model"
	"github.com/AEnjoy/IoT-lubricant/internal/model/repo"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/file"
	"github.com/joho/godotenv"
)

const (
	GATEWAY_ID_STR            = "GATEWAY_ID"
	MQ_LISTEN_PORT_STR        = "GATEWAY_MQ_PORT"
	CORE_HOST_STR             = "CORE_HOST"
	CORE_GRPC_LISTEN_PORT_STR = "CORE_PORT"
)

func main() {
	var envFilePath string
	var confFilePath string
	flag.StringVar(&envFilePath, "env", "", "Path to .env file")
	flag.StringVar(&confFilePath, "conf", "", "Path to .yaml file")
	flag.Parse()
	printBuildInfo()

	if envFilePath != "" {
		logger.Info("load env")
		err := godotenv.Load(envFilePath)
		if err != nil {
			logger.Info("Failed to load .env file, using system ones.")
		} else {
			logger.Infof("Loaded .env file from: %s", envFilePath)
		}
	}
	id := os.Getenv(GATEWAY_ID_STR)

	var config *model.ServerInfo
	if confFilePath != "" {
		config = new(model.ServerInfo)
		logger.Info("load conf")
		err := file.ReadYamlFile(confFilePath, config)
		if err != nil {
			logger.Info("Failed to load .yaml file, using system ones.")
		} else {
			logger.Infof("Loaded .yaml file from: %s", confFilePath)
		}
		id = config.GatewayID
	}

	app := gateway.NewApp(
		gateway.UseServerInfo(config),
		gateway.SetGatewayId(id),
		gateway.UseDB(repo.NewGatewayDb(nil)),
		gateway.LinkCoreServer(),
	)
	panic(app.Run())
}
