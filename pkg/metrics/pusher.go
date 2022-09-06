package metrics

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/unanoc/blockchain-indexer/pkg/client"
	"github.com/unanoc/blockchain-indexer/pkg/worker"
)

type PusherClient struct {
	client client.Request
}

func NewPusherClient(pushURL, key string, errorHandler client.HTTPErrorHandler) *PusherClient {
	client := client.InitClient(pushURL, errorHandler)
	client.AddHeader("X-API-Key", key)

	return &PusherClient{
		client: client,
	}
}

func (c *PusherClient) Do(req *http.Request) (*http.Response, error) {
	for key, value := range c.client.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do http request: %w", err)
	}

	return resp, nil
}

type Pusher interface {
	Push(ctx context.Context) error
	Close() error
}

type pusher struct {
	pusher *push.Pusher
}

func NewPusher(pushgatewayURL, jobName string) Pusher {
	return &pusher{
		pusher: push.New(pushgatewayURL, jobName).
			Grouping("instance", instanceID()).
			Gatherer(prometheus.DefaultGatherer),
	}
}

func NewPusherWithCustomClient(pushgatewayURL, jobName string, client client.HTTPClient) Pusher {
	return &pusher{
		pusher: push.New(pushgatewayURL, jobName).
			Grouping("instance", instanceID()).
			Gatherer(prometheus.DefaultGatherer).
			Client(client),
	}
}

func (p *pusher) Push(ctx context.Context) error {
	if err := p.pusher.Push(); err != nil {
		return fmt.Errorf("failed to push metrics: %w", err)
	}

	return nil
}

func (p *pusher) Close() error {
	if err := p.pusher.Delete(); err != nil {
		return fmt.Errorf("failed to send DELETE to pusher: %w", err)
	}

	return nil
}

func instanceID() string {
	instance := os.Getenv("DYNO")
	if instance == "" {
		instance = os.Getenv("INSTANCE_ID")
	}
	if instance == "" {
		instance = "local"
	}

	return instance
}

func InitDefaultMetricsPusher(pushgatewayURL, pushgatewayKey, serviceName string, pushInterval time.Duration) (
	worker.Worker, error,
) {
	client := NewPusherClient(pushgatewayURL, pushgatewayKey, nil)
	pusher := NewPusherWithCustomClient(pushgatewayURL, serviceName, client)

	if err := pusher.Push(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to PushGateway, metrics won't be pushed: %w", err)
	}

	metricsPusher := NewMetricsPusherWorker(worker.DefaultOptions(pushInterval), pusher)

	return metricsPusher, nil
}

func NewMetricsPusherWorker(options *worker.Options, pusher Pusher) worker.Worker {
	return worker.NewWorkerBuilder("metrics_pusher", pusher.Push).
		WithOptions(options).
		WithStop(pusher.Close).
		Build()
}
