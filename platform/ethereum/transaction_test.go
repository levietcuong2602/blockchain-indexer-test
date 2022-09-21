package ethereum

import (
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/unanoc/blockchain-indexer/pkg/mock"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

func TestNormalizeTransactions(t *testing.T) {
	tests := []struct {
		name              string
		inputFile         string
		receiptsInputFile string
		outputFile        string
		coin              coin.Coin
		block             int64
	}{
		{
			name:              "ethereum",
			inputFile:         "mock/txs_input_ethereum.json",
			receiptsInputFile: "mock/tx_receipts_input_ethereum.json",
			outputFile:        "mock/txs_expected_output_ethereum.json",
			coin:              coin.Ethereum(),
		},
		{
			name:              "smartchain",
			inputFile:         "mock/txs_input_smartchain.json",
			receiptsInputFile: "mock/tx_receipts_input_smartchain.json",
			outputFile:        "mock/txs_expected_output_smartchain.json",
			coin:              coin.Smartchain(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var txs []Transaction
			err := mock.JSONModelFromFilePath(tc.inputFile, &txs)
			assert.NoError(t, err)

			var receipts TransactionReceipts
			if tc.receiptsInputFile != "" {
				err = mock.JSONModelFromFilePath(tc.receiptsInputFile, &receipts)
				assert.NoError(t, err)
			}

			bl := Block{
				Timestamp:     (*types.HexNumber)(new(big.Int).SetInt64(1621950566)),
				Transactions:  txs,
				BaseFeePerGas: (*types.HexNumber)(new(big.Int).SetInt64(123)),
			}

			p := Init(tc.coin.ID, "")

			result := p.NormalizeBlock(bl, receipts.Map())
			resultJSON, err := json.Marshal(result)
			assert.NoError(t, err)

			expectedJSON, err := mock.JSONStringFromFilePath(tc.outputFile)
			assert.NoError(t, err)

			assert.JSONEq(t, expectedJSON, string(resultJSON))
		})
	}
}

func TestHexDataToBigInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *big.Int
		errFn    assert.ErrorAssertionFunc
	}{
		{
			name:     "zero",
			input:    "0",
			expected: new(big.Int).SetInt64(0),
			errFn:    assert.NoError,
		},
		{
			name:     "hex_zero",
			input:    "0x0",
			expected: new(big.Int).SetInt64(0),
			errFn:    assert.NoError,
		},
		{
			name:     "one",
			input:    "1",
			expected: new(big.Int).SetInt64(1),
			errFn:    assert.NoError,
		},
		{
			name:     "hex_one",
			input:    "0x1",
			expected: new(big.Int).SetInt64(1),
			errFn:    assert.NoError,
		},
		{
			name:     "ten",
			input:    "a",
			expected: new(big.Int).SetInt64(10),
			errFn:    assert.NoError,
		},
		{
			name:     "hex_ten",
			input:    "0xa",
			expected: new(big.Int).SetInt64(10),
			errFn:    assert.NoError,
		},
		{
			name:     "zero_prefixed_hex_ten",
			input:    "0x0000000000a",
			expected: new(big.Int).SetInt64(10),
			errFn:    assert.NoError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := hexDataToBigInt(tc.input)
			tc.errFn(t, err)
			if tc.expected != nil {
				assert.True(t, tc.expected.Cmp(result) == 0)
			}
		})
	}
}

func TestNormalizeTransaction_CatchPanics(t *testing.T) {
	num1234 := types.HexNumber(*big.NewInt(1234))
	ts := types.HexNumber(*big.NewInt(time.Now().Unix()))

	testCases := []struct {
		name      string
		txs       []Transaction
		receipts  map[string]TransactionReceipt
		timestamp *types.HexNumber
	}{
		{name: "all empty"},
		{name: "tx with empty fields", txs: []Transaction{
			{},
		}},
		{timestamp: &ts, name: "tx with empty fields with nonce", txs: []Transaction{
			{
				BlockNumber: &num1234,
				Nonce:       &num1234,
			},
		}},
		{timestamp: &ts, name: "contract call with empty fields", txs: []Transaction{
			{
				BlockNumber: &num1234,
				Nonce:       &num1234,
				Gas:         &num1234,
				GasPrice:    &num1234,
			},
		}},
		{timestamp: &ts, name: "contract call with value", txs: []Transaction{
			{
				BlockNumber: &num1234,
				Nonce:       &num1234,
				Gas:         &num1234,
				GasPrice:    &num1234,
				Value:       &num1234,
			},
		}},
		{timestamp: &ts, name: "transfer with empty fields", txs: []Transaction{
			{
				BlockNumber: &num1234,
				Nonce:       &num1234,
				Gas:         &num1234,
				GasPrice:    &num1234,
				Input:       "0x",
			},
		}},
		{
			timestamp: &ts,
			name:      "token transfer with empty fields",
			receipts: map[string]TransactionReceipt{
				"": {},
			},
			txs: []Transaction{
				{
					BlockNumber: &num1234,
					Nonce:       &num1234,
					Gas:         &num1234,
					GasPrice:    &num1234,
				},
			},
		},
		{
			timestamp: &ts,
			name:      "token transfer with empty fields",
			receipts: map[string]TransactionReceipt{
				"tx-hash": {
					Logs: []EventLog{{}},
				},
			},
			txs: []Transaction{
				{
					Hash:        "tx-hash",
					BlockNumber: &num1234,
					Nonce:       &num1234,
					Gas:         &num1234,
					GasPrice:    &num1234,
				},
			},
		},
		{
			timestamp: &ts,
			name:      "token transfer with some fields filled",
			receipts: map[string]TransactionReceipt{
				"tx-hash": {
					Logs: []EventLog{{
						Topics: []string{hashEventTypeTransfer, "", ""},
					}},
				},
			},
			txs: []Transaction{
				{
					Hash:        "tx-hash",
					BlockNumber: &num1234,
					Nonce:       &num1234,
					Gas:         &num1234,
					GasPrice:    &num1234,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			block := Block{
				Timestamp:     (*types.HexNumber)(new(big.Int).SetInt64(1621950566)),
				Transactions:  tc.txs,
				BaseFeePerGas: (*types.HexNumber)(new(big.Int).SetInt64(123)),
			}

			p := Platform{}
			assert.NotPanics(t, func() { _ = p.NormalizeBlock(block, tc.receipts) })
		})
	}
}
