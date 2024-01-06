package openapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const path = "../../../test/mock_driver/clock/api.json"

func TestApiCli(t *testing.T) {
	t.Skip("E2E是成功的,在单元测试还需要做些修改")
	// todo:需要将httptestserver的端口设置到 path.json-server-url中
	assert := assert.New(t)
	cli, err := NewOpenApiCli(path)
	assert.NoError(err)
	t.Log("Show Api info")
	paths := cli.GetPaths()
	for path, item := range paths {
		if item.Get != nil {
			t.Log("GET:", path, item.GetGet().GetSummary())
			item, ok := cli.Paths[path]
			assert.True(ok)
			parm := item.GetGet().Parameters
			if len(parm) == 1 {
				assert.Equal("time", parm[0].Name)
				m := make(map[string]Property)
				m["time"] = Property{Type: "2024-01-01 20:37:17"}
				parm[0].Schema.Properties = m
			}
			resp, err := cli.SendGETMethod(path, parm)
			assert.NoError(err)
			t.Log(string(resp))
		}
		if item.Post != nil {
			t.Log("POST:", path, item.GetPost().GetSummary())
			body := item.GetPost().GetRequestBody()
			c := make(map[string]MediaType)
			kv := map[string]interface{}{"time": "2024-01-01 20:37:17"}
			body.Content = c
			body.Content["application/json"] = MediaType{kv}
			resp, err := cli.SendPOSTMethod(path, *body)
			assert.NoError(err)
			t.Log(string(resp))
		}
	}

}
