package protocol

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/statechannels/go-nitro/channel/state"
)

// State Proposal represents information about proposed state.
type StateProposal struct {
	State          *state.State
	liabilityState *LiabilityState
}

// NewStateProposal creates state proposal from state.
func NewStateProposal(state *state.State) (*StateProposal, error) {
	liabilityState, err := DecodeLiabilityFromBytes(state.AppData)
	if errors.Is(err, ErrEmptyByteArray) {
		liabilityState = NewLiabilityState()
	} else if err != nil {
		return &StateProposal{}, err
	}

	return &StateProposal{
		State:          state,
		liabilityState: liabilityState,
	}, nil
}

// SetFinal sets proposed state to final.
func (sp *StateProposal) SetFinal() {
	sp.State.IsFinal = true
}

// SetFinal sets proposed state with appData.
func (sp *StateProposal) SetAppData(appData []byte) {
	sp.State.AppData = appData
}

// RequestLiability add request liability to state proposal.
func (sp *StateProposal) RequestLiability(from, to uint, asset Asset, amount decimal.Decimal) error {
	return sp.liabilityState.AddRequestLiability(from, to, asset, amount)
}

// AcknowledgeLiability add acknowledge liability to state proposal.
func (sp *StateProposal) AcknowledgeLiability(from, to uint, asset Asset, amount decimal.Decimal) error {
	return sp.liabilityState.AddAcknowledgeLiability(from, to, asset, amount)
}

// RevertLiability add revert liability to state proposal.
func (sp *StateProposal) RevertLiability(from, to uint, asset Asset, amount decimal.Decimal) error {
	return sp.liabilityState.AddRevertLiability(from, to, asset, amount)
}

// ApproveLiabilities approves all requested liabilities.
func (sp *StateProposal) ApproveLiabilities() error {
	appDataBytes, err := sp.liabilityState.EncodeToBytes()
	if err != nil {
		return err
	}

	sp.State.AppData = appDataBytes

	return nil
}
