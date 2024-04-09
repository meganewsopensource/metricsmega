package monitoring

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Test_middleware_Metrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(Config().Metrics)
	r.GET("/metrics", handler())
	r.GET("/test", response)
	tests := []struct {
		name           string
		exceptedString string
		wantErr        bool
	}{
		{
			name:           "Success",
			exceptedString: "http_request_duration_seconds_bucket{code=\"200\",endpoint=\"/test\",method=\"GET\",le=\"0.1\"} 1",
			wantErr:        false,
		},
		{
			name:           "Fail, could not find that text",
			exceptedString: "test-123456878-just-for-break",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Error("Error on http.NewRequest")
			}
			response := httptest.NewRecorder()
			r.ServeHTTP(response, request)

			requestMetrics, err := http.NewRequest("GET", "/metrics", nil)
			if err != nil {
				t.Error("Error on http.NewRequest")
			}
			responseMetrics := httptest.NewRecorder()
			r.ServeHTTP(responseMetrics, requestMetrics)

			body, err := io.ReadAll(responseMetrics.Body)
			if err != nil {
				fmt.Printf("Erro ao ler o corpo da resposta: %s\n", err)
				return
			}

			bodyStr := string(body)

			if !tt.wantErr && !strings.Contains(bodyStr, tt.exceptedString) {
				t.Errorf("Error on CompanyLinkerController: expected status code: %s, returned status code: %s", tt.exceptedString, bodyStr)
			}
		})
	}
}

func response(context *gin.Context) {
	context.JSON(http.StatusOK, "It works!")
}

func handler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
