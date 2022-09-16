package blockconsumer

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/internal/services"
	"github.com/unanoc/blockchain-indexer/pkg/metrics"
	"github.com/unanoc/blockchain-indexer/pkg/service"
	"github.com/unanoc/blockchain-indexer/pkg/worker"
	"github.com/unanoc/blockchain-indexer/platform"
)

type App struct {
	metricsPusher worker.Worker
	workers       []worker.Worker
}

func NewApp() *App {
	services.Setup()

	db, err := postgres.New(config.Default.Database.URL, config.Default.Database.Log)
	if err != nil {
		log.WithError(err).Fatal("Database init error")
	}

	platforms := platform.InitPlatforms()

	prometheus := prometheus.NewPrometheus(config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem)
	prometheus.RegisterBlocksConsumerMetrics()

	metricsPusher, err2 := metrics.InitDefaultMetricsPusher(
		config.Default.Prometheus.PushGateway.URL,
		config.Default.Prometheus.PushGateway.Key,
		fmt.Sprintf("%s_%s", config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem),
		config.Default.Prometheus.PushGateway.PushInterval,
	)
	if err2 != nil {
		log.WithError(err2).Warn("Metrics pusher init error")
	}

	workers := make([]worker.Worker, 0, len(platforms))
	for _, pl := range platforms {
		kafka := kafka.NewReader(kafka.ReaderConfig{
			Brokers:       strings.Split(config.Default.Kafka.Brokers, ","),
			MaxAttempts:   config.Default.Kafka.MaxAttempts,
			Topic:         fmt.Sprintf("%s%s", config.Default.Kafka.BlocksTopicPrefix, pl.Coin().Handle),
			GroupID:       pl.Coin().Handle,
			StartOffset:   kafka.FirstOffset,
			RetentionTime: config.Default.Kafka.RetentionTime,
		})

		workers = append(workers, NewWorker(db, kafka, prometheus, pl))
	}

	return &App{
		metricsPusher: metricsPusher,
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
