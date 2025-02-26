package middleware

import (
	gprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

func GetSrvMetrics() *gprom.ServerMetrics {
	srvMetrics := gprom.NewServerMetrics(
		gprom.WithServerHandlingTimeHistogram(
			gprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)
	return srvMetrics
}
func GetRegistry(srvMetrics *gprom.ServerMetrics) *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)
	return reg
}
