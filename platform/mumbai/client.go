package mumbai

import (
	"fmt"
	"github.com/unanoc/blockchain-indexer/pkg/client"
	"math/big"
	"strings"
)

type Client struct {
	client.Request
}

func (c *Client) GetCurrentBlockNumber() (int64, error) {
	blockNumberBytes, _ := c.RPCCallRaw("eth_blockNumber", map[string]string{"finality": "final"})
	// Get the block number from the response
	blockNumberHex := string(blockNumberBytes)
	blockNumberHex = strings.Trim(blockNumberHex, "\"")

	// Convert the block number from hexadecimal to decimal
	blockNumber := new(big.Int)
	blockNumber.SetString(blockNumberHex[2:], 16)
	blockNumber.Int64()
	fmt.Println("Current block number:", blockNumber.Int64())

	return blockNumber.Int64(), nil
}

func (c *Client) GetBlockByNumber(num int64) (ChunkDetail, error) {
	var block Block

	if err := c.RPCCall(&block, "block", map[string]int64{"block_id": num}); err != nil || len(block.Chunks) == 0 {
		return ChunkDetail{}, err
	}

	var chunk ChunkDetail
	if err := c.RPCCall(&chunk, "chunk", []string{block.Chunks[0].Hash}); err != nil {
		return ChunkDetail{}, err
	}

	chunk.Header.Timestamp = block.Header.Timestamp

	return chunk, nil
}

func (c *Client) GetVersion() (string, error) {
	var nodeStatus NodeStatus

	if err := c.RPCCall(&nodeStatus, "status", nil); err != nil {
		return "", err
	}

	return nodeStatus.Version.Version, nil
}
