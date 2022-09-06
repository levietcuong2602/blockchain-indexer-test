package main

import (
	"context"

	"github.com/unanoc/blockchain-indexer/internal/services/blockproducer"
)

func main() {
	blockproducer.NewApp().Run(context.Background())
}
