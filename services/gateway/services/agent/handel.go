package agent

import "github.com/aenjoy/iot-lubricant/pkg/types/exception"

var handelFunc func(<-chan *exception.Exception)

func SetErrorHandelFunc(f func(<-chan *exception.Exception)) {
	handelFunc = f
}
