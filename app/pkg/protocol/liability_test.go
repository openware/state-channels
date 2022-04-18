package protocol

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestLiability(t *testing.T) {
	t.Run("new liability", func(t *testing.T) {
		liability := NewLiability()
		assert.NotEmpty(t, liability)
	})

	t.Run("request liability", func(t *testing.T) {
		liability := NewLiability()
		assert.NotEmpty(t, liability)

		liability.AddRequestLiability("ETH", decimal.NewFromFloat(1))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(1)}, liability.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liability.ACK)

		liability.AddRequestLiability("ETH", decimal.NewFromFloat(2))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(3)}, liability.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liability.ACK)

		liability.AddRequestLiability("BTC", decimal.NewFromFloat(12))
		assert.Equal(t, map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(3),
			"BTC": decimal.NewFromFloat(12),
		}, liability.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liability.ACK)
	})

	t.Run("acknowledge liability", func(t *testing.T) {
		liability := NewLiability()
		assert.NotEmpty(t, liability)

		// Acknowledge, no req asset existing
		err := liability.AddAcknowledgeLiability("ETH", decimal.NewFromFloat(1))
		assert.Error(t, err, ErrNonExistingLiability)

		// Acknowledge, same amount
		liability.AddRequestLiability("ETH", decimal.NewFromFloat(2))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, liability.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liability.ACK)

		err = liability.AddAcknowledgeLiability("ETH", decimal.NewFromFloat(2))
		assert.NoError(t, err)
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, liability.ACK)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liability.REQ)

		// Acknowledge, bigger amount
		liability.AddRequestLiability("BTC", decimal.NewFromFloat(22))
		assert.Equal(t, map[Asset]decimal.Decimal{"BTC": decimal.NewFromFloat(22)}, liability.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, liability.ACK)

		err = liability.AddAcknowledgeLiability("BTC", decimal.NewFromFloat(23))
		assert.Error(t, err, ErrAcknowledgeOperation)
	})

	t.Run("revert liability", func(t *testing.T) {
		liability := NewLiability()
		assert.NotEmpty(t, liability)

		liability.AddRequestLiability("ETH", decimal.NewFromFloat(1))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(1)}, liability.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liability.ACK)

		// Revert non existing liability
		err := liability.AddRevertLiability("BTC", decimal.NewFromFloat(1))
		assert.Error(t, err, ErrNonExistingLiability)

		// Revert liability with bigger amount that actual one
		err = liability.AddRevertLiability("ETH", decimal.NewFromFloat(3))
		assert.Error(t, err, ErrReverseOperation)

		// Revert liability which was acknowledged
		err = liability.AddAcknowledgeLiability("ETH", decimal.NewFromFloat(1))
		assert.NoError(t, err)
		err = liability.AddRevertLiability("ETH", decimal.NewFromFloat(1))
		assert.Error(t, err, ErrReverseOperation)

		// Successfull revert of liability
		liability.AddRequestLiability("BTC", decimal.NewFromFloat(1))
		assert.Equal(t, map[Asset]decimal.Decimal{"BTC": decimal.NewFromFloat(1)}, liability.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(1)}, liability.ACK)

		err = liability.AddRevertLiability("BTC", decimal.NewFromFloat(1))
		assert.NoError(t, err)
		assert.Equal(t, map[Asset]decimal.Decimal{}, liability.REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(1)}, liability.ACK)
	})

	t.Run("add liability", func(t *testing.T) {
		// first liability
		liability1 := NewLiability()
		assert.NotEmpty(t, liability1)

		liability1.AddRequestLiability("ETH", decimal.NewFromFloat(1))
		liability1.AddRequestLiability("BTC", decimal.NewFromFloat(0.2))
		liability1.AddRequestLiability("LTC", decimal.NewFromFloat(2))

		err := liability1.AddAcknowledgeLiability("BTC", decimal.NewFromFloat(0.2))
		assert.NoError(t, err)

		// second liability
		liability2 := NewLiability()
		assert.NotEmpty(t, liability2)

		liability2.AddRequestLiability("ETH", decimal.NewFromFloat(1))
		liability2.AddRequestLiability("USDT", decimal.NewFromFloat(1))
		liability2.AddRequestLiability("LTC", decimal.NewFromFloat(2))

		err = liability2.AddAcknowledgeLiability("LTC", decimal.NewFromFloat(2))
		assert.NoError(t, err)

		liability1.Add(liability2)

		expectedResult := Liability{
			REQ: map[Asset]decimal.Decimal{
				"ETH":  decimal.NewFromFloat(2),
				"LTC":  decimal.NewFromFloat(2),
				"USDT": decimal.NewFromFloat(1),
			},
			ACK: map[Asset]decimal.Decimal{
				"BTC": decimal.NewFromFloat(0.2),
				"LTC": decimal.NewFromFloat(2),
			},
		}

		fmt.Println(liability1.ACK)
		fmt.Println(liability1.REQ)
		assert.Equal(t, expectedResult.ACK, liability1.ACK)
		assert.Equal(t, expectedResult.REQ, liability1.REQ)
	})
}

