package mumbai

import (
	"encoding/json"
	"fmt"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unanoc/blockchain-indexer/pkg/mock"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
)

func TestNormalizeChunk(t *testing.T) {
	var txs []Transaction
	err := mock.JSONModelFromFilePath("mocks/tx_input.json", &txs)
	assert.NoError(t, err)

	var receipts TransactionReceipts
	err = mock.JSONModelFromFilePath("mocks/tx_receipts_input.json", &receipts)
	assert.NoError(t, err)

	block := Block{
		Timestamp:     (*types.HexNumber)(new(big.Int).SetInt64(1621950566)),
		Transactions:  txs,
		BaseFeePerGas: (*types.HexNumber)(new(big.Int).SetInt64(123)),
	}

	p := Init(coin.MUMBAI, "")

	result := p.NormalizeBlock(block, receipts.Map())

	resultJSON, err := json.Marshal(result)
	assert.NoError(t, err)

	fmt.Println(string(resultJSON))

	expectedJSON, err := mock.JSONStringFromFilePath("mocks/tx_expected_output.json")
	assert.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(resultJSON))
}
