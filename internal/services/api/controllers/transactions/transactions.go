package transactions

import (
	"net/http"

	"github.com/unanoc/blockchain-indexer/internal/repository"
	"github.com/unanoc/blockchain-indexer/internal/services/api/httperr"
)

type Controller struct {
	db repository.Storage
}

func NewController(db repository.Storage) *Controller {
	return &Controller{db: db}
}

//nolint:unparam
func (i *Controller) GetTransactions() (*GetTransactionsResp, *httperr.Error) {
	return &GetTransactionsResp{
		StatusCode: http.StatusOK,
	}, nil
}
