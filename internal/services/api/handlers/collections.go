package handlers

import (
	"context"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/platform"

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
	collectionDB, err := s.db.FindCollectionByContract(ctx, collection.Contract)
	if collectionDB != nil {
		log.WithError(err).Error("Collection: ", collection.Contract, " existed!")
		return nil, httperr.ErrCollectionExisted
	}

	coll, err := s.db.InsertCollection(ctx, collection)
	if err != nil {
		log.WithError(err).Error("Insert collection error")
		return nil, httperr.ErrInternalServer
	}

	return coll, nil
}

func (s *CollectionService) GetDetectSmartcontractStandard(contract string, chain string) models.ContractStandard {
	platforms := platform.InitPlatforms()
	platform := platforms[chain]

	standard, err := platform.DetectSmartcontractStandard(contract)
	if err != nil {
		log.WithError(err).Error("Detect smart contract standard err")
		return models.ContractStandard("UNKNOWN")
	}

	return models.ContractStandard(standard)
}
