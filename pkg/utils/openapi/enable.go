package openapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AEnjoy/IoT-lubricant/pkg/types/errs"
)

type Enable struct {
	Get  map[string][]Parameter  `json:"get"`
	Post map[string]*RequestBody `json:"post"`
	Slot map[int]string          `json:"slot"` // slot_id:method-path
}

type EnableParams struct {
	// 不需要考虑同时设置get和post参数的情况
	GetParams  map[string]string `json:"get_params"` // key is request param-name and value is param-value
	PostParams *RequestBody      `json:"post_params"`
	Slot       int               `json:"slot"`
}

func (e *Enable) SlotGetEnable(slot int) (method, path string) {
	if strings.HasPrefix(e.Slot[slot], http.MethodGet) {
		return http.MethodGet, strings.TrimPrefix(e.Slot[slot], fmt.Sprintf("%s:", http.MethodGet))
	}
	if strings.HasPrefix(e.Slot[slot], http.MethodPost) {
		return http.MethodPost, strings.TrimPrefix(e.Slot[slot], fmt.Sprintf("%s:", http.MethodPost))
	}
	return "", ""
}

// EnableApi 启用api中指定的方法,并设置参数,返回处理后的api
//
// 如果不需要添加参数，params 仍然应该不为 nil，但里面的参数全为空即可
func EnableApi(doc OpenApi, params *EnableParams, path string) (OpenApi, error) {
	if params == nil {
		return nil, errs.ErrInvalidInput
	}
	item, ok := doc.GetPaths()[path]
	if !ok {
		return nil, errs.ErrInvalidPath
	}
	apiInfo := doc.(*ApiInfo)
	apiInfo.l.Lock()
	defer apiInfo.l.Unlock()

	if apiInfo.Slot == nil {
		apiInfo.Slot = make(map[int]string)
	}
	if apiInfo.Get == nil {
		apiInfo.Get = make(map[string][]Parameter)
	}
	if apiInfo.Post == nil {
		apiInfo.Post = make(map[string]*RequestBody)
	}

	if params.GetParams != nil && item.GetGet() != nil {
		s := int(params.Slot)
		_, ok := apiInfo.Enable.Slot[s]
		if ok {
			return nil, fmt.Errorf("slot %d is already used", s)
		}

		var parameters []Parameter
		for k, v := range params.GetParams {
			var p Parameter
			p.Set(k, v)
			parameters = append(parameters, p)
		}
		apiInfo.Enable.Get[path] = parameters
		apiInfo.Enable.Slot[s] = fmt.Sprintf("%s:%s", http.MethodGet, path)
	}

	if params.PostParams != nil && item.GetPost() != nil {
		s := int(params.Slot)
		_, ok := apiInfo.Enable.Slot[s]
		if ok {
			return nil, fmt.Errorf("slot %d is already used", s)
		}

		apiInfo.Enable.Post[path] = params.PostParams
		apiInfo.Enable.Slot[s] = fmt.Sprintf("%s:%s", http.MethodPost, path)
	}
	return apiInfo, nil
}

func DisableApi(doc OpenApi, slot int) (OpenApi, error) {
	apiInfo := doc.(*ApiInfo)
	apiInfo.l.Lock()
	defer apiInfo.l.Unlock()

	if apiInfo.Slot == nil {
		return nil, errs.ErrInvalidSlot
	}

	method, path := apiInfo.SlotGetEnable(slot)
	switch method {
	case http.MethodGet:
		delete(apiInfo.Enable.Get, path)
	case http.MethodPost:
		delete(apiInfo.Enable.Post, path)
	}
	delete(apiInfo.Enable.Slot, slot)
	return apiInfo, nil
}

// UpdateApi 更新api中指定的方法,并设置参数,返回处理后的api
//
// 与EnableApi类似但只能更新已启用的api和slot，,
// 如果slot = -1,api会自动寻找对应的slot，否则为覆盖 slot
func UpdateApi(doc OpenApi, params *EnableParams, path string) (OpenApi, error) {
	if params == nil {
		return nil, errs.ErrInvalidInput
	}
	item, ok := doc.GetPaths()[path]
	if !ok {
		return nil, errs.ErrInvalidPath
	}
	apiInfo := doc.(*ApiInfo)
	apiInfo.l.Lock()
	defer apiInfo.l.Unlock()

	if apiInfo.Slot == nil {
		return nil, errs.ErrNotInit
	}

	if params.GetParams != nil && item.GetGet() != nil {
		if params.Slot == -1 {
			for s, v := range apiInfo.Slot {
				if v == fmt.Sprintf("%s:%s", http.MethodGet, path) {
					params.Slot = s
					break
				}
			}
		}
		if params.Slot == -1 {
			return nil, errs.ErrInvalidSlot
		}

		var parameters []Parameter
		for k, v := range params.GetParams {
			var p Parameter
			p.Set(k, v)
			parameters = append(parameters, p)
		}
		apiInfo.Enable.Get[path] = parameters
		apiInfo.Enable.Slot[params.Slot] = fmt.Sprintf("%s:%s", http.MethodGet, path)
	}

	if params.PostParams != nil && item.GetPost() != nil {
		if params.Slot == -1 {
			for s, v := range apiInfo.Slot {
				if v == fmt.Sprintf("%s:%s", http.MethodPost, path) {
					params.Slot = s
					break
				}
			}
		}
		if params.Slot == -1 {
			return nil, errs.ErrInvalidSlot
		}

		apiInfo.Enable.Post[path] = params.PostParams
		apiInfo.Enable.Slot[params.Slot] = fmt.Sprintf("%s:%s", http.MethodPost, path)
	}
	return apiInfo, nil
}
