package protocol

import (
	"math/big"

	"github.com/statechannels/go-nitro/types"
)

// Participant stores information about participant.
type Participant struct {
	Address      types.Address
	Destination  types.Destination
	LockedAmount *big.Int
	Index        uint
}

// NewParticipant returns a new Participant from supplied params.
func NewParticipant(address types.Address, destination types.Destination, index uint, lockedAmount *big.Int) *Participant {
	return &Participant{
		Address:      address,
		Destination:  destination,
		Index:        index,
		LockedAmount: lockedAmount,
	}
}
