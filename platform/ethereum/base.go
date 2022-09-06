package ethereum

import (
	"github.com/unanoc/blockchain-indexer/pkg/client"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/blockchain/types"
	"github.com/unanoc/blockchain-indexer/pkg/sentry"
)

type Platform struct {
	chain  types.ChainType
	client Client
}

func InitPlatform(chain types.ChainType, url string) *Platform {
	return &Platform{
		chain:  chain,
		client: Client{client.InitClient(url, sentry.DefaultSentryErrorHandler())},
	}
}

func (p Platform) GetChain() types.ChainType {
	return p.chain
}
