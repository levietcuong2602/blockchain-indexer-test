package ethereum

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/unanoc/blockchain-indexer/pkg/client"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/strings"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

const (
	MethodBlockNumber        = "eth_blockNumber"
	MethodBlockByNumber      = "eth_getBlockByNumber"
	MethodTransactionReceipt = "eth_getTransactionReceipt"
	MethodClientVersion      = "web3_clientVersion"

	rpcBatchSize = 45
)

type Client struct {
	client.Request
}

func (c *Client) GetCurrentBlockNumber() (int64, error) {
	var blockNumber types.HexNumber

	err := c.RPCCall(&blockNumber, MethodBlockNumber, nil)
	if err != nil {
		return 0, err
	}

	return (*big.Int)(&blockNumber).Int64(), nil
}

func (c *Client) GetBlockByNumber(num int64) (*Block, error) {
	var block Block
	params := []interface{}{(*types.HexNumber)(new(big.Int).SetInt64(num)), true}

	if err := c.RPCCall(&block, MethodBlockByNumber, params); err != nil {
		return nil, err
	}

	return &block, nil
}

func (c *Client) GetTransactionReceipts(hash ...string) (TransactionReceipts, error) {
	if len(hash) == 0 {
		return nil, nil
	}

	requestChunks := client.MakeBatchRequests(strings.StringSliceToInterfaces(hash...), rpcBatchSize,
		c.hashToRPCRequestMapper(MethodTransactionReceipt))

	var responses []client.RPCResponse
	for _, requestsChunk := range requestChunks {
		chunkResponses, err := c.RPCBatchCall(requestsChunk)
		if err != nil {
			return nil, err
		}

		responses = append(responses, chunkResponses...)
	}

	responseResults := make([]interface{}, len(responses))
	for i, response := range responses {
		responseResults[i] = response.Result
	}

	responseResultBytes, err := json.Marshal(responseResults)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	var result []TransactionReceipt
	if err = json.Unmarshal(responseResultBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return result, nil
}

func (c *Client) hashToRPCRequestMapper(method string) func(interface{}) client.RPCRequest {
	return func(i interface{}) client.RPCRequest {
		array := []interface{}{i}

		return client.RPCRequest{
			Method: method,
			Params: array,
		}
	}
}

func (c *Client) GetVersion() (string, error) {
	var resp string
	if err := c.RPCCall(&resp, MethodClientVersion, nil); err != nil {
		return "", err
	}

	return resp, nil
}
