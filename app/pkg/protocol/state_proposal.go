package protocol

import (
	"app/internal/liability"
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/shopspring/decimal"
	st "github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

var ErrEmptyByteArray = errors.New("empty byte array")

// StateProposal represents information about proposed state.
type StateProposal struct {
	state            *st.State
	liabilitiesState liability.LiabilitiesState
}

// NewStateProposal creates state proposal from state.
func NewStateProposal(state *st.State) (*StateProposal, error) {
	liabilitiesState, err := liability.DecodeFromBytes(state.AppData)
	if errors.Is(err, liability.ErrEmptyByteArray) {
		liabilitiesState = make(liability.LiabilitiesState)
	} else if err != nil {
		return &StateProposal{}, err
	}

	return &StateProposal{
		state:            state,
		liabilitiesState: liabilitiesState,
	}, nil
}

// TurnNum returns proposed state turn num.
func (sp *StateProposal) TurnNum() uint64 {
	return sp.state.TurnNum
}

// SetFinal sets proposed state to final.
func (sp *StateProposal) SetFinal() {
	sp.state.IsFinal = true
}

// IsFinal returns either proposed state is final or not.
func (sp *StateProposal) IsFinal() bool {
	return sp.state.IsFinal
}

// AppData returns proposed state app data.
func (sp *StateProposal) AppData() types.Bytes {
	return sp.state.AppData
}

// SetFinal sets proposed state with appData.
func (sp *StateProposal) SetAppData(appData []byte) {
	sp.state.AppData = appData
}

// LiabilityState returns proposed state liability.
func (sp *StateProposal) LiabilityState() liability.LiabilitiesState {
	return sp.liabilitiesState
}

// State returns state proposal state.
func (sp *StateProposal) State() st.State {
	return *sp.state
}

// PendingLiability add pending liability to state proposal.
func (sp *StateProposal) PendingLiability(from, to uint, asset liability.Asset, amount decimal.Decimal) {
	sp.liabilitiesState.AddPendingLiability(from, to, asset, amount)
}

// ExecutedLiability add executed liability to state proposal.
func (sp *StateProposal) ExecutedLiability(from, to uint, asset liability.Asset, amount decimal.Decimal) error {
	return sp.liabilitiesState.AddExecutedLiability(from, to, asset, amount)
}

// RevertLiability add revert liability to state proposal.
func (sp *StateProposal) RevertLiability(from, to uint, asset liability.Asset, amount decimal.Decimal) error {
	return sp.liabilitiesState.AddRevertLiability(from, to, asset, amount)
}

// ApproveLiabilities approves all requested liabilities.
func (sp *StateProposal) ApproveLiabilities() error {
	appDataBytes, err := sp.liabilitiesState.EncodeToBytes()
	if err != nil {
		return err
	}

	sp.state.AppData = appDataBytes

	return nil
}

// EncodeToBytes tranform stateProposal struct to bytes.
func (sp *StateProposal) EncodeToBytes() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(sp)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// DecodeStateProposalFromBytes tranform bytes to stateProposal struct.
func DecodeStateProposalFromBytes(stateProposalData []byte) (StateProposal, error) {
	if len(stateProposalData) == 0 {
		return StateProposal{}, ErrEmptyByteArray
	}

	buf := bytes.NewBuffer(stateProposalData)
	dec := gob.NewDecoder(buf)
	var stateProposal StateProposal

	if err := dec.Decode(&stateProposalData); err != nil {
		return StateProposal{}, err
	}

	return stateProposal, nil
}
