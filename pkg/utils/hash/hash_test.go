package hash

import (
	"testing"
)

func TestMD5(t *testing.T) {
	passwd := "123456"
	password, err := HashPassword(passwd)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(password)
}
