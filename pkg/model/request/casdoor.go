package request

type CasdoorPasswordAuthRequest struct {
	Application  string `json:"application"`
	Organization string `json:"organization"`
	Username     string `json:"username"`
	AutoSignin   bool   `json:"autoSignin"`
	Password     string `json:"password"`
	SigninMethod string `json:"signinMethod"`
	Type         string `json:"type"`
}
