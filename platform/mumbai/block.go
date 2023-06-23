package mumbai

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
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

	mapping := make(map[string]Transaction)
	hashes := make([]string, 0, len(block.Transactions))
	for _, tx := range block.Transactions {
		hashes = append(hashes, tx.Hash)

		mapping[tx.Hash] = tx
	}
	receipts, err := p.client.GetTransactionReceipts(hashes...)
	if err != nil {
		return nil, err
	}
	block.TxReceipts = receipts

	for _, receipt := range receipts {
		log.Println("block number:", receipt.BlockNumber)
		log.Println("transaction hash:", receipt.TransactionHash)

		if receipt.TransactionHash == "0x6590ee2f6c2a05ced04763bc4136890d7b560e74089f8e24cc69145d02d2398d" {
			log.Println("ok em Dung", receipt.TransactionHash)

			input := mapping[receipt.TransactionHash].Input

			// Known function signatures
			transferSignature := "0xa9059cbb"     // transfer(address,uint256)
			transferFromSignature := "0x23b872dd" // transferFrom(address,address,uint256)
			approveSignature := "0x095ea7b3"      // approve(address,uint256)

			// Check the input field against known function signatures
			if strings.HasPrefix(input, transferSignature) {
				fmt.Println("Token Transfer: transfer")
			} else if strings.HasPrefix(input, transferFromSignature) {
				fmt.Println("Token Transfer: transferFrom")
			} else if strings.HasPrefix(input, approveSignature) {
				fmt.Println("Token Transfer: approve")
			} else {
				fmt.Println("Unknown Token Transfer Type")
			}
		}

		if receipt.To == "0x4E086f01508d8d2a65c5D18A06959ebe98EeB848" {
			log.Println("ok em Dung", receipt.To)
		}
	}

	data, err := json.Marshal(block)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal block: %w", err)
	}
	return data, nil
}
