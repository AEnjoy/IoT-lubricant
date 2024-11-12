package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
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
func SaveConfig(isEnable bool) error {
	var errs error
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
	if isEnable && Config.EnableConfig != nil {
		if !strings.HasSuffix(Config.FileName, ".enable") {
			Config.FileName = Config.FileName + ".enable"
		}
		enConfFile, err := os.OpenFile(Config.FileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			errs = errors.Join(errs, json.NewEncoder(enConfFile).Encode(Config.EnableConfig.(*openapi.ApiInfo)))
		}
	}
	if !isEnable && Config.Config != nil {
		if strings.HasSuffix(Config.FileName, ".enable") {
			Config.FileName = strings.TrimSuffix(Config.FileName, ".enable")
		}
		confFile, err := os.OpenFile(Config.FileName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			errs = errors.Join(errs, err)
		} else {
			errs = errors.Join(errs, json.NewEncoder(confFile).Encode(Config.Config.(*openapi.ApiInfo)))
		}
	}
	return errs
}
