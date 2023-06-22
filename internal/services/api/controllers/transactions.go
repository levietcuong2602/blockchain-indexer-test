package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/services/api/handlers"
	"github.com/unanoc/blockchain-indexer/internal/services/api/validations"
)

type ITransactionController interface {
	GetTransactions(*gin.Context)
	GetTransactionByHash(*gin.Context)
	GetTransactionsByUser(*gin.Context)
}

type TransactionController struct {
	service *handlers.TransactionService
}

func NewTransactionController(db repository.Storage) *TransactionController {
	return &TransactionController{service: handlers.NewTransactionService(&db)}
}

// GetTransactions godoc
// @Description  Returns all transaction list by creation date order(asc/desc)
// @Tags         Transactions
// @Produce      json
// @Param        chain query string true "Chain"
// @Param        page query int false "Page for pagination"
// @Param        limit query int false "The limit of the number of items"
// @Param        recent query bool false "Enable desc order"
// @Success      200  {object}  transactions.TxsResp
// @Failure      400  {object}  httperr.Error
// @Failure      500  {object}  httperr.Error
// @Router       /api/v1/transactions [get]
func (api *TransactionController) GetTransactions(c *gin.Context) {
	params, err := validations.ValidateTransactionsParams(c)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	response, err := api.service.GetTransactions(c.Request.Context(),
		params.Chain, params.Page, params.Limit, params.Recent)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	c.JSON(response.StatusCode, response)
}

// GetTransactionByHash godoc
// @Description  Returns transaction details by hash
// @Tags         Transactions
// @Produce      json
// @Param        chain query string true "Chain"
// @Param        hash  path string true "Transaction hash"
// @Success      200  {object}  transactions.TxResp
// @Failure      400  {object}  httperr.Error
// @Failure      404  {object}  httperr.Error
// @Failure      500  {object}  httperr.Error
// @Router       /api/v1/transactions/{hash} [get]
func (api *TransactionController) GetTransactionByHash(c *gin.Context) {
	params, err := validations.ValidateTransactionParams(c)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	response, err := api.service.GetTransactionByHash(c.Request.Context(), params.Chain, params.Hash)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	c.JSON(response.StatusCode, response)
}

// GetUserTransactionsByUser godoc
// @Description  Returns all user transaction list by creation date order(asc/desc)
// @Tags         Transactions
// @Produce      json
// @Param        chain query string true "Chain"
// @Param        address query string true "Address that made transactions"
// @Param        page query int false "Page for pagination"
// @Param        limit query int false "The limit of the number of items"
// @Param        recent query bool false "Enable desc order"
// @Success      200  {object}  transactions.TxsResp
// @Failure      400  {object}  httperr.Error
// @Failure      500  {object}  httperr.Error
// @Router       /api/v1/transactions/user [get]
func (api *TransactionController) GetTransactionsByUser(c *gin.Context) {
	params, err := validations.ValidateUserTransactionsParams(c)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	response, err := api.service.GetTransactionsByUser(c.Request.Context(),
		params.Chain, params.Address, params.Page, params.Limit, params.Recent)
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	c.JSON(response.StatusCode, response)
}
