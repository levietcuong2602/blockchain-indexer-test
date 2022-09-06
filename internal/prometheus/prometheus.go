package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const (
	labelChain     = "chain"
	labelPath      = "path"
	labelStatus    = "status"
	labelHost      = "host"
	labelVersion   = "version"
	labelTopic     = "topic"
	labelPartition = "partition"
	labelOffset    = "offset"
)

// Prometheus is a struct for prometheus metrics.
type Prometheus struct {
	namespace string
	subsystem string

	// API metrics
	totalRequests  *prometheus.CounterVec
	responseStatus *prometheus.CounterVec
	responseTime   *prometheus.HistogramVec

	// Blockproducer metrics
	lastFetchedBlock      *prometheus.GaugeVec
	currentNodeBlock      *prometheus.GaugeVec
	kafkaMessageSizeBytes *prometheus.GaugeVec

	// Blockconsumer metrics
	topicPartitionOffset *prometheus.GaugeVec
	topicLag             *prometheus.GaugeVec
}

// NewPrometheus return an instance of Prometheus with registered metrics.
func NewPrometheus(namespace, subsystem string) *Prometheus {
	prometheus.DefaultRegisterer.Unregister(collectors.NewGoCollector())
	prometheus.DefaultRegisterer.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	return &Prometheus{namespace: namespace, subsystem: subsystem}
}
