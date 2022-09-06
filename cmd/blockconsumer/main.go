package main

import (
	"context"

	"github.com/unanoc/blockchain-indexer/internal/services/blockconsumer"
)

func main() {
	blockconsumer.NewApp().Run(context.Background())
}
