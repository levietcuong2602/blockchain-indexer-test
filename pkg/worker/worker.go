package worker

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Builder interface {
	WithOptions(options *Options) Builder
	WithStop(func() error) Builder
	Build() Worker
}

type builder struct {
	worker *worker
}

func NewWorkerBuilder(name string, workerFn func(context.Context) error) Builder {
	return &builder{
		worker: &worker{
			name:     name,
			workerFn: workerFn,
			options:  DefaultOptions(1 * time.Minute),
		},
	}
}

func (b *builder) WithOptions(options *Options) Builder {
	b.worker.options = options

	return b
}

func (b *builder) WithStop(stopFn func() error) Builder {
	b.worker.stopFn = stopFn

	return b
}

func (b *builder) Build() Worker {
	return b.worker
}

// Worker interface can be constructed using worker.NewBuilder("worker_name", workerFn).Build()
// or allows custom implementation (e.g. one-off jobs)
type Worker interface {
	Name() string
	Start(ctx context.Context, wg *sync.WaitGroup)
}

type worker struct {
	name     string
	workerFn func(context.Context) error
	stopFn   func() error
	options  *Options
}

func (w *worker) Name() string {
	return w.name
}

//nolint:gocognit
func (w *worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(w.options.Interval)
		defer ticker.Stop()

		if w.options.RunImmediately {
			log.WithField("worker", w.name).Info("Run immediately")

			if err := w.workerFn(ctx); err != nil {
				log.WithError(err).WithField("worker", w.name).Error("Error occurred while running the worker")
			}
		}

		for {
			select {
			case <-ctx.Done():
				if w.stopFn != nil {
					log.WithField("worker", w.name).Info("Stopping...")
					if err := w.stopFn(); err != nil {
						log.WithField("worker", w.name).WithError(err).Warn("Error occurred while stopping the worker")
					}
				}
				log.WithField("worker", w.name).Info("Stopped")

				return
			case <-ticker.C:
				if w.options.RunConsequently {
					ticker.Stop()
				}

				log.WithField("worker", w.name).Info("Processing")

				if err := w.workerFn(ctx); err != nil {
					log.WithError(err).WithField("worker", w.name).Error("Error occurred while running the worker")
				}

				if w.options.RunConsequently {
					ticker = time.NewTicker(w.options.Interval)
				}
			}
		}
	}()
}
