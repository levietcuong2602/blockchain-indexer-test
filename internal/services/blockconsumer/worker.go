package blockconsumer

import (
	"context"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/pkg/worker"
	"github.com/unanoc/blockchain-indexer/platform"
)

const (
	workerName = "block_consumer"
)

type Worker struct {
	log        *log.Entry
	db         *postgres.Database
	kafka      *kafka.Reader
	prometheus *prometheus.Prometheus
	API        platform.Platform
}

func NewWorker(db *postgres.Database, kafka *kafka.Reader,
	p *prometheus.Prometheus, pl platform.Platform,
) worker.Worker {
	w := &Worker{
		log: log.WithFields(log.Fields{
			"worker": workerName,
			"chain":  pl.GetChain(),
		}),
		db:         db,
		kafka:      kafka,
		prometheus: p,
		API:        pl,
	}

	opts := &worker.Options{
		Interval:        config.Default.BlockConsumer.Interval,
		RunImmediately:  true,
		RunConsequently: false,
	}

	return worker.NewWorkerBuilder(workerName, w.run).WithOptions(opts).Build()
}

func (w *Worker) run(ctx context.Context) error {
	return nil
}
