package response

type CasdoorLoginResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Sub    string `json:"sub"`
	Name   string `json:"name"`
	Data   string `json:"data"`
	Data2  bool   `json:"data2"`
}
type Token struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
