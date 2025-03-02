package model

import (
	"errors"
	"time"

	"github.com/rs/xid"
)

const (
	COOKIE_TOKEY_KEY   = "token"
	REFRESH_HEADER_KEY = "X-REFRESH-TOKEN"
	GIN_TOKEN_KEY_NAME = "token"
)

type Token struct {
	// 该Token是颁发
	UserId string `json:"id" gorm:"column:user_id"` // uuid
	// 办法给用户的访问令牌(用户需要携带Token来访问接口)
	AccessToken string `json:"token" gorm:"column:access_token"`
	// 过期时间(2h), 单位是秒
	AccessTokenExpiredAt int `json:"access_token_expired_at" gorm:"column:access_token_expired_at"`
	// 刷新Token
	RefreshToken string `json:"refreshToken" gorm:"column:refresh_token"`
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
func (t *Token) IssueTime() time.Time {
	return time.Unix(t.CreatedAt, 0)
}

func (t *Token) AccessTokenDuration() time.Duration {
	return time.Duration(t.AccessTokenExpiredAt * int(time.Second))
}

func (t *Token) RefreshTokenDuration() time.Duration {
	return time.Duration(t.RefreshTokenExpiredAt * int(time.Second))
}

func (t *Token) AccessTokenIsExpired() error {
	// 过期时间: 颁发时间 + 过期时长
	expiredTime := t.IssueTime().Add(t.AccessTokenDuration())

	if time.Since(expiredTime).Seconds() > 0 {
		return errors.New("AccessTokenExpired")
	}

	return nil
}

func (t *Token) RefreshTokenIsExpired() error {
	// 过期时间: 颁发时间 + 过期时长
	expiredTime := t.IssueTime().Add(t.RefreshTokenDuration())

	if time.Since(expiredTime).Seconds() > 0 {
		return errors.New("RefreshTokenExpired")
	}
	return nil
}
func NewToken(u *User) *Token {
	if u == nil {
		return &Token{}
	}

	return &Token{
		UserId: u.UserId,
		// 使用随机UUID
		AccessToken:           xid.New().String(),
		AccessTokenExpiredAt:  3600,
		RefreshToken:          xid.New().String(),
		RefreshTokenExpiredAt: 3600 * 4,
		CreatedAt:             time.Now().Unix(),
	}
}
