package user

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/repo"
)

type Api struct {
	Db repo.ICoreDb
}
