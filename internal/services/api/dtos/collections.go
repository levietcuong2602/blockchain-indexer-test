package dtos

import (
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
)

type CreateCollectionBodyDtos struct {
	Name            string `form:"name" json:"name" xml:"name" binding:"required"`
	Slug            string `form:"slug" json:"slug" xml:"slug" binding:"required"`
	Metadata        string `form:"metadata" json:"metadata" xml:"metadata" binding:"required"`
	Contract        string `form:"contract" json:"contract" xml:"contract" binding:"required"`
	TokenCount      int64  `form:"token_count" json:"token_count" xml:"token_count" binding:"required"`
	MintedTimestamp int64  `form:"minted_timestamp" json:"minted_timestamp" xml:"minted_timestamp"`
}

type GetCollectionQueryDtos struct {
	Name   string `form:"name" json:"name" xml:"name" binding:"required"`
	Page   int    `form:"name" json:"name" xml:"name" binding:"required"`
	Limit  int    `form:"limit" json:"limit" xml:"limit"`
	Recent bool   `form:"recent" json:"recent" xml:"recent" `
}

func CreatedCollectionPagedResponse(collections []models.Collection, page, page_size, count int) interface{} {
	var resources = make([]interface{}, len(collections))
	for index, product := range collections {
		resources[index] = product
	}
	return CreatePagedResponse(resources, "collections", page, page_size, count)
}
