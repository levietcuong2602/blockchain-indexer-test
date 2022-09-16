package platform

import (
	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/blockchain/types"
	"github.com/unanoc/blockchain-indexer/platform/ethereum"
)

type (
	Platform interface {
		GetChain() types.ChainType
		GetCurrentBlockNumber() (int64, error)
		GetBlockByNumber(num int64) ([]byte, error)
		// NormalizeRawBlock(rawBlock []byte) (types.Txs, error)
	}

	Platforms map[types.ChainType]Platform
)

func InitPlatforms() Platforms {
	return Platforms{
		types.BSC: ethereum.InitPlatform(types.BSC, config.Default.Platforms.Smartchain.Node),
	}
}
