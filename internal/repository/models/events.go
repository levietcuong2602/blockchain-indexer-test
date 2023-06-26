package models

import (
	"github.com/jackc/pgtype"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
	"math/big"
)

// Event - Events emitted from smart contracts to be held in this table
type Event struct {
	BlockHash       string           `gorm:"column:blockhash;type:char(66);not null;primaryKey"`
	Index           uint64           `gorm:"column:index;type:bigint;not null;primaryKey"`
	Origin          string           `gorm:"column:origin;type:char(42);not null;index"`
	Topics          pgtype.TextArray `gorm:"column:topics;type:text[];not null;index:,type:gin"`
	Data            []byte           `gorm:"column:data;type:bytea"`
	TransactionHash string           `gorm:"column:txhash;type:char(66);not null;index"`
}

func NormalizeEvent(e types.Event) (*Event, error) {
	index := new(big.Int)
	iByte, err := e.LogIndex.MarshalJSON()
	if err != nil {
		return nil, err
	}
	index.SetBytes(iByte)

	oBytes, err := e.BlockNumber.MarshalJSON()
	if err != nil {
		return nil, err
	}

	topics := pgtype.TextArray{}
	topics.Set(e.Topics)
	event := Event{
		BlockHash:       e.BlockHash,
		Index:           index.Uint64(),
		Origin:          string(oBytes),
		Data:            []byte(e.Data),
		TransactionHash: e.TransactionHash,
		Topics:          topics,
	}

	return &event, nil
}

func NormalizeEvents(events []types.Event) ([]Event, error) {
	result := make([]Event, len(events))
	for i := range events {
		normalizedEvent, err := NormalizeEvent(events[i])
		if err != nil {
			return nil, err
		}

		result[i] = *normalizedEvent
	}

	return result, nil
}
