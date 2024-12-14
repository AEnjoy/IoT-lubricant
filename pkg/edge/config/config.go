package config

import (
	"errors"
	"os"
	"sync"

	"github.com/AEnjoy/IoT-lubricant/internal/model"
	def "github.com/AEnjoy/IoT-lubricant/pkg/default"
	"github.com/AEnjoy/IoT-lubricant/pkg/logger"
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	json "github.com/bytedance/sonic/encoder"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

var (
	Config       *model.EdgeSystem
	GatewayID    string
	GatherLock   sync.Mutex
	OriginConfig openapi.OpenAPICli
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
func NullConfig() *model.EdgeSystem {
	return &model.EdgeSystem{
		Config: &openapi.ApiInfo{},
	}
}
func SaveConfig(t SaveType) error {
	logger.Infoln("Saving config... with SaveType:", t)
	var errs error
	if t == SaveType_ALL || t == SaveType_Config {
		fileName := os.Getenv("CONFIG")
		if fileName == "" {
			fileName = def.AgentDefaultConfigFileName
		}
		if Config.FileName == "" {
			Config.FileName = def.AgentDefaultOpenapiFileName
		}
		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		err1 := yaml.NewEncoder(f).Encode(Config)
		errs = errors.Join(err1, err)

		diff := cmp.Diff(OriginConfig, Config.Config.(*openapi.ApiInfo).OpenAPICli)
		if diff != "" {
			logger.Info("Config changed:", diff)
			confFile, err := os.OpenFile(Config.FileName, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				errs = errors.Join(errs, err)
			} else {
				errs = errors.Join(errs, json.NewStreamEncoder(confFile).Encode(Config.Config.(*openapi.ApiInfo)))
			}
		}
	}

	if t == SaveType_ALL || t == SaveType_EnableConfig && Config.Config.GetEnable() != nil {
		fileName := Config.FileName + ".enable"
		enConfFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			errs = errors.Join(errs, json.NewStreamEncoder(enConfFile).Encode(Config.Config.GetEnable()))
		}
	}
	return errs
}
func RefreshSlot() error {
	if Config.Config != nil && Config.Config.GetEnable() != nil && Config.Config.GetEnable().Slot != nil {
		var slot []int
		for i := range Config.Config.GetEnable().Slot {
			slot = append(slot, i)
		}
		Config.EnableSlot = slot
		return nil
	} else {
		return errors.New("config error")
	}
}
