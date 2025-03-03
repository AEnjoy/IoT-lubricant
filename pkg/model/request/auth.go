package request

type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
	Membered bool   `json:"membered"`
}
type Token struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
