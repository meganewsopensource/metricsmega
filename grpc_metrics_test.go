package monitoring

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_grpcMiddleware_UnaryServerInterceptor(t *testing.T) {
	// Initialize the metrics
	m := GrpcConfig()

	// reset the metric for testing purposes if needed, but since we use a global registry,
	// we should be careful. NewHistogramVec registers to DefaultRegisterer.
	// For this test, we can just check if the metric count increases.

	tests := []struct {
		name           string
		handler        grpc.UnaryHandler
		info           *grpc.UnaryServerInfo
		expectedCode   string
		expectedMethod string
		expectedPath   string
		expectError    bool
	}{
		{
			name: "Success",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "response", nil
			},
			info: &grpc.UnaryServerInfo{
				FullMethod: "/package.Service/Method",
			},
			expectedCode:   "0", // OK
			expectedMethod: "unary",
			expectedPath:   "/package.Service/Method",
			expectError:    false,
		},
		{
			name: "Error",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, status.Error(codes.InvalidArgument, "invalid argument")
			},
			info: &grpc.UnaryServerInfo{
				FullMethod: "/package.Service/MethodError",
			},
			expectedCode:   "3", // InvalidArgument
			expectedMethod: "unary",
			expectedPath:   "/package.Service/MethodError",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// We need to measure the count before.
			// Construct the expected label values to query the metric.
			// Note: testutil.ToFloat64 might be hard with Labels.
			// Let's use Gather() and check the output string.

			interceptor := m.UnaryServerInterceptor()
			_, err := interceptor(context.Background(), nil, tt.info, tt.handler)

			if (err != nil) != tt.expectError {
				t.Errorf("UnaryServerInterceptor() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// Allow some time for metric observation (though it should be synchronous)
			time.Sleep(10 * time.Millisecond)

			// Validate metrics
			// We expect: grpc_request_duration_seconds_bucket{code="...",endpoint="...",method="unary",le="..."}

			// Since we cannot easily isolate this test's metrics from others in the same process
			// (if run in parallel or shared registry), we check if the metric exists.

			// We can use testutil.CollectAndCompare but we need the expected output.
			// Or just check if the metric families contain our metric with correct labels.

			// A simpler approach for this unit test:
			// 1. Create a dummy registry? No, the code uses prometheus.MustRegister which uses DefaultRegisterer.
			//    Refactoring to allow injecting registerer would be better but changes API.

			// Let's inspect the DefaultGatherer.
			mfs, err := prometheus.DefaultGatherer.Gather()
			if err != nil {
				t.Fatalf("Gather failed: %v", err)
			}

			found := false
			for _, mf := range mfs {
				if mf.GetName() == "grpc_request_duration_seconds" {
					for _, m := range mf.GetMetric() {
						foundCode := false
						foundEndpoint := false
						foundMethod := false

						for _, label := range m.GetLabel() {
							if label.GetName() == "code" && label.GetValue() == tt.expectedCode {
								foundCode = true
							}
							if label.GetName() == "endpoint" && label.GetValue() == tt.expectedPath {
								foundEndpoint = true
							}
							if label.GetName() == "method" && label.GetValue() == tt.expectedMethod {
								foundMethod = true
							}
						}

						if foundCode && foundEndpoint && foundMethod {
							found = true
							break
						}
					}
				}
			}

			if !found {
				t.Errorf("Metric not found for %s", tt.name)
			}
		})
	}
}
