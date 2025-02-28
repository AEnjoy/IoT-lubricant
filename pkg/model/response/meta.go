package response

type Meta struct {
	Code int    `json:"http_code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Name string `json:"name,omitempty"` // module name
	Data any    `json:"data,omitempty"`
}
type Failed struct {
	Meta
}
type Success struct {
	Meta
}

type Empty struct {
}
