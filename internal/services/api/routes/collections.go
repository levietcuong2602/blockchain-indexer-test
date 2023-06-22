package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/services/api/controllers"
)

type CollectionsRouter struct {
	controller *controllers.CollectionController
}

func NewCollectionsAPI(db repository.Storage) API {
	return &CollectionsRouter{
		controller: controllers.NewCollectionController(db),
	}
}

func (api *CollectionsRouter) Setup(router *gin.RouterGroup) {
	router.POST("/", api.controller.CreateCollection)
}
