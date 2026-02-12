package monitoring

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type grpcMiddleware struct {
	duration *prometheus.HistogramVec
}

var (
	grpcDuration *prometheus.HistogramVec
	grpcOnce     sync.Once
)

func GrpcConfig() grpcMiddleware {
	grpcOnce.Do(func() {
		grpcDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Time taken to process gRPC requests",
			Buckets: []float64{0.001, 0.002, 0.004, 0.008, 0.016, 0.032, 0.064, 0.128, 0.256, 0.512, 1.024, 2.048, 4.096, 8.192, 16.384, 32.768},
		}, []string{"code", "method", "endpoint"})

		prometheus.MustRegister(grpcDuration)
	})

	return grpcMiddleware{
		duration: grpcDuration,
	}
}

func (m grpcMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()

		resp, err := handler(ctx, req)

		elapsed := time.Since(startTime).Seconds()
		st, _ := status.FromError(err)

		m.duration.With(prometheus.Labels{
			"code":     strconv.Itoa(int(st.Code())),
			"method":   "unary", // You might want to extract the method name more specifically if needed, or stick to this convention
			"endpoint": info.FullMethod,
		}).Observe(elapsed)

		return resp, err
	}
}
