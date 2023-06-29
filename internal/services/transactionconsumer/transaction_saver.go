package transactionconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/unanoc/blockchain-indexer/internal/prometheus"
	"github.com/unanoc/blockchain-indexer/internal/repository/models"
	"github.com/unanoc/blockchain-indexer/internal/repository/postgres"
	"github.com/unanoc/blockchain-indexer/pkg/mq"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/types"
)

const serviceName = "transaction_saver"

type TransactionSaver struct {
	db         *postgres.Database
	prometheus *prometheus.Prometheus
}

func NewTransactionSaver(db *postgres.Database, p *prometheus.Prometheus) mq.MessageProcessor {
	return &TransactionSaver{
		db:         db,
		prometheus: p,
	}
}

func (ts *TransactionSaver) Process(message mq.Message) error {
	var block types.Block

	if err := json.Unmarshal(message, &block); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	log.WithFields(log.Fields{
		"service": serviceName,
		"txs":     len(block.Txs),
	}).Info("Consumed")

	normalizedBlock, err := models.NormalizeBlock(block)
	if err != nil {
		return fmt.Errorf("failed to normalized block: %w", err)
	}
	if err = ts.db.InsertBlock(context.Background(), *normalizedBlock); err != nil {
		return fmt.Errorf("failed to insert block: %w", err)
	}
	if len(block.Txs) == 0 {
		return nil
	}

	chain := block.Txs[0].Chain
	normalizedTxs, err := models.NormalizeTransactions(block.Txs, chain)
	if err != nil {
		return fmt.Errorf("failed to normalized txs: %w", err)
	}

	if err = ts.db.InsertTransactions(context.Background(), normalizedTxs); err != nil {
		return fmt.Errorf("failed to insert txs: %w", err)
	}

	for _, tx := range block.Txs {
		normalizeEvents, err := models.NormalizeEvents(tx.Events)
		if err != nil {
			return fmt.Errorf("failed to normalized event: %w", err)
		}
		if err = ts.db.InsertEvents(context.Background(), normalizeEvents); err != nil {
			return fmt.Errorf("failed to insert events: %w", err)
		}

		ts.HandleDecodeTransaction(tx)
	}

	ts.prometheus.SetParsedTxs(chain, len(normalizedTxs))

	return nil
}

func (ts *TransactionSaver) HandleDecodeTransaction(tx types.Tx) {
	// TODO: Optimize detect smart contract address
	contractAddress := tx.To
	// Check in collection
	log.Println("Handle contract: ", strings.ToLower(contractAddress))
	collection, err := ts.db.FindCollectionByContract(context.Background(), strings.ToLower(contractAddress))
	if err != nil {
		return
	}

	if collection.ID == 0 {
		return
	}

	log.Println("Handle contract: ", contractAddress)

	var transferEvents []types.TransferLog
	for _, vLog := range tx.Events {
		var transferEvent types.TransferLog
		transferEventHash := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

		if strings.Compare(vLog.Topics[0], transferEventHash.Hex()) == 0 && len(vLog.Topics) >= 4 {
			func() {
				transferEvent.From = common.HexToAddress(vLog.Topics[1]).String()
				transferEvent.To = common.HexToAddress(vLog.Topics[2]).String()
				transferEvent.Contract = collection.Contract

				if collection.Standard == "ERC20" {
					transferEvent.Amount = (*types.HexNumber)(common.HexToHash(vLog.Topics[3]).Big())
				} else if collection.Standard == "ERC721" {
					transferEvent.TokenId = (*types.HexNumber)(common.HexToHash(vLog.Topics[3]).Big())
				}
				transferEvents = append(transferEvents, transferEvent)
			}()
		}
	}

	for _, event := range transferEvents {
		ts.processEventInfo(event.From, event.To, event.TokenId, event.Amount, event.Contract)
	}
}

func (ts *TransactionSaver) processEventInfo(from string, to string, tokenId *types.HexNumber, amount *types.HexNumber, contract string) error {
	tokenID := new(big.Int).Set((*big.Int)(tokenId)).Uint64()
	amountNft := new(big.Int).Set((*big.Int)(amount)).Int64()
	nftBalanceFrom, err := ts.db.FindNftBalanceByOwnerContractAndTokenId(context.Background(), from, contract, tokenID)
	if nftBalanceFrom.ID == 0 {
		nftBalanceFrom, err = ts.db.InsertNftBalance(context.Background(), models.NftBalance{
			Contract: contract,
			TokenId:  tokenID,
			Owner:    from,
			Amount:   amountNft,
		})
		if err != nil {
			return fmt.Errorf("failed to insert nft balance from: %w", err)
		}
	}

	nftBalanceTo, err := ts.db.FindNftBalanceByOwnerContractAndTokenId(context.Background(), to, contract, tokenID)
	if nftBalanceTo.ID == 0 {
		nftBalanceTo, err = ts.db.InsertNftBalance(context.Background(), models.NftBalance{
			Contract: contract,
			TokenId:  tokenID,
			Owner:    from,
			Amount:   0,
		})
		if err != nil {
			return fmt.Errorf("failed to insert nft balance to: %w", err)
		}
	}

	//	update amount balances
	nftBalanceFrom.Amount = nftBalanceFrom.Amount - amountNft
	nftBalanceTo.Amount = nftBalanceTo.Amount + amountNft
	ts.db.UpdateNftBalance(context.Background(), *nftBalanceFrom)
	ts.db.UpdateNftBalance(context.Background(), *nftBalanceTo)
	return nil
}
