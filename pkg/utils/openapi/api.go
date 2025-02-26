package openapi

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/aenjoy/iot-lubricant/pkg/types/errs"
	"github.com/aenjoy/iot-lubricant/pkg/utils/file"
	json "github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
)

type ApiInfo struct {
	filename string
	l        *sync.Mutex

	OpenAPICli `json:"open_api_cli"`
	Enable     `json:"enable"`
}

// 定义OpenAPI文档的结构体
type OpenAPICli struct {
	OpenAPI string              `json:"openapi"`
	Info    Info                `json:"info"`
	Servers []Server            `json:"servers,omitempty"`
	Paths   map[string]PathItem `json:"paths"`
}

func (api *OpenAPICli) GetApiInfo() Info {
	return api.Info
}
func (api *OpenAPICli) GetPaths() map[string]PathItem {
	if api.Paths == nil {
		api.Paths = make(map[string]PathItem)
	}
	return api.Paths
}

// Server 结构体定义了服务器的URL和描述
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// PathItem 包含多个HTTP方法的操作
type PathItem struct {
	Get  *Operation `json:"get,omitempty"`
	Post *Operation `json:"post,omitempty"`
}

func (api *PathItem) GetGet() *Operation {
	return api.Get
}
func (api *PathItem) GetPost() *Operation {
	return api.Post
}

// Operation 定义了操作的详细信息，包括请求和响应
type Operation struct {
	Summary     string              `json:"summary"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"` // POST
	Parameters  []Parameter         `json:"parameters,omitempty"`  // GET
	Responses   map[string]Response `json:"responses"`             // 200, 400等
}

func (api *Operation) GetSummary() string {
	return api.Summary
}
func (api *Operation) GetRequestBody() *RequestBody {
	return api.RequestBody
}
func (api *Operation) GetParameters() []Parameter {
	return api.Parameters
}
func (api *Operation) GetResponses() map[string]Response {
	return api.Responses
}

// RequestBody 定义了POST请求的body结构
type RequestBody struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content"`
}

func (api *RequestBody) SetBodyWithJson(kv map[string]string) {
	if kv == nil {
		return
	}
	api.Content["application/json"] = MediaType{
		Schema: kv,
	}
}

func (api *RequestBody) GetDescription() string {
	return api.Description
}
func (api *RequestBody) GetContent() map[string]MediaType {
	return api.Content
}

// MediaType 定义了请求内容的结构
type MediaType struct {
	Schema interface{} `json:"schema"` // 使用interface{}来表示未知的内容
}

func (api *MediaType) GetSchema() interface{} {
	return api
}
func (api *MediaType) GetSchemaContent() ([]byte, error) {
	return json.Marshal(api)
}

// Parameter 定义了GET请求的参数
type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"` // 通常是"query"、"header"等
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Schema      Schema `json:"schema"`
}

func (p *Parameter) Set(k, v string) {
	p.Name = k
	p.Schema = Schema{
		Type:       "string",
		Properties: map[string]Property{p.Name: {Type: v}},
	}
}

// Schema 定义了JSON请求body的schema
type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
}

func (api *Schema) GetProperties() map[string]Property {
	return api.Properties
}

// Property 定义了body中字段的类型
type Property struct {
	Type string `json:"type"`
}

// Response 定义了接口的响应
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content"` //key: such as application/json
}

func (api *Response) GetContent() map[string]MediaType {
	return api.Content
}

