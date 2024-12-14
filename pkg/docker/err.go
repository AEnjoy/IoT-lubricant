package docker

import "errors"

var (
	ErrNotInit = errors.New("docker api not init")
	ErrPullWay = errors.New("pull way not support")
)
