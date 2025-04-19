package middleware

import (
	"time"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
	"github.com/gin-gonic/gin"
)

func GetMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		AllowCORS,
		CustomRecoveryWithZap(logger.L(), true, defaultHandleRecovery),
		Ginzap(logger.L(), time.RFC3339, true),
	}
}
