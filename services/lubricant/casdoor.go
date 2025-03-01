package lubricant

import (
	"fmt"
	"os"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/services/lubricant/config"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

func initCasdoor() error {
	c := config.GetConfig()
	file, err := os.ReadFile(c.AuthPublicKeyFile)
	if err != nil {
		return err
	}
	if os.Getenv("RUNNING_LEVEL") == "debug" {
		logger.Debug("CERT FILE:")
		fmt.Println(string(file))
	}

	casdoorsdk.InitConfig(c.AuthEndpoint, c.AuthClientID, c.AuthClientSecret, string(file), c.AuthOrganization, c.AuthApplicationName)
	return nil
}
