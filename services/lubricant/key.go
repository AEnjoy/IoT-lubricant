package lubricant

import (
	"os"

	def "github.com/aenjoy/iot-lubricant/pkg/constant"
	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/aenjoy/iot-lubricant/pkg/utils/file"
	"github.com/aenjoy/iot-lubricant/pkg/utils/hash"
	"github.com/aenjoy/iot-lubricant/services/lubricant/config"
)

func initKeys() error {
	readFile, err := file.ReadFile(def.ServerKeyFileName)
	if err == nil {
		config.ServerKeys = []byte(readFile)
		return nil
	}
	key, err := hash.GenerateKey(32)
	if err != nil {
		logger.Errorln("GenerateKey error, error Info is: ", err)
		return err
	}
	config.ServerKeys = key
	return os.WriteFile(def.ServerKeyFileName, key, 0644)
}
