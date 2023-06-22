package dtos

import (
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"net/http"
)

func CreatedCollectionPagedResponse(request *http.Request, collections []models.Collection, page, page_size, count int) interface{} {
	var resources = make([]interface{}, len(collections))
	for index, product := range collections {
		resources[index] = product
	}
	return CreatePagedResponse(request, resources, "collections", page, page_size, count)
}
