package liability

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestLiabilities(t *testing.T) {
	t.Run("new liabilities", func(t *testing.T) {
		liabilities := NewLiabilities()
		assert.NotEmpty(t, liabilities)
	})

	t.Run("request liability", func(t *testing.T) {
		liabilities := NewLiabilities()
		assert.NotEmpty(t, liabilities)

		liabilities.AddRequestLiability("ETH", decimal.NewFromFloat(1))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(1)}, liabilities.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liabilities.ACK)

		liabilities.AddRequestLiability("ETH", decimal.NewFromFloat(2))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(3)}, liabilities.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liabilities.ACK)

		liabilities.AddRequestLiability("BTC", decimal.NewFromFloat(12))
		assert.Equal(t, map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(3),
			"BTC": decimal.NewFromFloat(12),
		}, liabilities.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liabilities.ACK)
	})

	t.Run("acknowledge liability", func(t *testing.T) {
		liabilities := NewLiabilities()
		assert.NotEmpty(t, liabilities)

		// Acknowledge, no req asset existing
		err := liabilities.AddAcknowledgeLiability("ETH", decimal.NewFromFloat(1))
		assert.Error(t, err, ErrNoReqLiability)

		// Acknowledge, same amount
		liabilities.AddRequestLiability("ETH", decimal.NewFromFloat(2))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, liabilities.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liabilities.ACK)

		err = liabilities.AddAcknowledgeLiability("ETH", decimal.NewFromFloat(2))
		assert.NoError(t, err)
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, liabilities.ACK)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liabilities.REQ)

		// Acknowledge, bigger amount
		liabilities.AddRequestLiability("BTC", decimal.NewFromFloat(22))
		assert.Equal(t, map[Asset]decimal.Decimal{"BTC": decimal.NewFromFloat(22)}, liabilities.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, liabilities.ACK)

		err = liabilities.AddAcknowledgeLiability("BTC", decimal.NewFromFloat(23))
		assert.Error(t, err, ErrInvalidOperation)
	})

	t.Run("revert liability", func(t *testing.T) {
		liabilities := NewLiabilities()
		assert.NotEmpty(t, liabilities)

		liabilities.AddRequestLiability("ETH", decimal.NewFromFloat(1))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(1)}, liabilities.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liabilities.ACK)

		// Revert non existing liability
		err := liabilities.AddRevertLiability("BTC", decimal.NewFromFloat(1))
		assert.Error(t, err, ErrNoReqLiability)

		// Revert liability with bigger amount that actual one
		err = liabilities.AddRevertLiability("ETH", decimal.NewFromFloat(3))
		assert.Error(t, err, ErrInvalidOperation)

		// Revert liability which was acknowledged
		err = liabilities.AddAcknowledgeLiability("ETH", decimal.NewFromFloat(1))
		assert.NoError(t, err)
		err = liabilities.AddRevertLiability("ETH", decimal.NewFromFloat(1))
		assert.Error(t, err, ErrInvalidOperation)

		// Successfull revert of liability
		liabilities.AddRequestLiability("BTC", decimal.NewFromFloat(1))
		assert.Equal(t, map[Asset]decimal.Decimal{"BTC": decimal.NewFromFloat(1)}, liabilities.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(1)}, liabilities.ACK)

		err = liabilities.AddRevertLiability("BTC", decimal.NewFromFloat(1))
		assert.NoError(t, err)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liabilities.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(1)}, liabilities.ACK)
	})
}

func TestLiabilitiesState(t *testing.T) {
	t.Run("adds request liabilities", func(t *testing.T) {
		state := make(LiabilitiesState)

		state.AddRequestLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, state[0][1].REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state[0][1].ACK)

		state.AddRequestLiability(0, 1, "ETH", decimal.NewFromFloat(3))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(5)}, state[0][1].REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state[0][1].ACK)

		state.AddRequestLiability(0, 1, "BTC", decimal.NewFromFloat(0.3))
		assert.Equal(t, map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(5),
			"BTC": decimal.NewFromFloat(0.3)},
			state[0][1].REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state[0][1].ACK)
	})

	t.Run("adds acknowledge liability", func(t *testing.T) {
		state := make(LiabilitiesState)

		// there is no to/from field
		err := state.AddAcknowledgeLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.Error(t, err, ErrNonExistingLiabilities)

		state.AddRequestLiability(0, 1, "ETH", decimal.NewFromFloat(2))

		// there is no to field
		err = state.AddAcknowledgeLiability(0, 4, "ETH", decimal.NewFromFloat(2))
		assert.Error(t, err, ErrNonExistingLiabilities)

		err = state.AddAcknowledgeLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.NoError(t, err)
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, state[0][1].ACK)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state[0][1].REQ)
	})

	t.Run("adds revert liability", func(t *testing.T) {
		state := make(LiabilitiesState)

		// there is no to/from field
		err := state.AddRevertLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.Error(t, err, ErrNonExistingLiabilities)

		state.AddRequestLiability(0, 1, "ETH", decimal.NewFromFloat(2))

		// there is no to field
		err = state.AddRevertLiability(0, 4, "ETH", decimal.NewFromFloat(2))
		assert.Error(t, err, ErrNonExistingLiabilities)

		err = state.AddRevertLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.NoError(t, err)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state[0][1].ACK)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state[0][1].REQ)
	})

	t.Run("encode to bytes", func(t *testing.T) {
		state := make(LiabilitiesState)
		state.AddRequestLiability(0, 1, "ETH", decimal.NewFromFloat(0.4))

		bytes, err := state.EncodeToBytes()

		assert.NoError(t, err)
		assert.NotEmpty(t, bytes)
	})

	t.Run("decode from bytes", func(t *testing.T) {
		state := make(LiabilitiesState)
		state.AddRequestLiability(0, 1, "ETH", decimal.NewFromFloat(0.4))

		bytes, err := state.EncodeToBytes()
		assert.NoError(t, err)

		decodedState, err := DecodeFromBytes(bytes)
		assert.NoError(t, err)
		assert.Equal(t, state, decodedState)
	})
}
