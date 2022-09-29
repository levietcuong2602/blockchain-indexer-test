package cosmos

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type (
	TxType       string
	EventType    string
	AttributeKey string
	DenomType    string
)

const (
	// Types of messages
	MsgSend                    TxType = "/cosmos.bank.v1beta1.MsgSend"
	MsgDelegate                TxType = "/cosmos.staking.v1beta1.MsgDelegate"
	MsgUndelegate              TxType = "/cosmos.staking.v1beta1.MsgUndelegate"
	MsgWithdrawDelegatorReward TxType = "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"

	DenomAtom    DenomType = "uatom"
	DenomKava    DenomType = "ukava"
	DenomOsmosis DenomType = "uosmo"
	DenomLuna    DenomType = "uluna"
)

// Current Block
type (
	BlockResponse struct {
		Block Block `json:"block"`
	}

	Block struct {
		LastCommit LastCommit `json:"last_commit"`
	}

	LastCommit struct {
		Height string `json:"height"`
	}
)

// Block By Number
type (
	TxPage struct {
		TxResponses []TxResponse `json:"tx_responses"`
		Pagination  struct {
			Total string `json:"total"`
		}
	}

	TxResponse struct {
		Height    string `json:"height"`
		TxHash    string `json:"txhash"`
		Code      int    `json:"code"`
		Date      string `json:"timestamp"`
		Logs      Logs   `json:"logs"`
		GasWanted string `json:"gas_wanted"`
		GasUsed   string `json:"gas_used"`
		Tx        Tx     `json:"tx"`
	}

	Logs []Log

	Log struct {
		Events []Event `json:"events"`
	}

	Event struct {
		Type       EventType  `json:"type"`
		Attributes Attributes `json:"attributes"`
	}

	Attributes []Attribute

	Attribute struct {
		Key   AttributeKey `json:"key"`
		Value string       `json:"value"`
	}

	Tx struct {
		Body     Body     `json:"body"`
		AuthInfo AuthInfo `json:"auth_info"`
	}

	Body struct {
		Messages []Message `json:"messages"`
		Memo     string    `json:"memo"`
	}

	Message struct {
		MessageValue
	}

	MessageValue interface{}

	AuthInfo struct {
		Fee         Fee          `json:"fee"`
		SignerInfos []SignerInfo `json:"signer_infos"`
	}

	Fee struct {
		Amount []Amount `json:"amount"`
	}

	SignerInfo struct {
		Sequence string `json:"sequence"`
	}

	MessageValueSend struct {
		Type     TxType   `json:"@type"`
		FromAddr string   `json:"from_address"`
		ToAddr   string   `json:"to_address"`
		Amount   []Amount `json:"amount"`
	}

	MessageValueDelegate struct {
		Type          TxType `json:"@type"`
		DelegatorAddr string `json:"delegator_address"`
		ValidatorAddr string `json:"validator_address"`
		Amount        Amount `json:"amount,omitempty"`
	}

	Amount struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
)

func (l Logs) GetWithdrawRewardValue(denom DenomType) string {
	result := int64(0)
	for _, log := range l {
		for _, att := range log.Events {
			if att.Type == "withdraw_rewards" {
				result += att.Attributes.GetWithdrawRewardValue(denom)
			}
		}
	}

	return strconv.FormatInt(result, 10)
}

// GetWithdrawRewardValue logic explanation:
// for all attributes with key "amount", parse the value to get values of provided denom
// example: "179731247ukrw,1174881uluna,1448483umnt", for "uluna" denom value will be 1174881
func (a Attributes) GetWithdrawRewardValue(denom DenomType) int64 {
	result := int64(0)
	for _, att := range a {
		if att.Key == "amount" {
			values := strings.Split(att.Value, ",")
			for _, value := range values {
				idx := strings.IndexByte(value, 'u')
				if idx < 0 {
					continue
				}

				v := value[:idx]
				d := value[idx:]
				if d == string(denom) {
					amount, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						log.WithFields(log.Fields{"denom": denom, "value": value}).
							Error("Invalid amount value for cosmos-like chain")

						continue
					}
					result += amount
				}
			}
		}
	}

	return result
}

// UnmarshalJSON reads different message types
func (m *Message) UnmarshalJSON(buf []byte) error {
	messageType := struct {
		Type TxType `json:"@type"`
	}{}
	if err := json.Unmarshal(buf, &messageType); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	switch messageType.Type {
	case MsgSend:
		var msgSend MessageValueSend
		if err := json.Unmarshal(buf, &msgSend); err != nil {
			return fmt.Errorf("failed to unmarshal json to MessageValueSend type: %w", err)
		}

		m.MessageValue = msgSend
	case MsgDelegate, MsgUndelegate, MsgWithdrawDelegatorReward:
		var msgDelegate MessageValueDelegate
		if err := json.Unmarshal(buf, &msgDelegate); err != nil {
			return fmt.Errorf("failed to unmarshal json to MessageValueDelegate type: %w", err)
		}

		m.MessageValue = msgDelegate
	}

	return nil
}

// MarshalJSON reads different message types
func (m *Message) MarshalJSON() (data []byte, err error) {
	msgSend, ok := m.MessageValue.(MessageValueSend)
	if ok {
		data, err = json.Marshal(msgSend)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}

		return data, nil
	}

	msgDelegate, ok := m.MessageValue.(MessageValueDelegate)
	if ok {
		data, err = json.Marshal(msgDelegate)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}

		return data, nil
	}

	data, err = json.Marshal(m.MessageValue)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	return data, nil
}

type (
	NodeInfo struct {
		AppVersion AppVersion `json:"application_version"`
	}

	AppVersion struct {
		Version string `json:"version"`
	}
)
