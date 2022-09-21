package main

import (
	"context"

	"github.com/unanoc/blockchain-indexer/internal/services/nodes"
)

func main() {
	nodes.NewApp().Run(context.Background())
}
