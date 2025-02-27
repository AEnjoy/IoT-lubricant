package router

import (
	"github.com/aenjoy/iot-lubricant/services/lubricant/ioc"
	"github.com/gin-gonic/gin"
)

var _ ioc.Object = (*WebService)(nil)

type WebService struct {
	*gin.Engine
}

func (w *WebService) Init() error {
	router, err := CoreRouter()
	w.Engine = router
	return err
}

func (WebService) Weight() uint16 {
	return ioc.CoreWebServer
}

func (WebService) Version() string {
	return "dev"
}