func TestLiabilityMap(t *testing.T) {
	t.Run("creates new liabilityMap", func(t *testing.T) {
		liability := NewLiability()
		liability.AddRequestLiability("BTC", decimal.NewFromFloat(1))
		liabilityMap := NewLiabilitiesMap(0, liability)

		assert.NotEmpty(t, liabilityMap)
		assert.Equal(t, liability, liabilityMap.Liabilities[0])
	})

	t.Run("adds liabilityMap", func(t *testing.T) {
		liability1 := NewLiability()
		liability2 := NewLiability()
		liability3 := NewLiability()

		liability1.AddRequestLiability("ETH", decimal.NewFromFloat(0.4))
		liability2.AddRequestLiability("BTC", decimal.NewFromFloat(0.2))
		liability3.AddRequestLiability("LTC", decimal.NewFromFloat(10))

		err := liability1.AddAcknowledgeLiability("ETH", decimal.NewFromFloat(0.4))
		assert.NoError(t, err)
		err = liability3.AddAcknowledgeLiability("LTC", decimal.NewFromFloat(10))
		assert.NoError(t, err)

		liabilityMap1 := NewLiabilitiesMap(0, liability1)
		liabilityMap2 := NewLiabilitiesMap(1, liability2)
		liabilityMap3 := NewLiabilitiesMap(0, liability3)

		liabilityMap1.Add(liabilityMap2)
		liabilityMap1.Add(liabilityMap3)

		resultLiability := liability1
		resultLiability.Add(liability3)

		assert.Equal(t, len(liabilityMap1.Liabilities), 2)
		assert.Equal(t, liability2, liabilityMap1.Liabilities[1])
		assert.Equal(t, resultLiability, liabilityMap1.Liabilities[0])
	})
}

