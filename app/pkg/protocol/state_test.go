package protocol

import (
	"app/pkg/nitro"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

func TestBuildState(t *testing.T) {
	t.Run("build state with several participants", func(t *testing.T) {
		participant1 := NewParticipant(common.HexToAddress("0x01"), types.Destination(common.HexToHash("0x01")), uint(1), big.NewInt(2))
		participant2 := NewParticipant(common.HexToAddress("0x02"), types.Destination(common.HexToHash("0x02")), uint(1), big.NewInt(3))
		contract := NewContract(nitro.Client{}, common.HexToAddress("0x"))
		channelNonce := big.NewInt(time.Now().UnixMilli())

		expectedOutomeExit := outcome.Exit{
			outcome.SingleAssetExit{
				Asset: contract.AssetAddress,
				Allocations: []outcome.Allocation{
					{Destination: participant1.Destination, Amount: participant1.LockedAmount},
					{Destination: participant2.Destination, Amount: participant2.LockedAmount},
				},
			},
		}

		state := buildState(contract, []*Participant{participant1, participant2}, channelNonce, []byte{}, 1, true)
		assert.Equal(t, true, state.IsFinal)
		assert.Equal(t, uint64(1), state.TurnNum)
		assert.Equal(t, []common.Address{participant1.Address, participant2.Address}, state.Participants)
		assert.Equal(t, expectedOutomeExit, state.Outcome)
	})
}
