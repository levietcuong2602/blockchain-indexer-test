package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/unanoc/blockchain-indexer/pkg/primitives/blockchain/types"
)

func (p *Prometheus) RegisterBlocksProducerMetrics() {
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
	p.kafkaMessageSizeBytes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "kafka_message_size_bytes"),
			Help: "Kafka message size in bytes",
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

func (p *Prometheus) SetKafkaMessageSizeBytes(chain types.ChainType, size int) {
	p.kafkaMessageSizeBytes.With(prometheus.Labels{
		labelChain: string(chain),
	}).Set(float64(size))
}
