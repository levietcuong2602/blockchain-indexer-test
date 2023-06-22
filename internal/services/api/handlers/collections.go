package handlers

import (
	"context"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"

	log "github.com/sirupsen/logrus"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/services/api/dtos"
	"github.com/unanoc/blockchain-indexer/internal/services/api/httperr"
)

type CollectionService struct {
	db repository.Storage
}

func NewCollectionService(dbConnector *repository.Storage) *CollectionService {
	return &CollectionService{*dbConnector}
}

func (s *CollectionService) GetCollections(ctx context.Context, name string,
	page, limit int, recent bool,
) (interface{}, *httperr.Error) {
	collections, err := s.db.GetCollections(ctx, name, page, limit, recent)
	if err != nil {
		log.WithError(err).Error("Getting collections error")

		return nil, httperr.ErrInternalServer
	}

	totalCount, err := s.db.GetCollectionTotalCount(ctx, name)
	if err != nil {
		log.WithError(err).Error("Getting of txs count error")

		return nil, httperr.ErrInternalServer
	}

	return dtos.CreatedCollectionPagedResponse(collections, page, limit, int(totalCount)), nil
}

func (s *CollectionService) CreateCollection(ctx context.Context, collection models.Collection) (*models.Collection, *httperr.Error) {
	collection, err := s.db.InsertCollection(ctx, collection)
	if err != nil {
		log.WithError(err).Error("Getting of txs count error")

		return nil, httperr.ErrInternalServer
	}

	return &collection, nil
}
