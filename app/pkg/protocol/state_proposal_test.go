package protocol

import (
	"app/internal/liability"
	"app/pkg/nitro"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

func getStateProposal() (StateProposal, error) {
	participant := NewParticipant(common.HexToAddress("0x01"), types.Destination(common.HexToHash("0x01")), uint(1), big.NewInt(2))
	contract := NewContract(nitro.Client{}, common.HexToAddress("0x"))
	proposal := NewInitProposal(*participant, *contract)
	state := proposal.State

	newStateProposal, err := NewStateProposal(state)
	return *newStateProposal, err
}

func TestSetFinal(t *testing.T) {
	stateProposal, err := getStateProposal()
	assert.NoError(t, err)
	stateProposal.SetFinal()

	assert.Equal(t, true, stateProposal.IsFinal())
}

func TestSetAppData(t *testing.T) {
	stateProposal, err := getStateProposal()
	assert.NoError(t, err)
	assert.Equal(t, types.Bytes(types.Bytes{}), stateProposal.AppData())

	appData := []byte{1, 2, 3}
	stateProposal.SetAppData(appData)

	expectedResult := types.Bytes(types.Bytes{0x1, 0x2, 0x3})
	assert.Equal(t, expectedResult, stateProposal.AppData())
}

func TestAddLiability(t *testing.T) {
	t.Run("add liability to state without app data", func(t *testing.T) {
		stateProposal, err := getStateProposal()
		assert.NoError(t, err)
		assert.Equal(t, types.Bytes(types.Bytes{}), stateProposal.AppData())

		stateProposal.RequestLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		stateProposal.ApproveLiabilities()
		assert.NotEmpty(t, stateProposal.AppData())
	})

	t.Run("add liability to state with app data", func(t *testing.T) {
		stateProposal, err := getStateProposal()
		assert.NoError(t, err)
		assert.Equal(t, types.Bytes(types.Bytes{}), stateProposal.AppData())

		err = stateProposal.RequestLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.NoError(t, err)
		err = stateProposal.RequestLiability(0, 1, "GOLD", decimal.NewFromFloat(1))
		assert.NoError(t, err)
		err = stateProposal.RequestLiability(0, 1, "LTC", decimal.NewFromFloat(2))
		assert.NoError(t, err)
		err = stateProposal.AcknowledgeLiability(0, 1, "LTC", decimal.NewFromFloat(2))
		assert.NoError(t, err)
		err = stateProposal.RevertLiability(0, 1, "GOLD", decimal.NewFromFloat(1))
		assert.NoError(t, err)

		err = stateProposal.ApproveLiabilities()
		assert.NoError(t, err)

		liab, err := liability.DecodeFromBytes(stateProposal.AppData())
		assert.NoError(t, err)

		assert.Equal(t, map[liability.Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, liab.State[0].Liabilities[1].REQ)
		assert.Equal(t, map[liability.Asset]decimal.Decimal{"LTC": decimal.NewFromFloat(2)}, liab.State[0].Liabilities[1].ACK)
	})
}
