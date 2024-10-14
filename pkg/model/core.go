package model

type Server struct {
	Info map[string]string `json:"server_info" gorm:"column:server_info;serializer:json"`
}

func (Server) TableName() string {
	return "server"
}

type Token struct {
	// 该Token是颁发
	UserId string `json:"id" gorm:"column:user_id"` // uuid
	// 办法给用户的访问令牌(用户需要携带Token来访问接口)
	AccessToken string `json:"access_token" gorm:"column:access_token"`
	// 过期时间(2h), 单位是秒
	AccessTokenExpiredAt int `json:"access_token_expired_at" gorm:"column:access_token_expired_at"`
	// 刷新Token
	RefreshToken string `json:"refresh_token" gorm:"column:refresh_token"`
	// 刷新Token过期时间(7d)
	RefreshTokenExpiredAt int `json:"refresh_token_expired_at" gorm:"column:refresh_token_expired_at"`

	// 创建时间
	CreatedAt int64 `json:"created_at" gorm:"column:created_at"`
	// 更新实现
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at"`
}

func (Token) TableName() string {
	return "token"
}

type User struct {
	UserId   string `json:"id" gorm:"column:user_id"` // uuid
	UserName string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`

	CreatedAt int64 `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "user"
}

type Data struct {
	AgentID  string `json:"agent_id" gorm:"column:agent_id"`
	DeviceID string `json:"device_id" gorm:"column:device_id"`

	Content string `json:"data" gorm:"column:data;serializer:json"`

	CreatedAt int64 `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at"`
}

func (Data) TableName() string {
	return "data"
}

type Gateway struct {
	GatewayID   string `json:"id" gorm:"column:id"`
	UserId      string `json:"user_id" gorm:"column:user_id"`
	Description string `json:"description" gorm:"column:description"`

	Address           string `json:"address" gorm:"column:address"`                         // SSH: ip:port or domain:port
	UserNameAndPasswd string `json:"username_and_passwd" gorm:"column:username_and_passwd"` //

	TlsConfig string `json:"tls_config" gorm:"column:tls_config;serializer:json"` // grpc tls config

	CreatedAt int64 `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at"`
}

func (Gateway) TableName() string {
	return "gateway"
}
