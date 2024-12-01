package edge

import (
	"encoding/json"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
	"github.com/getkin/kin-openapi/openapi3"
)

func CheckConfigInvalidGet(a openapi.OpenApi) bool {
	if a == nil {
		return false
	}
	// 检查至少一个选项启用且配置有效
	apiInfo := a.(*openapi.ApiInfo)
	marshal, err := json.Marshal(apiInfo.OpenAPICli)
	if err != nil {
		return false
	}
	if !IsOpenAPIDoc(marshal) {
		return false
	}
	enable := a.GetEnable()
	if enable.Get == nil || len(enable.Get) == 0 {
		return false
	}
	return true
}
func CheckConfigInvalidPost(a openapi.OpenApi) bool {
	if a == nil {
		return false
	}
	// 检查至少一个选项启用且配置有效
	apiInfo := a.(*openapi.ApiInfo)
	marshal, err := json.Marshal(apiInfo.OpenAPICli)
	if err != nil {
		return false
	}
	if !IsOpenAPIDoc(marshal) {
		return false
	}
	enable := a.GetEnable()
	if enable.Post == nil || len(enable.Post) == 0 {
		return false
	}
	return true
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
	err = doc.Validate(loader.Context, openapi3.DisableExamplesValidation())
	return err == nil
}
