package mumbai

import (
	"encoding/json"
	"fmt"
)

const errorMissingBlock = -32000

func (p *Platform) GetCurrentBlockNumber() (int64, error) {
	return p.client.GetCurrentBlockNumber()
}

func (p *Platform) GetBlockByNumber(num int64) ([]byte, error) {
	block, err := p.client.GetBlockByNumber(num)
	if err != nil { // mumbai won't return old blocks
		return nil, err
	}
	if block.Timestamp == nil { // pending block
		return nil, fmt.Errorf("pending block %d... timestamp is nil", num)
	}

	hashes := make([]string, 0, len(block.Transactions))
	for _, tx := range block.Transactions {
		hashes = append(hashes, tx.Hash)
	}
	receipts, err := p.client.GetTransactionReceipts(hashes...)
	if err != nil {
		return nil, err
	}
	block.TxReceipts = receipts

	data, err := json.Marshal(block)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal block: %w", err)
	}
	return data, nil
}