func (api *ApiInfo) InitApis(filename string) error {
	api.filename = filename
	api.OpenAPICli = OpenAPICli{}

	if filename != "" {
		err := file.ReadJsonFile(filename, &api.OpenAPICli)
		if err != nil {
			panic(err)
		}
	} else {
		api.OpenAPICli.Paths = make(map[string]PathItem)
	}

	return nil
}
func (api *ApiInfo) SendGETMethod(path string, parameters []Parameter) ([]byte, error) {
	if _, ok := api.Paths[path]; !ok {
		return nil, errs.ErrNotFound
	}
	if api.Paths[path].Get == nil {
		return nil, errs.ErrInvalidMethod
	}

	fullPath := api.Servers[0].URL + path //
	cli := http.Client{}
	q := url.Values{}
	for _, param := range parameters {
		q.Add(param.Name, param.Schema.GetProperties()[param.Name].Type)
	}
	u, err := url.Parse(fullPath)
	if err != nil {
		return nil, err
	}
	u.RawQuery = q.Encode()
	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	respHttp, err := cli.Do(request)
	if err != nil {
		return nil, err
	}
	defer respHttp.Body.Close()
	all, err := io.ReadAll(respHttp.Body)
	if err != nil {
		return nil, err
	}
	//logger.Infoln(string(all))
	//resp := make(map[string]Response)
	//content := make(map[string]MediaType)
	//
	//for k, _ := range api.Paths[path].Get.Responses["200"].Content {
	//	content[k] = MediaType{all}
	//}
	//resp[respHttp.Status] = Response{
	//	Description: api.Paths[path].Get.Responses[respHttp.Status].Description,
	//	Content:     content,
	//}
	return all, nil
}
func (api *ApiInfo) SendPOSTMethod(path string, body RequestBody) ([]byte, error) {
	if _, ok := api.Paths[path]; !ok {
		return nil, errs.ErrNotFound
	}
	if api.Paths[path].Post == nil {
		return nil, errs.ErrInvalidMethod
	}
	fullPath := api.Servers[0].URL + path
	cli := http.Client{}

	var ct string
	for c := range body.GetContent() {
		ct = c
	}

	byte, err := json.Marshal(body.GetContent()[ct].Schema)
	if err != nil {
		return nil, err
	}

	bytesReader := bytes.NewReader(byte)
	request, err := http.NewRequest(http.MethodPost, fullPath, bytesReader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", ct)

	respHttp, err := cli.Do(request)

	if err != nil {
		return nil, err
	}
	defer respHttp.Body.Close()

	all, err := io.ReadAll(respHttp.Body)
	if err != nil {
		return nil, err
	}
	//
	//resp := make(map[string]Response)
	//content := make(map[string]MediaType)
	//
	//for k, _ := range api.Paths[path].Get.Responses["200"].Content {
	//	content[k] = MediaType{all}
	//}
	//resp[respHttp.Status] = Response{
	//	Description: api.Paths[path].Get.Responses[respHttp.Status].Description,
	//	Content:     content,
	//}
	return all, nil
}
func (api *ApiInfo) SendPOSTMethodEx(path, ct string, body []byte) ([]byte, error) {
	if _, ok := api.Paths[path]; !ok {
		return nil, errs.ErrNotFound
	}
	if api.Paths[path].Post == nil {
		return nil, errs.ErrInvalidMethod
	}

	bytesReader := bytes.NewReader(body)
	request, err := http.NewRequest(http.MethodPost, path, bytesReader)
	if err != nil {
		return nil, err
	}

	cli := http.Client{}
	request.Header.Set("Content-Type", ct)

	respHttp, err := cli.Do(request)

	if err != nil {
		return nil, err
	}
	defer respHttp.Body.Close()

	all, err := io.ReadAll(respHttp.Body)
	if err != nil {
		return nil, err
	}
	//
	//resp := make(map[string]Response)
	//content := make(map[string]MediaType)
	//
	//for k, _ := range api.Paths[path].Get.Responses["200"].Content {
	//	content[k] = MediaType{all}
	//}
	//resp[respHttp.Status] = Response{
	//	Description: api.Paths[path].Get.Responses[respHttp.Status].Description,
	//	Content:     content,
	//}
	return all, nil
}
func (api *ApiInfo) GetEnable() *Enable {
	return &api.Enable
}
func (api *ApiInfo) InitEnable(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return decoder.NewStreamDecoder(file).Decode(&api.Enable)
}
