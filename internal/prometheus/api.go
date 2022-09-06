package prometheus

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (p *Prometheus) RegisterAPIMetrics() {
	p.totalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "http_requests_total"),
			Help: "Total number of HTTP requests.",
		},
		[]string{labelPath},
	)
	p.responseStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "http_response_status"),
			Help: "Status of HTTP response",
		},
		[]string{labelStatus},
	)
	p.responseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "http_response_time_seconds"),
			Help: "Duration of HTTP requests.",
		},
		[]string{labelPath},
	)
}

func GinMetricsMiddleware(prometheus *Prometheus) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.String()
		for _, param := range c.Params {
			path = strings.ReplaceAll(path, param.Value, ":"+param.Key)
		}

		timer := prometheus.SetReponseTime(path)
		c.Next()
		timer.ObserveDuration()

		prometheus.SetReponseStatus(c.Writer.Status())
		prometheus.SetTotalRequests(path)
	}
}

func (p *Prometheus) SetTotalRequests(path string) {
	p.totalRequests.WithLabelValues(path).Inc()
}

func (p *Prometheus) SetReponseStatus(status int) {
	p.responseStatus.WithLabelValues(strconv.Itoa(status)).Inc()
}

func (p *Prometheus) SetReponseTime(path string) *prometheus.Timer {
	return prometheus.NewTimer(p.responseTime.WithLabelValues(path))
}
