package blockproducer

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
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
	metricsPusher worker.Worker
	workers       []worker.Worker
}

func NewApp() *App {
	services.InitConfig()
	services.InitLogging()
	services.InitSentry()
	services.InitDatabase()

	db, err := postgres.New(config.Default.Database.URL, config.Default.Database.Log)
	if err != nil {
		log.WithError(err).Fatal("Database init error")
	}

	platforms := platform.InitPlatforms()

	if err = initBlockTrackers(context.Background(), db, platforms); err != nil {
		log.WithError(err).Fatal("Block trackers init error")
	}

	creatKafkaTopics(platforms)

	kafka := kafka.NewWriter(kafka.WriterConfig{
		Brokers: strings.Split(config.Default.Kafka.Brokers, ","),
		// Brokers:      []string{"localhost:9092"},
		MaxAttempts:  config.Default.Kafka.MaxAttempts,
		BatchBytes:   config.Default.Kafka.MessageMaxBytes,
		RequiredAcks: -1,
		Logger:       log.New(),
	})
	// kafka.AllowAutoTopicCreation = true

	prometheus := prometheus.NewPrometheus(config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem)
	prometheus.RegisterBlocksProducerMetrics()

	metricsPusher, err := metrics.InitDefaultMetricsPusher(
		config.Default.Prometheus.PushGateway.URL,
		config.Default.Prometheus.PushGateway.Key,
		fmt.Sprintf("%s_%s", config.Default.Prometheus.NameSpace, config.Default.Prometheus.SubSystem),
		config.Default.Prometheus.PushGateway.PushInterval,
	)
	if err != nil {
		log.WithError(err).Warn("Metrics pusher init error")
	}

	workers := make([]worker.Worker, 0, len(platforms))
	for _, pl := range platforms {
		workers = append(workers, NewWorker(db, kafka, prometheus, pl))
	}

	return &App{
		metricsPusher: metricsPusher,
		workers:       workers,
	}
}

func (a *App) Run(ctx context.Context) {
	service.RunWithGracefulShutdown(ctx, func(ctx context.Context, wg *sync.WaitGroup) {
		if a.metricsPusher != nil {
			a.metricsPusher.Start(ctx, wg)
		}

		for _, worker := range a.workers {
			go worker.Start(ctx, wg)
		}
	})
}

func initBlockTrackers(ctx context.Context, db repository.Storage, platforms platform.Platforms) error {
	for _, pl := range platforms {
		tracker, err := db.GetBlockTracker(ctx, pl.Coin().Handle)
		if err != nil {
			if !postgres.IsErrNotFound(err) {
				return fmt.Errorf("failed to get block tracker: %w", err)
			}
		}

		if tracker != nil {
			continue
		}

		var currentBlock int64

		if config.Default.BlockProducer.StartFromLastBlock {
			currentBlock, err = pl.GetCurrentBlockNumber()
			if err != nil {
				return fmt.Errorf("failed to get current block number: %w", err)
			}
		}

		if err = db.UpsertBlockTracker(ctx, pl.Coin().Handle, currentBlock); err != nil {
			return fmt.Errorf("failed to insert block tracker: %w", err)
		}
	}

	return nil
}

func creatKafkaTopics(platforms platform.Platforms) {
	conn, err := kafka.Dial("tcp", strings.Split(config.Default.Kafka.Brokers, ",")[0])
	if err != nil {
		log.WithError(err).Fatal("Kafka dial error")
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.WithError(err).Fatal("Kafka Controller init error")
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	//controllerConn, err = kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		log.WithError(err).Fatal("Kafka dial error")
	}
	defer controllerConn.Close()

	topicConfigs := make([]kafka.TopicConfig, 0, len(platforms))
	for _, pl := range platforms {
		topic := fmt.Sprintf("%s%s", config.Default.Kafka.BlocksTopicPrefix, pl.Coin().Handle)

		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     config.Default.Kafka.Partitions,
			ReplicationFactor: config.Default.Kafka.ReplicationFactor,
		})
	}

	if err = controllerConn.CreateTopics(topicConfigs...); err != nil {
		log.WithError(err).Fatal("Topic creation error")
	}
}
