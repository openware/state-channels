package protocol

import (
	"app/pkg/nitro"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

func TestInitProposal(t *testing.T) {
	participant := NewParticipant(common.HexToAddress("0x01"), types.Destination(common.HexToHash("0x01")), uint(1), big.NewInt(2))
	contract := NewContract(nitro.Client{}, common.HexToAddress("0x"))
	proposal := NewInitProposal(*participant, *contract)

	assert.NotEmpty(t, proposal)
	assert.Equal(t, []Participant{*participant}, proposal.Participants)
	assert.Equal(t, *contract, proposal.Contract)
}

func TestAddParticipant(t *testing.T) {
	participant1 := NewParticipant(common.HexToAddress("0x01"), types.Destination(common.HexToHash("0x01")), uint(1), big.NewInt(2))
	participant2 := NewParticipant(common.HexToAddress("0x01"), types.Destination(common.HexToHash("0x01")), uint(1), big.NewInt(2))

	contract := NewContract(nitro.Client{}, common.HexToAddress("0x"))
	proposal := NewInitProposal(*participant1, *contract)
	assert.NotEmpty(t, proposal)

	proposal.AddParticipant(*participant2)
	assert.Equal(t, 2, len(proposal.Participants))
	assert.Equal(t, []Participant{*participant1, *participant2}, proposal.Participants)
	assert.Equal(t, []common.Address{participant1.Address, participant2.Address}, proposal.State.Participants)
}
