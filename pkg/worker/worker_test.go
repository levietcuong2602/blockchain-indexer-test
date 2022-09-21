package worker_test

import (
	"context"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"gotest.tools/assert"

	"github.com/unanoc/blockchain-indexer/pkg/worker"
)

func TestWorkerWithDefaultOptions(t *testing.T) {
	counter := 0
	logger := log.WithFields(log.Fields{"worker": "test_worker"})
	worker := worker.NewWorkerBuilder("test", logger, func(context.Context) error {
		counter++
		return nil
	}).WithOptions(worker.DefaultOptions(100 * time.Millisecond)).Build()

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	worker.Start(ctx, wg)

	time.Sleep(350 * time.Millisecond)
	cancel()
	wg.Wait()

	assert.Equal(t, 4, counter, "Should execute 4 times - 1st immidietly, and 3 after")
}

func TestWorkerStartsConsequently(t *testing.T) {
	counter := 0
	logger := log.WithFields(log.Fields{"worker": "test_worker"})
	options := worker.DefaultOptions(100 * time.Millisecond)
	options.RunConsequently = true

	worker := worker.NewWorkerBuilder("test", logger, func(context.Context) error {
		time.Sleep(100 * time.Millisecond)
		counter++
		return nil
	}).WithOptions(options).Build()

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	worker.Start(ctx, wg)

	time.Sleep(350 * time.Millisecond)
	cancel()
	wg.Wait()

	assert.Equal(t, 3, counter, "Should execute 3 times - 1st immidietly, and 2 after with delay between runs")
}
