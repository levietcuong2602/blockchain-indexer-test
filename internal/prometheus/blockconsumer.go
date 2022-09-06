package prometheus

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/segmentio/kafka-go"
)

func (p *Prometheus) RegisterBlocksConsumerMetrics() {
	p.topicPartitionOffset = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "topic_partition_offset"),
			Help: "Offset of topic's partition",
		},
		[]string{labelChain, labelTopic, labelPartition},
	)
	p.topicLag = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(p.namespace, p.subsystem, "topic_lag"),
			Help: "Topic lag",
		},
		[]string{labelChain, labelTopic},
	)
}

func (p *Prometheus) SetBlocksConsumerTopicPartitionOffset(chain, topic string, partition int, offset int64) {
	p.topicPartitionOffset.With(prometheus.Labels{
		labelChain:     chain,
		labelTopic:     topic,
		labelPartition: strconv.Itoa(partition),
	}).Set(float64(offset))
}

func (p *Prometheus) SetBlocksConsumerMetrics(ctx context.Context, kafkaReader *kafka.Reader, chain string) {
	go func(ctx context.Context, kafkaReader *kafka.Reader, chain string) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				stats := kafkaReader.Stats()

				p.topicLag.With(prometheus.Labels{
					labelChain: chain,
					labelTopic: stats.Topic,
				}).Set(float64(stats.Lag))

				time.Sleep(10 * time.Second)
			}
		}
	}(ctx, kafkaReader, chain)
}
