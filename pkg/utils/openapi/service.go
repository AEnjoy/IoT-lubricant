package openapi

var _ OpenApi = (*ApiInfo)(nil)

type OpenApi interface {
	SendGETMethod(path string, parameters []Parameter) ([]byte, error)
	SendPOSTMethod(path string, body RequestBody) ([]byte, error)
	SendPOSTMethodEx(path, ct string, body []byte) ([]byte, error)

	GetApiInfo() Info
	GetPaths() map[string]PathItem
}

func NewOpenApiCli(fileName string) (*ApiInfo, error) {
	retVal := &ApiInfo{}
	err := retVal.InitApis(fileName)
	if err != nil {
		return nil, err
	}
	return retVal, err
}
