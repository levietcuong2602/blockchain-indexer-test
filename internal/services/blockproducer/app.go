package blockproducer

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/internal/services"
	"github.com/unanoc/blockchain-indexer/pkg/metrics"
	"github.com/unanoc/blockchain-indexer/pkg/service"
	"github.com/unanoc/blockchain-indexer/pkg/worker"
	"github.com/unanoc/blockchain-indexer/platform"
)

type App struct {
	db            repository.Storage
	prometheus    *prometheus.Prometheus
	platforms     platform.Platforms
	metricsPusher worker.Worker
	kafkaProducer *kafka.Writer
	workers       []worker.Worker
}

func NewApp() *App {
	services.Setup()

	db, err := postgres.New(config.Default.Database.URL, config.Default.Database.Log)
	if err != nil {
		log.WithError(err).Fatal("Database init error")
	}

	platforms := platform.InitPlatforms()

	if err := initBlockTrackers(context.Background(), db, platforms); err != nil {
		log.WithError(err).Warn("Block trackers init error")
	}

	prometheus := prometheus.NewPrometheus(config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem)
	prometheus.RegisterParserMetrics()

	metricsPusher, err2 := metrics.InitDefaultMetricsPusher(
		config.Default.Prometheus.PushGateway.URL,
		config.Default.Prometheus.PushGateway.Key,
		fmt.Sprintf("%s_%s", config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem),
		config.Default.Prometheus.PushGateway.PushInterval,
	)
	if err2 != nil {
		log.WithError(err2).Warn("Metrics pusher init error")
	}

	kafkaProducer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      config.Default.Kafka.Brokers,
		MaxAttempts:  config.Default.Kafka.MaxAttempts,
		BatchBytes:   config.Default.Kafka.MessageMaxBytes,
		RequiredAcks: -1,
	})

	workers := make([]worker.Worker, 0, len(platforms))
	for _, pl := range platforms {
		workers = append(workers, NewWorker(db, kafkaProducer, prometheus, pl))
	}

	return &App{
		db:            db,
		prometheus:    prometheus,
		platforms:     platforms,
		metricsPusher: metricsPusher,
		kafkaProducer: kafkaProducer,
		workers:       workers,
	}
}

func (a *App) Run(ctx context.Context) {
	service.RunWithGracefulShutdown(ctx, func(ctx context.Context, wg *sync.WaitGroup) {
		for _, worker := range a.workers {
			go worker.Start(ctx, wg)
		}

		if a.metricsPusher != nil {
			a.metricsPusher.Start(ctx, wg)
		}
	})
}

func initBlockTrackers(ctx context.Context, db repository.Storage, platforms platform.Platforms) error {
	for _, pl := range platforms {
		tracker, err := db.GetBlockTracker(ctx, pl.GetChain())
		if err != nil {
			if !postgres.IsErrNotFound(err) {
				return fmt.Errorf("failed to get block tracker: %w", err)
			}
		}

		if tracker != nil {
			continue
		}

		if err = db.UpsertBlockTracker(ctx, pl.GetChain(), 0); err != nil {
			return fmt.Errorf("failed to insert block tracker: %w", err)
		}
	}

	return nil
}
