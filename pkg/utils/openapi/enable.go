package openapi

import (
	"errors"

	"github.com/AEnjoy/IoT-lubricant/protobuf/meta"
)

type Enable struct {
	Get  map[string][]Parameter  `json:"get"`
	Post map[string]*RequestBody `json:"post"`
}

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidPath  = errors.New("invalid path")
)

type EnableParams struct {
	// 不需要考虑同时设置get和post参数的情况
	GetParams  map[string]string `json:"get_params"` // key is request param-name and value is param-value
	PostParams *RequestBody      `json:"post_params"`
	Slot       []meta.KvInt      `json:"slot"`
}

// EnableApi 启用api中指定的方法,并设置参数,返回处理后的api
//
// 如果不需要添加参数，params 仍然应该不为 nil，但里面的参数全为空即可
func EnableApi(doc OpenApi, params *EnableParams, path string) (OpenApi, error) {
	if params == nil {
		return nil, ErrInvalidInput
	}
	item, ok := doc.GetPaths()[path]
	if !ok {
		return nil, ErrInvalidPath
	}
	apiInfo := doc.(*ApiInfo)
	apiInfo.l.Lock()
	defer apiInfo.l.Unlock()

	if params.GetParams != nil && item.GetGet() != nil {
		var parameters []Parameter
		for k, v := range params.GetParams {
			var p Parameter
			p.Set(k, v)
			parameters = append(parameters, p)
		}
		apiInfo.Enable.Get[path] = parameters
	}

	if params.PostParams != nil && item.GetPost() != nil {
		apiInfo.Enable.Post[path] = params.PostParams
	}
	return apiInfo, nil
}
