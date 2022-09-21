package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (p *Prometheus) RegisterNodesMetrics() {
	p.nodeCurrentBlock = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "node_current_block"),
			Help: "Node current block",
		},
		[]string{labelChain, labelHost, labelVersion},
	)
	p.nodeLatency = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "node_latency"),
			Help: "Node latency (ms)",
		},
		[]string{labelChain, labelHost},
	)
	p.nodeStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "node_status"),
			Help: "Node status",
		},
		[]string{labelChain, labelHost},
	)
}

func (p *Prometheus) SetNodeCurrentBlock(chain, host, version string, blockNumber int64) {
	p.nodeCurrentBlock.With(prometheus.Labels{
		labelChain:   chain,
		labelHost:    host,
		labelVersion: version,
	}).Set(float64(blockNumber))
}

func (p *Prometheus) SetNodeLatency(chain, host string, latency int64) {
	p.nodeLatency.With(prometheus.Labels{
		labelChain: chain,
		labelHost:  host,
	}).Set(float64(latency))
}

func (p *Prometheus) SetNodeStatus(chain, host string, httpStatus int) {
	p.nodeStatus.With(prometheus.Labels{
		labelChain: chain,
		labelHost:  host,
	}).Set(float64(httpStatus))
}

func (p *Prometheus) ResetNodeCurrentBlock() {
	p.nodeCurrentBlock.Reset()
}

func (p *Prometheus) ResetNodeStatus() {
	p.nodeStatus.Reset()
}
