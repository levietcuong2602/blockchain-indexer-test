package blockproducer

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/pkg/http"
	"github.com/unanoc/blockchain-indexer/pkg/worker"
	"github.com/unanoc/blockchain-indexer/platform"
)

const (
	workerName = "block_producer"
)

type Worker struct {
	log          *log.Entry
	db           *postgres.Database
	kafka        *kafka.Writer
	prometheus   *prometheus.Prometheus
	API          platform.Platform
	nodes        []string
	lastNodeUsed int
}

func NewWorker(db *postgres.Database, k *kafka.Writer, p *prometheus.Prometheus, pl platform.Platform) worker.Worker {
	w := &Worker{
		log:          log.WithFields(log.Fields{"worker": workerName, "chain": pl.Coin().Handle}),
		db:           db,
		kafka:        k,
		prometheus:   p,
		API:          pl,
		nodes:        make([]string, 0),
		lastNodeUsed: 0,
	}

	opts := &worker.Options{
		Interval:        config.Default.BlockProducer.Interval,
		RunImmediately:  true,
		RunConsequently: false,
	}

	return worker.NewWorkerBuilder(workerName, w.run).WithOptions(opts).Build()
}

func (w *Worker) run(ctx context.Context) error {
	chain := w.API.Coin().Handle

	nodes, err := w.db.GetNodesByChain(ctx, chain)
	if err != nil {
		log.WithError(err).WithField("chain", chain).Error("Getting nodes list error")

		return nil
	}

	for _, node := range nodes {
		w.nodes = append(w.nodes, http.BuildURL(node.Scheme, node.Host))
	}

	if err := w.fetch(ctx); err != nil {
		time.Sleep(config.Default.BlockProducer.BackoffInterval)

		return err
	}

	return nil
}

func (w *Worker) fetch(ctx context.Context) error {
	chain := w.API.Coin().Handle

	tracker, err := w.db.GetBlockTracker(ctx, chain)
	if err != nil {
		return fmt.Errorf("failed to get block tracker: %w", err)
	}

	fromBlock, toBlock, err := w.getBlocksIntervalToFetch(tracker)
	if err != nil {
		return fmt.Errorf("failed to get blocks interval: %w", err)
	}

	if fromBlock == toBlock {
		log.WithField("current_block", toBlock).Info("No new blocks")

		return nil
	}

	blocks, err := w.fetchBlocks(fromBlock, toBlock)
	if err != nil {
		return fmt.Errorf("failed to get fetch blocks: %w", err)
	}

	if len(blocks) == 0 {
		return nil
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Num < blocks[j].Num
	})

	var lastFetchedBlock int64
	for _, block := range blocks {
		if err = w.writeBlockToKafka(ctx, block); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"block_num":  block.Num,
				"block_size": len(block.Data),
			}).Error("Sending blocks to Kafka error")

			break
		}

		lastFetchedBlock = block.Num
	}

	if err = w.db.UpsertBlockTracker(ctx, chain, lastFetchedBlock); err != nil {
		return fmt.Errorf("failed to update block tracker: %w", err)
	}

	w.prometheus.SetLastFetchedBlock(chain, lastFetchedBlock)

	return nil
}

func (w *Worker) getBlocksIntervalToFetch(tracker *models.BlockTracker) (int64, int64, error) {
	chain := w.API.Coin().Handle
	lastParsedBlock := tracker.Height

	currentBlock, err := w.API.GetCurrentBlockNumber()
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			if nextNode := w.getNextNode(); nextNode != "" {
				w.API.UpdateNodeConnection(nextNode)
			}
		}

		return 0, 0, fmt.Errorf("failed to get current block number: %w", err)
	}

	w.prometheus.SetCurrentNodeBlock(chain, currentBlock)

	fromBlock, toBlock := getNextBlocksToParse(lastParsedBlock,
		currentBlock, config.Default.BlockProducer.FetchBlocksMax)

	return fromBlock, toBlock, nil
}

func getNextBlocksToParse(lastParsedBlock, currentBlock, maxBlocks int64) (int64, int64) {
	// if there are no new blocks since last time
	if lastParsedBlock == currentBlock {
		return lastParsedBlock, currentBlock
	}

	// if current block is 0 or node has problems
	if lastParsedBlock > currentBlock {
		return lastParsedBlock, lastParsedBlock
	}

	fromBlock := lastParsedBlock + 1
	toBlock := currentBlock

	if currentBlock-lastParsedBlock > maxBlocks {
		toBlock = lastParsedBlock + maxBlocks
	}

	return fromBlock, toBlock
}

