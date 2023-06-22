package handlers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/services/api/dtos"
	"github.com/unanoc/blockchain-indexer/internal/services/api/httperr"
	"math"
	"net/http"
)

type CollectionService struct {
	db repository.Storage
}

func NewCollectionService(dbConnector *repository.Storage) *CollectionService {
	return &CollectionService{*dbConnector}
}

func (s *TransactionService) GetCollections(ctx context.Context, name string,
	page, limit int, recent bool,
) (*dtos.TxsResp, *httperr.Error) {
	collections, err := s.db.GetCollections(ctx, name, page, limit, recent)
	if err != nil {
		log.WithError(err).Error("Getting collections error")

		return nil, httperr.ErrInternalServer
	}

	transactions, err := ToTxs(txs)
	if err != nil {
		log.WithError(err).Error("Txs normalizing error")

		return nil, httperr.ErrInternalServer
	}

	totalCount, err := s.db.GetCollectionTotalCount(ctx, name)
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
