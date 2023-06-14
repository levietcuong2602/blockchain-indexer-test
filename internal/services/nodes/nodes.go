package nodes

import (
	"context"
	"fmt"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
)

//nolint:gochecknoglobals
var nodesList = map[string][]models.Node{
	"binance": {
		{
			Chain:      "binance",
			Scheme:     "https",
			Host:       "dex.binance.org",
			Enabled:    false,
			Monitoring: false,
		},
	},
	"cosmos": {
		{
			Chain:      "cosmos",
			Scheme:     "https",
			Host:       "us-atom-restapi.binancechain.io",
			Enabled:    false,
			Monitoring: false,
		},
		{
			Chain:      "cosmos",
			Scheme:     "https",
			Host:       "cosmos-mainnet-rpc.allthatnode.com:1317",
			Enabled:    false,
			Monitoring: false,
		},
	},
	"ethereum": {
		{
			Chain:      "ethereum",
			Scheme:     "https",
			Host:       "ethereum-mainnet-rpc.allthatnode.com",
			Enabled:    false,
			Monitoring: false,
		},
	},
	"smartchain": {
		{
			Chain:      "smartchain",
			Scheme:     "https",
			Host:       "bsc-dataseed1.binance.org",
			Enabled:    false,
			Monitoring: false,
		},
		{
			Chain:      "smartchain",
			Scheme:     "https",
			Host:       "bsc-dataseed2.binance.org",
			Enabled:    false,
			Monitoring: false,
		},
		{
			Chain:      "smartchain",
			Scheme:     "https",
			Host:       "bsc-dataseed3.binance.org",
			Enabled:    false,
			Monitoring: false,
		},
		{
			Chain:      "smartchain",
			Scheme:     "https",
			Host:       "bsc-dataseed4.binance.org",
			Enabled:    false,
			Monitoring: false,
		},
	},
	"near": {
		{
			Chain:      "near",
			Scheme:     "https",
			Host:       "rpc.mainnet.near.org",
			Enabled:    false,
			Monitoring: false,
		},
	},
	"mumbai": {
		{
			Chain:      "mumbai",
			Scheme:     "https",
			Host:       "polygon-mumbai.g.alchemy.com/v2/nf4tMvfGGwWRDhtGUcmBXaY0L-VJVseD",
			Enabled:    true,
			Monitoring: true,
		},
	},
}

func AddNodesListToDB(db *postgres.Database) error {
	for _, nodesList := range nodesList {
		if err := db.InsertNodes(context.Background(), nodesList); err != nil {
			return fmt.Errorf("failed to insert nodes: %w", err)
		}
	}

	return nil
}
