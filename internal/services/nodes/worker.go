package nodes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/pkg/client"
	pkghttp "github.com/unanoc/blockchain-indexer/pkg/http"
	"github.com/unanoc/blockchain-indexer/pkg/worker"
	"github.com/unanoc/blockchain-indexer/platform"
)

const workerName = "nodes_checker"

type Worker struct {
	log        *log.Entry
	db         repository.Storage
	prometheus *prometheus.Prometheus
}

func NewWorker(db repository.Storage, p *prometheus.Prometheus) worker.Worker {
	w := &Worker{
		log:        log.WithField("worker", workerName),
		db:         db,
		prometheus: p,
	}

	opts := &worker.Options{
		Interval:        config.Default.Nodes.Interval,
		RunImmediately:  true,
		RunConsequently: false,
	}

	return worker.NewWorkerBuilder(workerName, w.run).WithOptions(opts).Build()
}

func (w *Worker) run(ctx context.Context) error {
	nodes, err := w.db.GetNodes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	platforms := getPlatforms(nodes)

	w.prometheus.ResetNodeCurrentBlock()
	w.prometheus.ResetNodeStatus()

	for host, platform := range platforms {
		go w.checkNode(host, platform)
	}

	return nil
}

func (w *Worker) checkNode(host string, nodeAPI platform.Platform) {
	var blockNumber int64
	var latency time.Duration
	var version string
	var err error

	chain := nodeAPI.Coin().Handle

	version, err = nodeAPI.GetVersion()
	if err != nil {
		log.WithFields(log.Fields{"host": host, "error": err}).Error("Getting node version error")
	}

	now := time.Now()
	blockNumber, err = nodeAPI.GetCurrentBlockNumber()
	if err != nil {
		log.WithFields(log.Fields{"host": host, "error": err}).Error("Getting current block number error")
		w.prometheus.SetNodeStatus(chain, host, getHTTPErrorCode(err))

		return
	}

	latency = time.Since(now)

	w.prometheus.SetNodeCurrentBlock(chain, host, version, blockNumber)
	w.prometheus.SetNodeLatency(chain, host, latency.Milliseconds())
	w.prometheus.SetNodeStatus(chain, host, http.StatusOK)

	log.WithFields(log.Fields{
		"coin":         chain,
		"host":         host,
		"version":      version,
		"block_number": blockNumber,
	}).Info("success")
}

func getHTTPErrorCode(err error) int {
	httpErr := &client.HTTPError{}
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode
	}

	return http.StatusInternalServerError
}

func getPlatforms(nodes []models.Node) map[string]platform.Platform {
	platforms := make(map[string]platform.Platform)

	for _, node := range nodes {
		if node.Monitoring {
			platform := platform.GetPlatform(node.Chain, pkghttp.BuildURL(node.Scheme, node.Host))
			if platform != nil {
				platforms[node.Host] = platform
			}
		}
	}

	return platforms
}
