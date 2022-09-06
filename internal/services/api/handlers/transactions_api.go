package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/services/api/controllers/transactions"
)

type TransactionsAPI struct {
	controller *transactions.Controller
}

func NewTransactionsAPI(db repository.Storage) API {
	return &TransactionsAPI{
		controller: transactions.NewController(db),
	}
}

// GetTransactions godoc
// @Description  Returns transactions
// @Tags         Transactions
// @Produce      json
// @Success      200  {object}  transactions.GetTransactionsResp
// @Failure      500  {object}  httperr.Error
// @Router       /api/v1/transactions [get]
func (api *TransactionsAPI) GetTransactions(c *gin.Context) {
	response, err := api.controller.GetTransactions()
	if err != nil {
		c.JSON(err.GetStatusCode(), err)

		return
	}

	c.JSON(response.StatusCode, response)
}
