package request

import (
	"github.com/aenjoy/iot-lubricant/pkg/utils/hash"
)

type CreateUserRequest struct {
	UserName string            `json:"username" gorm:"column:username"`
	Password string            `json:"password" gorm:"column:password"`
	Label    map[string]string `json:"label" gorm:"column:label;serializer:json"`
	hashed   bool
}

func (req *CreateUserRequest) HashPassword(key []byte) (err error) {
	if req.hashed {
		return nil
	}
	// hash the password
	e, err := hash.Encrypt(key, []byte(req.Password))
	if err != nil {
		return err
	}
	req.Password = string(e)
	req.hashed = true
	return
}
