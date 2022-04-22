package protocol

import (
	"math/big"
	"time"

	st "github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
)

// InitProposal represents information about initial state, contract, participants.
type InitProposal struct {
	Participants []*Participant
	State        *st.State
	Contract     *Contract
	ChannelNonce *big.Int
}

// NewInitProposal returns InitProposal object from income params.
func NewInitProposal(p *Participant, contract *Contract) *InitProposal {
	channelNonce := big.NewInt(time.Now().UnixMilli())
	participants := []*Participant{p}
	// Build initial state, called PreFund state in go-nitro
	state := buildState(contract, participants, channelNonce, []byte{}, 0, false)

	return &InitProposal{
		Contract:     contract,
		ChannelNonce: channelNonce,
		State:        &state,
		Participants: participants,
	}
}

// AddParticipant adds participant into proposed state and participant array.
func (ip *InitProposal) AddParticipant(p *Participant) {
	ip.Participants = append(ip.Participants, p)

	ip.State.Participants = append(ip.State.Participants, p.Address)
	ip.State.Outcome[0].Allocations = append(ip.State.Outcome[0].Allocations,
		outcome.Allocation{
			Destination: p.Destination,
			Amount:      p.LockedAmount,
		})
}
