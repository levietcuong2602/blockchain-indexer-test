package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/unanoc/blockchain-indexer/pkg/primitives/blockchain/types"
)

func (p *Prometheus) RegisterParserMetrics() {
	p.lastFetchedBlock = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "last_fetched_block"),
			Help: "Last fetched block",
		},
		[]string{labelChain},
	)
	p.currentNodeBlock = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "current_node_block"),
			Help: "Current node block",
		},
		[]string{labelChain},
	)
	p.messageSizeBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "message_size_bytes"),
			Help: "Message size in bytes",
		},
		[]string{labelChain},
	)
}

func (p *Prometheus) SetLastFetchedBlock(chain types.ChainType, block int64) {
	p.lastFetchedBlock.With(prometheus.Labels{
		labelChain: string(chain),
	}).Set(float64(block))
}

func (p *Prometheus) SetCurrentNodeBlock(chain types.ChainType, block int64) {
	p.currentNodeBlock.With(prometheus.Labels{
		labelChain: string(chain),
	}).Set(float64(block))
}

func (p *Prometheus) SetBlocksParserMessageSizeBytes(chain types.ChainType, size int) {
	p.messageSizeBytes.With(prometheus.Labels{
		labelChain: string(chain),
	}).Set(float64(size))
}
