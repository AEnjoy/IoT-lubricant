package middleware

import (
	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

func GetLoggerInterceptor() grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(logger.InterceptorLogger())
}
