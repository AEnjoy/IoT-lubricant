package core

import (
	"github.com/AEnjoy/IoT-lubricant/internal/app/core/config"
	def "github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/file"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/hash"
	"os"
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
