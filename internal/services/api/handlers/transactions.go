package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/services/api/dtos"
	"math"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/internal/services/api/httperr"
)

var ErrTxDoesNotExist = errors.New("transaction does not exist")

type TransactionService struct {
	db repository.Storage
}

func NewTransactionService(dbConnector *repository.Storage) *TransactionService {
	return &TransactionService{*dbConnector}
}

func (s *TransactionService) GetTransactions(ctx context.Context, chain string,
	page, limit int, recent bool,
) (*dtos.TxsResp, *httperr.Error) {
	txs, err := s.db.GetTransactions(ctx, chain, page, limit, recent)
	if err != nil {
		log.WithError(err).Error("Getting txs error")

		return nil, httperr.ErrInternalServer
	}

	transactions, err := ToTxs(txs)
	if err != nil {
		log.WithError(err).Error("Txs normalizing error")

		return nil, httperr.ErrInternalServer
	}

	totalCount, err := s.db.GetTransactionTotalCount(ctx, chain)
	if err != nil {
		log.WithError(err).Error("Getting of txs count error")

		return nil, httperr.ErrInternalServer
	}

	return &dtos.TxsResp{
		StatusCode: http.StatusOK,
		TotalCount: totalCount,
		TotalPages: int(math.Ceil(float64(totalCount) / float64(limit))),
		PageNumber: page,
		Limit:      limit,
		Txs:        transactions,
	}, nil
}

func (s *TransactionService) GetTransactionByHash(ctx context.Context, chain, hash string) (*dtos.TxResp, *httperr.Error) {
	tx, err := s.db.GetTransactionByHash(ctx, chain, hash)
	if err != nil {
		if postgres.IsErrNotFound(err) {
			return nil, httperr.NewError(http.StatusNotFound, ErrTxDoesNotExist.Error())
		}

		log.WithError(err).Error("Getting tx error")

		return nil, httperr.ErrInternalServer
	}

	var metadata interface{}
	if err := json.Unmarshal(tx.Metadata.RawMessage, &metadata); err != nil {
		log.WithError(err).Error("Tx normalization error")

		return nil, httperr.ErrInternalServer
	}

	return &dtos.TxResp{
		StatusCode: http.StatusOK,
		Tx: dtos.Tx{
			Hash:      tx.Hash,
			Chain:     tx.Chain,
			Height:    tx.Block,
			From:      tx.Sender,
			To:        tx.Recipient,
			Status:    string(tx.Status),
			Type:      string(tx.Type),
			Sequence:  tx.Sequence,
			Fee:       tx.Fee,
			Data:      metadata,
			Timestamp: tx.BlockCreatedAt,
		},
	}, nil
}

func (s *TransactionService) GetTransactionsByUser(ctx context.Context, chain string,
	address string, page, limit int, recent bool,
) (*dtos.TxsResp, *httperr.Error) {
	txs, err := s.db.GetTransactionsByAddress(ctx, chain, address, page, limit, recent)
	if err != nil {
		log.WithError(err).Error("Getting txs by address error")

		return nil, httperr.ErrInternalServer
	}

	transactions, err := ToTxs(txs)
	if err != nil {
		log.WithError(err).Error("Txs normalizing error")

		return nil, httperr.ErrInternalServer
	}

	totalCount, err := s.db.GetTransactionByAddressTotalCount(ctx, chain, address)
	if err != nil {
		log.WithError(err).Error("Getting of user txs count error")

		return nil, httperr.ErrInternalServer
	}

	return &dtos.TxsResp{
		StatusCode: http.StatusOK,
		TotalCount: totalCount,
		TotalPages: int(math.Ceil(float64(totalCount) / float64(limit))),
		PageNumber: page,
		Limit:      limit,
		Txs:        transactions,
	}, nil
}

func ToTxs(txs []models.Transaction) ([]dtos.Tx, error) {
	transactions := make([]dtos.Tx, 0, len(txs))

	for _, tx := range txs {
		var metadata interface{}
		if err := json.Unmarshal(tx.Metadata.RawMessage, &metadata); err != nil {
			return nil, fmt.Errorf("failed to normalize db.Transaction to Tx: %w", err)
		}

		transactions = append(transactions, dtos.Tx{
			Hash:      tx.Hash,
			Chain:     tx.Chain,
			Height:    tx.Block,
			From:      tx.Sender,
			To:        tx.Recipient,
			Status:    string(tx.Status),
			Type:      string(tx.Type),
			Sequence:  tx.Sequence,
			Fee:       tx.Fee,
			Data:      metadata,
			Timestamp: tx.BlockCreatedAt,
		})
	}

	return transactions, nil
}
