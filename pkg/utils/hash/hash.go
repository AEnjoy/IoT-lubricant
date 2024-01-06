package hash

import (
	"crypto/md5"
	"fmt"
)

func HashPassword(str string) (string, error) {
	data := []byte(str + "-salt-123456789")
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has), nil
}
