package main

import (
	"context"

	_ "github.com/unanoc/blockchain-indexer/docs"
	"github.com/unanoc/blockchain-indexer/internal/services/api"
)

// @title   Blockchain Indexer API
// @version 1.0

// @contact.name  Cuong Lee
// @contact.email Cuong@smartosc.com

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	api.NewApp().Run(context.Background())
}
