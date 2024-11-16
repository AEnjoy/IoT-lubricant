package edge

import (
	"errors"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidPath  = errors.New("invalid path")
)

type Params struct {
	// 不需要考虑同时设置get和post参数的情况
	GetParams  *map[string]string   `json:"get_params"`
	PostParams *openapi.RequestBody `json:"post_params"`
}

// EnableApi 启用api中指定的方法,并设置参数,返回处理后的api
func EnableApi(doc openapi.ApiInfo, params *Params, path string) (openapi.ApiInfo, error) {
	if params == nil {
		return openapi.ApiInfo{}, ErrInvalidInput
	}
	pathItem, ok := doc.GetPaths()[path]
	if !ok {
		return openapi.ApiInfo{}, ErrInvalidPath
	}
	if params.GetParams != nil {
		var parameters []openapi.Parameter
		for k, v := range *params.GetParams {
			var p openapi.Parameter
			p.Set(k, v)
			parameters = append(parameters, p)
		}
		pathItem.Get.Parameters = parameters
	}
	if params.PostParams != nil {
		pathItem.Post.RequestBody = params.PostParams
	}
	return doc, nil
}
