package middleware

import (
	"runtime/debug"

	"github.com/AEnjoy/IoT-lubricant/pkg/utils/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetRecovery(reg *prometheus.Registry) grpc.UnaryServerInterceptor {
	// Setup metric for panic recoveries.
	panicsTotal := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})

	grpcPanicRecoveryHandler := func(p any) (err error) {
		panicsTotal.Inc()
		logger.Errorf("recovered from internal panic: %v Stack:%s", p, string(debug.Stack()))
		return status.Errorf(codes.Internal, "%s", p)
	}
	return recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler))
}
