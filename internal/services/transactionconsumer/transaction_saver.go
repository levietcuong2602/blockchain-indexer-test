package transactionconsumer

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/pkg/mq"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

const serviceName = "transaction_saver"

type TransactionSaver struct {
	db         *postgres.Database
	prometheus *prometheus.Prometheus
}

func NewTransactionSaver(db *postgres.Database, p *prometheus.Prometheus) mq.MessageProcessor {
	return &TransactionSaver{
		db:         db,
		prometheus: p,
	}
}

func (ts *TransactionSaver) Process(message mq.Message) error {
	var block types.Block

	if err := json.Unmarshal(message, &block); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	log.WithFields(log.Fields{
		"service": serviceName,
		"txs":     len(block.Txs),
	}).Info("Consumed")

	normalizedBlock, err := models.NormalizeBlock(block)
	if err != nil {
		return fmt.Errorf("failed to normalized block: %w", err)
	}
	if err = ts.db.InsertBlock(context.Background(), *normalizedBlock); err != nil {
		return fmt.Errorf("failed to insert block: %w", err)
	}
	if len(block.Txs) == 0 {
		return nil
	}

	chain := block.Txs[0].Chain
	normalizedTxs, err := models.NormalizeTransactions(block.Txs, chain)
	if err != nil {
		return fmt.Errorf("failed to normalized txs: %w", err)
	}

	if err = ts.db.InsertTransactions(context.Background(), normalizedTxs); err != nil {
		return fmt.Errorf("failed to insert txs: %w", err)
	}

	for _, tx := range block.Txs {
		normalizeEvents, err := models.NormalizeEvents(tx.Events)
		if err != nil {
			return fmt.Errorf("failed to normalized event: %w", err)
		}
		if err = ts.db.InsertEvents(context.Background(), normalizeEvents); err != nil {
			return fmt.Errorf("failed to insert events: %w", err)
		}
	}

	ts.prometheus.SetParsedTxs(chain, len(normalizedTxs))

	return nil
}

func (ts *TransactionSaver) HandleDecodeTransaction(tx types.Tx) {
	// TODO: Optimize detect smart contract address
	contractAddress := tx.To
	// Check in collection
	conllection, err := ts.db.FindCollectionByContract(context.Background(), contractAddress)
	if err != nil {
		return
	}

	if conllection.ID == 0 {
		return
	}

	//for _, event := range tx.Events {
	//
	//}
}
