package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/services/api/controllers"
)

type TransactionsRouter struct {
	controller *controllers.TransactionController
}

func NewTransactionsAPI(db repository.Storage) API {
	return &TransactionsRouter{
		controller: controllers.NewTransactionController(db),
	}
}

func (api *TransactionsRouter) Setup(router *gin.RouterGroup) {
	router.GET("/", api.controller.GetTransactions)
	router.GET("/:hash", api.controller.GetTransactionByHash)
	router.GET("/user", api.controller.GetTransactionsByUser)
}