type BlockData struct {
	Num  int64
	Data []byte
}

func (w *Worker) fetchBlocks(fromBlock, toBlock int64) ([]BlockData, error) {
	chain := w.API.Coin().Handle

	blocksCount := toBlock - fromBlock + 1

	var (
		blocksChan = make(chan BlockData, blocksCount)
		errorsChan = make(chan error, blocksCount)
		totalCount int32
		wg         sync.WaitGroup
	)

	for i := fromBlock; i <= toBlock; i++ {
		wg.Add(1)

		go func(i int64, wg *sync.WaitGroup) {
			defer wg.Done()

			if err := w.fetchBlock(i, blocksChan); err != nil {
				errorsChan <- err

				return
			}

			atomic.AddInt32(&totalCount, 1)
		}(i, &wg)
	}

	wg.Wait()
	close(errorsChan)
	close(blocksChan)

	if len(errorsChan) > 0 {
		errorsList := make([]error, 0, len(errorsChan))
		for err := range errorsChan {
			errorsList = append(errorsList, err)
		}

		log.WithFields(log.Fields{
			"chain":  chain,
			"count":  len(errorsList),
			"errors": errorsList,
		}).Error("Fetch Blocks Errors")

		return nil, fmt.Errorf("failed to fetch blocks: %d: %d", fromBlock, toBlock) //nolint:goerr113
	}

	blocks := make([]BlockData, 0, len(blocksChan))
	for block := range blocksChan {
		blocks = append(blocks, block)
	}

	log.WithFields(log.Fields{
		"chain": chain,
		"from":  fromBlock,
		"to":    toBlock,
		"total": totalCount,
	}).Info("Fetched blocks batch")

	return blocks, nil
}

func (w *Worker) fetchBlock(num int64, blocksChan chan<- BlockData) error {
	block, err := w.getBlockByNumberWithRetry(config.Default.BlockProducer.BlockRetryNum,
		config.Default.BlockProducer.BlockRetryInterval, num)
	if err != nil {
		return fmt.Errorf("failed to get block by number %d: %w", num, err)
	}

	blocksChan <- BlockData{Num: num, Data: block}

	return nil
}

func (w *Worker) getBlockByNumberWithRetry(attempts int, sleep time.Duration, num int64) ([]byte, error) {
	chain := w.API.Coin().Handle

	block, err := w.API.GetBlockByNumber(num)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"chain":     chain,
			"block_num": num,
		}).Warn("Getting block by number error")

		if attempts--; attempts > 0 {
			newNode := w.getNextNode()

			log.WithFields(log.Fields{
				"chain":        chain,
				"number":       num,
				"attempts":     attempts,
				"sleep":        sleep.String(),
				"new_node_url": newNode,
			}).Warn("GetBlockByNumber retry")

			if newNode == "" {
				return w.getBlockByNumberWithRetry(attempts, sleep, num)
			}

			time.Sleep(sleep)
			w.API.UpdateNodeConnection(newNode)

			return w.getBlockByNumberWithRetry(attempts, sleep, num)
		}

		return nil, fmt.Errorf("failed to get block by number after retry: %w", err)
	}

	return block, nil
}

func (w *Worker) writeBlockToKafka(ctx context.Context, block BlockData) error {
	chain := w.API.Coin().Handle

	topic := fmt.Sprintf("%s%s", config.Default.Kafka.BlocksTopicPrefix, chain)

	if err := w.kafka.WriteMessages(ctx, kafka.Message{
		Value: block.Data,
		Topic: topic,
	}); err != nil {
		return fmt.Errorf("failed to write message to Kafka: %w", err)
	}

	log.WithFields(log.Fields{
		"chain":      chain,
		"block_num":  block.Num,
		"topic":      topic,
		"block_size": len(block.Data),
	}).Info("Produced to Kafka")

	w.prometheus.SetKafkaMessageSizeBytes(chain, len(block.Data))

	return nil
}

func (w *Worker) getNextNode() string {
	if len(w.nodes) > 0 {
		if w.lastNodeUsed+1 < len(w.nodes) {
			w.lastNodeUsed++

			return w.nodes[w.lastNodeUsed]
		}
		w.lastNodeUsed = 0

		return w.nodes[0]
	}

	return ""
}
