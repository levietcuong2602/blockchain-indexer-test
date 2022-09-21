package cosmos

import (
	"encoding/json"
	"fmt"
)

func (p Platform) GetCurrentBlockNumber() (int64, error) {
	return p.client.GetCurrentBlockNumber()
}

func (p Platform) GetBlockByNumber(num int64) ([]byte, error) {
	block, err := p.client.GetBlockByNumber(num)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(block)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal block: %w", err)
	}

	return data, nil
}
