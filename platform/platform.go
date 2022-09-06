package platform

import (
	"github.com/unanoc/blockchain-indexer/pkg/primitives/blockchain/types"
)

type (
	Platform interface {
		GetChain() types.ChainType
		GetCurrentBlockNumber() (int64, error)
		GetBlockByNumber(num int64) ([]byte, error)
	}

	Platforms map[types.ChainType]Platform
)

func InitPlatforms() Platforms {
	return Platforms{}
}
