package edge

import "github.com/AEnjoy/IoT-lubricant/pkg/utils/openapi"

func CheckConfigInvalidGet(a openapi.OpenApi) bool {
	if a == nil {
		return false
	}
	// 检查至少一个选项启用且配置有效
	for _, item := range a.GetPaths() {
		opera := item.GetGet()
		if item.GetPost() != nil && opera == nil { // POST
			continue
		}

		if opera == nil {
			return false
		}
		parameters := opera.GetParameters()
		for _, param := range parameters {
			t := param.Schema.GetProperties()[param.Name].Type
			if t == "" && param.Required {
				return false
			}
		}
	}
	return true
}
