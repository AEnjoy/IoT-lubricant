package lubricant

import (
	"errors"
	"fmt"
	"os"

	def "github.com/aenjoy/iot-lubricant/pkg/default"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

func initCasdoor() error {
	clientid, ok := os.LookupEnv(def.ENV_CORE_AUTH_CLIENT_ID)
	if !ok {
		return errors.New("missing client id")
	}
	secret, ok := os.LookupEnv(def.ENV_CORE_AUTH_CLIENT_SECRET)
	if !ok {
		return errors.New("missing client secret")
	}
	endpoint, ok := os.LookupEnv(def.ENV_CORE_AUTH_ENDPOINT)
	if !ok {
		return errors.New("missing endpoint")
	}
	organization, ok := os.LookupEnv(def.ENV_CORE_AUTH_ORGANIZATION)
	if !ok {
		return errors.New("missing organization")
	}
	publickeyfile, ok := os.LookupEnv(def.ENV_CORE_AUTH_PUBLICKEYFILE)
	if !ok {
		return errors.New("missing public key file")
	}
	file, err := os.ReadFile(publickeyfile)
	if err != nil {
		return err
	}
	if os.Getenv("RUNNING_LEVEL") == "debug" {
		logger.Debug("CERT FILE:")
		fmt.Println(string(file))
	}
	casdoorsdk.InitConfig(endpoint, clientid, secret, string(file), organization, "application_lubricant")
	return nil
}
