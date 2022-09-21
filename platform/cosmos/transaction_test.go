package cosmos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/unanoc/blockchain-indexer/pkg/mock"
	"github.com/unanoc/blockchain-indexer/pkg/primitives/coin"
)

func Test_NormalizeTxs(t *testing.T) {
	testCases := []struct {
		coin           uint
		denom          DenomType
		description    string
		input          string
		expectedOutput string
	}{
		{
			coin:           coin.COSMOS,
			denom:          DenomAtom,
			description:    "transfer",
			input:          "mocks/tx_transfer.json",
			expectedOutput: "mocks/tx_transfer_expected_output.json",
		},
		{
			coin:           coin.COSMOS,
			denom:          DenomAtom,
			description:    "delegate",
			input:          "mocks/tx_delegate.json",
			expectedOutput: "mocks/tx_delegate_expected_output.json",
		},
		{
			coin:           coin.COSMOS,
			denom:          DenomAtom,
			description:    "get rewards",
			input:          "mocks/tx_get_reward.json",
			expectedOutput: "mocks/tx_get_reward_expected_output.json",
		},
		{
			coin:           coin.COSMOS,
			denom:          DenomAtom,
			description:    "undelegate",
			input:          "mocks/tx_undelegate.json",
			expectedOutput: "mocks/tx_undelegate_expected_output.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			var page TxPage
			err := mock.JSONModelFromFilePath(tc.input, &page)
			assert.NoError(t, err)

			p := Init(tc.coin, tc.denom, "")

			result := p.NormalizeBlock(page)
			resultJSON, err := json.Marshal(result)
			assert.NoError(t, err)

			expectedJSON, err := mock.JSONStringFromFilePath(tc.expectedOutput)
			assert.NoError(t, err)

			assert.JSONEq(t, expectedJSON, string(resultJSON))
		})
	}
}
