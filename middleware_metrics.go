package monitoring

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type middleware struct {
	duration *prometheus.HistogramVec
}

func Config() middleware {

	vec := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time taken to process HTTP requests",
		Buckets: []float64{0.1, 0.15, 0.2, 0.25, 0.3},
	}, []string{"code", "method", "endpoint"})

	prometheus.MustRegister(vec)

	return middleware{
		duration: vec,
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
