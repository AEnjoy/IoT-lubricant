package main

import (
	"os"
	"time"

	"github.com/aenjoy/iot-lubricant/cmd/apiserver/internal"
	"github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/utils"
	"github.com/aenjoy/iot-lubricant/pkg/version"
	"github.com/aenjoy/iot-lubricant/services/apiserver"
	"github.com/aenjoy/iot-lubricant/services/corepkg/config"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

func main() {
	var envFilePath string
	pflag.StringVar(&envFilePath, "env", "", "Path to .env file")
	pflag.Parse()
	if envFilePath != "" {
		logger.Info("load env")
		err := godotenv.Load(envFilePath)
		if err != nil {
			logger.Info("Failed to load .env file, using system ones.")
		} else {
			logger.Infof("Loaded .env file from: %s", envFilePath)
		}
	}

	err := internal.AppInit()
	if err != nil {
		panic(err)
	}

	listenPort := os.Getenv(constant.ENV_HTTP_LISTEN_PORT_STR)
	hostName := os.Getenv(constant.ENV_LUBRICANT_HOSTNAME_STR)
	app := apiserver.NewApp(
		apiserver.SetHostName(hostName),
		apiserver.SetPort(listenPort),
		apiserver.UseGinEngine(),
		//core.UseDB(repo.DefaultCoreClient()),
		apiserver.UseServerKey(),
		apiserver.UseCasdoor(),
		apiserver.UseSignalHandler(utils.HandelExitSignal(nil, config.SaveConfig, nil, 30*time.Second)),
	)
	panic(app.Run())
}

func init() {
	version.ServiceName = "IoTLubricantCore-ApiServer"
	version.PrintVersionInfo()
}