func TestLiabilityState(t *testing.T) {
	t.Run("creates liability state", func(t *testing.T) {
		state := NewLiabilityState()
		assert.NotEmpty(t, state)
	})

	t.Run("adds liabilities", func(t *testing.T) {
		state1 := NewLiabilityState()
		assert.NotEmpty(t, state1)

		liability1 := NewLiability()
		liability1.AddRequestLiability("BTC", decimal.NewFromFloat(2))
		liability1.AddRequestLiability("ETH", decimal.NewFromFloat(0.1))

		state1.AddLiability(0, 1, liability1)
		assert.Equal(t, liability1, state1.State[0].Liabilities[1])

		liability2 := NewLiability()
		liability2.AddRequestLiability("BTC", decimal.NewFromFloat(2))
		liability2.AddRequestLiability("LTC", decimal.NewFromFloat(0.1))

		state1.AddLiability(0, 1, liability2)

		liability1.Add(liability2)
		assert.Equal(t, liability1, state1.State[0].Liabilities[1])
	})

	t.Run("adds liability state", func(t *testing.T) {
		liability1 := NewLiability()
		liability1.AddRequestLiability("ETH", decimal.NewFromFloat(0.4))
		state1 := NewLiabilityState()
		state1.AddLiability(0, 1, liability1)

		liability2 := NewLiability()
		liability2.AddRequestLiability("BTC", decimal.NewFromFloat(2))
		state2 := NewLiabilityState()
		state2.AddLiability(0, 2, liability2)

		state1.MergeLiabilityState(state2)

		assert.Equal(t, liability1, state1.State[0].Liabilities[1])
		assert.Equal(t, liability2, state1.State[0].Liabilities[2])

		liability3 := NewLiability()
		liability3.AddRequestLiability("BTC", decimal.NewFromFloat(2))
		state3 := NewLiabilityState()
		state3.AddLiability(0, 2, liability3)

		liability2.Add(liability3)
		state1.MergeLiabilityState(state3)
		assert.Equal(t, liability2, state1.State[0].Liabilities[2])
	})

	t.Run("adds request liabilities", func(t *testing.T) {
		state := NewLiabilityState()

		state.AddRequestLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, state.State[0].Liabilities[1].REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state.State[0].Liabilities[1].ACK)

		state.AddRequestLiability(0, 1, "ETH", decimal.NewFromFloat(3))
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(5)}, state.State[0].Liabilities[1].REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state.State[0].Liabilities[1].ACK)

		state.AddRequestLiability(0, 1, "BTC", decimal.NewFromFloat(0.3))
		assert.Equal(t, map[Asset]decimal.Decimal{
			"ETH": decimal.NewFromFloat(5),
			"BTC": decimal.NewFromFloat(0.3)},
			state.State[0].Liabilities[1].REQ)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state.State[0].Liabilities[1].ACK)
	})

	t.Run("adds acknowledge liability", func(t *testing.T) {
		state := NewLiabilityState()

		// there is no to/from field
		err := state.AddAcknowledgeLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.Error(t, err, ErrNonExistingLiability)

		state.AddRequestLiability(0, 1, "ETH", decimal.NewFromFloat(2))

		// there is no to field
		err = state.AddAcknowledgeLiability(0, 4, "ETH", decimal.NewFromFloat(2))
		assert.Error(t, err, ErrNonExistingLiability)

		err = state.AddAcknowledgeLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.NoError(t, err)
		assert.Equal(t, map[Asset]decimal.Decimal{"ETH": decimal.NewFromFloat(2)}, state.State[0].Liabilities[1].ACK)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state.State[0].Liabilities[1].REQ)
	})

	t.Run("adds revert liability", func(t *testing.T) {
		state := NewLiabilityState()

		// there is no to/from field
		err := state.AddRevertLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.Error(t, err, ErrNonExistingLiability)

		state.AddRequestLiability(0, 1, "ETH", decimal.NewFromFloat(2))

		// there is no to field
		err = state.AddRevertLiability(0, 4, "ETH", decimal.NewFromFloat(2))
		assert.Error(t, err, ErrNonExistingLiability)

		err = state.AddRevertLiability(0, 1, "ETH", decimal.NewFromFloat(2))
		assert.NoError(t, err)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state.State[0].Liabilities[1].ACK)
		assert.Equal(t, map[Asset]decimal.Decimal{}, state.State[0].Liabilities[1].REQ)
	})

	t.Run("encode to bytes", func(t *testing.T) {
		liability := NewLiability()
		liability.AddRequestLiability("ETH", decimal.NewFromFloat(0.4))

		state := NewLiabilityState()
		state.AddLiability(0, 1, liability)

		bytes, err := state.EncodeToBytes()

		assert.NoError(t, err)
		assert.NotEmpty(t, bytes)
	})

	t.Run("decode from bytes", func(t *testing.T) {
		liability := NewLiability()
		liability.AddRequestLiability("ETH", decimal.NewFromFloat(0.4))

		state := NewLiabilityState()
		state.AddLiability(0, 1, liability)
		bytes, err := state.EncodeToBytes()
		assert.NoError(t, err)

		decodedState, err := DecodeLiabilityFromBytes(bytes)
		assert.NoError(t, err)
		assert.Equal(t, state, decodedState)
	})
}
