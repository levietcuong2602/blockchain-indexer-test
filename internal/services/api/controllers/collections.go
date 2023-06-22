package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/services/api/handlers"
	"github.com/unanoc/blockchain-indexer/internal/services/api/validations"
)

type ICollectionController interface {
	CreateCollection(*gin.Context)
	UpdateCollection(*gin.Context)
	GetCollections(*gin.Context)
}

type CollectionController struct {
	service *handlers.CollectionService
}

func NewCollectionController(db repository.Storage) *CollectionController {
	return &CollectionController{service: handlers.NewCollectionService(&db)}
}

func (api *CollectionController) CreateCollection(c *gin.Context) {

}
func (api *CollectionController) UpdateCollection(c *gin.Context) {

}

// GetCollections godoc
// @Description  Returns all transaction list by creation date order(asc/desc)
// @Tags         Collection
// @Produce      json
// @Param        page query int false "Page for pagination"
// @Param        limit query int false "The limit of the number of items"
// @Param        recent query bool false "Enable desc order"
// @Success      200  {object}  collection.TxsResp
// @Failure      400  {object}  httperr.Error
// @Failure      500  {object}  httperr.Error
// @Router       /api/v1/collections [get]
func (api *CollectionController) GetCollections(c *gin.Context) {
	params, err := validations.ValidateCollectionsParams(c)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	response, err := api.service.GetCollections(c.Request.Context(),
		params.Chain, params.Page, params.Limit, params.Recent)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	c.JSON(response.StatusCode, response)
}
