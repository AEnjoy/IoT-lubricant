package openapi

import (
	"sync"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/file"
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
	}
	err = retVal.InitEnable(fileName)
	if err != nil {
		return nil, err
	}
	return retVal, err
}
