package monitoring

import (
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type middleware struct {
	duration *prometheus.HistogramVec
}

var (
	httpDuration *prometheus.HistogramVec
	once         sync.Once
)

func Config() middleware {
	once.Do(func() {
		httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Time taken to process HTTP requests",
			Buckets: []float64{0.001, 0.002, 0.004, 0.008, 0.016, 0.032, 0.064, 0.128, 0.256, 0.512, 1.024, 2.048, 4.096, 8.192, 16.384, 32.768},
		}, []string{"code", "method", "endpoint"})

		prometheus.MustRegister(httpDuration)
	})

	return middleware{
		duration: httpDuration,
	}

}

func (m middleware) Metrics(c *gin.Context) {

	startTime := time.Now()

	c.Next()

	elapsed := time.Since(startTime).Seconds()

	m.duration.With(prometheus.Labels{"code": strconv.Itoa(c.Writer.Status()),
		"method":   c.Request.Method,
		"endpoint": c.Request.URL.Path,
	}).Observe(elapsed)

}
