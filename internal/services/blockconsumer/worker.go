package blockconsumer

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
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
			"chain":  pl.Coin().Handle,
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
	message, err := w.kafka.FetchMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch kafka message: %w", err)
	}

	txs, err := w.API.NormalizeRawBlock(message.Value)
	if err != nil {
		return fmt.Errorf("failed to normalize raw block: %w", err)
	}

	// txs.CleanMemos()
	normalizedTxs, err := models.NormalizeTransactions(txs, w.API.Coin().Handle)
	if err != nil {
		return fmt.Errorf("failed to normalized txs: %w", err)
	}

	if err = w.db.InsertTransactions(ctx, normalizedTxs); err != nil {
		return fmt.Errorf("failed to insert txs: %w", err)
	}

	if err = w.kafka.CommitMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to commit kafka message (topic=%s offset=%d partition=%d): %w",
			message.Topic, message.Offset, message.Partition, err)
	}

	log.WithFields(log.Fields{
		"chain":     w.API.Coin().Handle,
		"txs":       len(txs),
		"partition": message.Partition,
		"offset":    message.Offset,
	}).Info("Transactions have been consumed")

	return nil
}
