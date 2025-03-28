package edge

import (
	"github.com/aenjoy/iot-lubricant/pkg/utils/openapi"
	"github.com/bytedance/sonic"
	"github.com/getkin/kin-openapi/openapi3"
)

func CheckConfigInvalidGet(a openapi.OpenApi) bool {
	if a == nil {
		return false
	}
	// 检查至少一个选项启用且配置有效
	apiInfo := a.(*openapi.ApiInfo)
	marshal, err := sonic.Marshal(apiInfo.OpenAPICli)
	if err != nil {
		return false
	}
	if !IsOpenAPIDoc(marshal) {
		return false
	}
	return len(a.GetEnable().Get) != 0
}
func CheckConfigInvalidPost(a openapi.OpenApi) bool {
	if a == nil {
		return false
	}
	// 检查至少一个选项启用且配置有效
	apiInfo := a.(*openapi.ApiInfo)
	marshal, err := sonic.Marshal(apiInfo.OpenAPICli)
	if err != nil {
		return false
	}
	if !IsOpenAPIDoc(marshal) {
		return false
	}
	return len(a.GetEnable().Post) != 0
}
func CheckConfigInvalid(a openapi.OpenApi) bool {
	return CheckConfigInvalidGet(a) || CheckConfigInvalidPost(a)
}
func IsOpenAPIDoc(data []byte) bool {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return false
	}
	return doc.Validate(loader.Context, openapi3.DisableExamplesValidation()) == nil
}
