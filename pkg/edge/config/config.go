package config

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	def "github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/types"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	"gopkg.in/yaml.v3"
)

var (
	Config     *types.EdgeSystem
	GatewayID  string
	GatherLock sync.Mutex
)

type SaveType uint8

const (
	SaveType_ALL SaveType = iota
	SaveType_EnableConfig
	SaveType_Config
)

func IsGathering() bool {
	if GatherLock.TryLock() {
		GatherLock.Unlock()
		return false
	}
	return true
}
func NullConfig() *types.EdgeSystem {
	return &types.EdgeSystem{
		Config:       &openapi.ApiInfo{},
		EnableConfig: &openapi.ApiInfo{},
	}
}
func SaveConfig(t SaveType) error {
	var errs error
	if t == SaveType_ALL || t == SaveType_Config {
		fileName := os.Getenv("CONFIG")
		if fileName == "" {
			fileName = def.AgentDefaultConfigFileName
		}

		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		err1 := yaml.NewEncoder(f).Encode(Config)
		errs = errors.Join(err1, err)

		confFile, err := os.OpenFile(Config.FileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			errs = errors.Join(errs, json.NewEncoder(confFile).Encode(Config.Config.(*openapi.ApiInfo)))
		}
	}

	if t == SaveType_EnableConfig && Config.Config.GetEnable() != nil {
		fileName := Config.FileName + ".enable"
		enConfFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			errs = errors.Join(errs, json.NewEncoder(enConfFile).Encode(Config.Config.GetEnable()))
		}
	}
	return errs
}
