package openapi

import (
	"sync"

	"github.com/aenjoy/iot-lubricant/pkg/utils/file"
	"github.com/bytedance/sonic"
)

var _ OpenApi = (*ApiInfo)(nil)

type OpenApi interface {
	SendGETMethod(path string, parameters []Parameter) ([]byte, error)
	SendPOSTMethod(path string, body RequestBody) ([]byte, error)
	SendPOSTMethodEx(path, ct string, body []byte) ([]byte, error)

	GetApiInfo() Info
	GetPaths() map[string]PathItem
	GetEnable() *Enable
}

func NewOpenApiCli(fileName string) (*ApiInfo, error) {
	retVal := &ApiInfo{l: &sync.Mutex{}}
	err := retVal.InitApis(fileName)
	if err != nil {
		return nil, err
	}
	if file.IsFileExists(fileName + ".enable") {
		fileName = fileName + ".enable"
		err = retVal.InitEnable(fileName)
		if err != nil {
			return nil, err
		}
	}

	return retVal, err
}

func NewOpenApiCliEx(apiJson, enableJson []byte) (*ApiInfo, error) {
	retVal := &ApiInfo{l: &sync.Mutex{}}
	if len(apiJson) != 0 {
		if err := sonic.Unmarshal(apiJson, &retVal.OpenAPICli); err != nil {
			return nil, err
		}
	}
	if len(enableJson) != 0 {
		if err := sonic.Unmarshal(enableJson, &retVal.Enable); err != nil {
			return nil, err
		}
	}
	return retVal, nil
}
