package platform

import (
	"github.com/unanoc/blockchain-indexer/internal/config"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
	"github.com/unanoc/blockchain-indexer/platform/ethereum"
)

type (
	Platform interface {
		Coin() coin.Coin
		GetCurrentBlockNumber() (int64, error)
		GetBlockByNumber(num int64) ([]byte, error)
		GetVersion() (string, error)
		NormalizeRawBlock(rawBlock []byte) (types.Txs, error)
		UpdateNodeConnection(url string)
	}

	Platforms map[string]Platform
)

func InitPlatforms() Platforms {
	return Platforms{
		coin.Smartchain().Handle: ethereum.InitPlatform(coin.SMARTCHAIN, config.Default.Platforms.Smartchain.Node),
	}
}

//nolint:gofumpt
func GetPlatform(chain string, url string) Platform {
	switch chain {
	case coin.Smartchain().Handle:
		return ethereum.InitPlatform(coin.SMARTCHAIN, url)
	default:
		return nil
	}
}
