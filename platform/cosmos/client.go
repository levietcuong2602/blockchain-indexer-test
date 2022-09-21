package cosmos

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/unanoc/blockchain-indexer/pkg/client"
)

type Client struct {
	client.Request
}

func (c *Client) GetCurrentBlockNumber() (num int64, err error) {
	var resp BlockResponse
	if err = c.Get(&resp, "/blocks/latest", nil); err != nil {
		return num, err
	}

	num, err = strconv.ParseInt(resp.Block.LastCommit.Height, 10, 64)
	if err != nil {
		return num, fmt.Errorf("failed to parse int: %w", err)
	}

	return
}

func (c *Client) GetBlockByNumber(num int64) (TxPage, error) {
	var txs []TxResponse
	var offset int

	for {
		txsPage, err := c.getBlockByNumberWithOffset(num, offset)
		if err != nil {
			return txsPage, err
		}

		txs = append(txs, txsPage.TxResponses...)
		if len(txsPage.TxResponses) == 0 {
			break
		}

		total, err := strconv.ParseInt(txsPage.Pagination.Total, 10, 64)
		if err != nil {
			return txsPage, fmt.Errorf("failed to parse int: %w", err)
		}

		if total <= int64(len(txs)) {
			break
		}

		offset = len(txs)
	}

	return TxPage{TxResponses: txs}, nil
}

func (c *Client) getBlockByNumberWithOffset(num int64, offset int) (txs TxPage, err error) {
	err = c.Get(&txs, "/cosmos/tx/v1beta1/txs", url.Values{
		"events":            {fmt.Sprintf("tx.height=%d", num)},
		"pagination.offset": {strconv.Itoa(offset)},
	})

	return
}

func (c *Client) GetVersion() (string, error) {
	var nodeInfo NodeInfo
	err := c.Get(&nodeInfo, "/node_info", nil)
	if err != nil {
		return "", err
	}

	return nodeInfo.AppVersion.Version, nil
}
