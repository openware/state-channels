package protocol

import (
	"app/pkg/eth/gasprice"
	"bytes"
	"encoding/gob"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	chl "github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
)

var (
	ErrCompletedState       = errors.New("channel: already completed state")
	ErrNotFinalState        = errors.New("channel: not final state")
	ErrInvalidSignature     = errors.New("channel: invalid signature")
	ErrSignatureIsNotInList = errors.New("channel: signature is not in participant list")
	ErrInvalidAmount        = errors.New("channel: fund amount is different from initial outcome allocation amount")
	ErrIncompleteState      = errors.New("channel: incomplete state")
)

// Channel represents information about current state, channel info.
type Channel struct {
	initProposal *InitProposal
	lastState    *state.State
	c            chl.Channel
}

// InitChannel opens channel with participant who was requested opening a channel.
func InitChannel(initProposal *InitProposal, participantIndex uint) (*Channel, error) {
	c, err := chl.New(*initProposal.State, participantIndex)
	if err != nil {
		return &Channel{}, err
	}

	return &Channel{
		c:            c,
		initProposal: initProposal,
		lastState:    initProposal.State,
	}, nil
}

// ApproveChannelInit add participant's signature to the prefund state.
// It returns signed state signature.
func (channel *Channel) ApproveInitChannel(privateKey []byte) (state.Signature, error) {
	if channel.c.PreFundComplete() {
		return state.Signature{}, ErrCompletedState
	}

	preFundState := channel.c.PreFundState()
	signature, err := channel.signState(&preFundState, privateKey)
	if err != nil {
		return state.Signature{}, err
	}

	return signature, nil
}

// FundChannel deposits funds to already opened state channel.
// It returns on-chain transaction with detailed information.
func (channel *Channel) FundChannel(p *Participant, privateKey []byte, opts ...gasprice.Station) (*types.Transaction, error) {
	if !channel.c.PreFundComplete() {
		return &types.Transaction{}, ErrIncompleteState
	}

	contract := channel.initProposal.Contract
	adjudicator := contract.Client.Adjudicator
	signerFn := signTransaction(channel.c.ChainId, privateKey)

	// Construct TransactionOpts based on options
	var transactOpts bind.TransactOpts
	if len(opts) == 0 {
		transactOpts = bind.TransactOpts{From: p.Address, Signer: signerFn, Value: p.LockedAmount}
	} else {
		transactOpts = bind.TransactOpts{
			GasPrice: opts[0].GasPrice, GasLimit: opts[0].GasLimit,
			From: p.Address, Signer: signerFn, Value: p.LockedAmount,
		}
	}

	expectedHeld, err := channel.CheckHoldings()
	if err != nil {
		return &types.Transaction{}, err
	}

	transaction, err := adjudicator.Deposit(&transactOpts, contract.AssetAddress, channel.c.Id, expectedHeld, p.LockedAmount)
	if err != nil {
		return &types.Transaction{}, err
	}

	return transaction, nil
}

// ApproveChannelFunding signs postfund state after funding channel.
// It returns signed state signature.
func (channel *Channel) ApproveChannelFunding(privateKey []byte) (state.Signature, error) {
	if channel.c.PostFundComplete() {
		return state.Signature{}, ErrCompletedState
	}

	postFundState := channel.c.PostFundState()
	signature, err := channel.signState(&postFundState, privateKey)
	if err != nil {
		return state.Signature{}, err
	}

	if !postFundState.Equal(*channel.lastState) {
		channel.lastState = &postFundState
	}

	return signature, nil
}

// ProposeState constructs new state with specified liability structure and signs proposed state
// by participant who initiated proposal, returns this state proposal.
// TODO
// Now system could generate as many states as can
// Need to block ability to generate new state without approving prev one
func (channel *Channel) ProposeState() (*StateProposal, error) {
	lastStateNum := uint64(channel.lastState.TurnNum + 1)
	channel.lastState.TurnNum = lastStateNum
	stProposal, err := NewStateProposal(channel.lastState)
	if err != nil {
		return &StateProposal{}, err
	}

	return stProposal, nil
}

// SignState adds a participant's signature to the proposed state and returns signed state signature.
// An error is thrown if the signature is invalid.
func (channel *Channel) SignState(stateProposal *StateProposal, privateKey []byte) (state.Signature, error) {
	signature, err := channel.signState(stateProposal.state, privateKey)
	if err != nil {
		return state.Signature{}, err
	}

	// if participant agrees only on specific state, system need to update last state in agreement
	if !channel.lastState.Equal(*stateProposal.state) {
		channel.lastState = stateProposal.state
	}

	return signature, nil
}

// Conclude transfer all participants funds to the destination addresses and close state channel.
// It returns on-chain transaction with detailed information.
func (channel *Channel) Conclude(p *Participant, privateKey []byte, participantSignatures map[common.Address]state.Signature, opts ...gasprice.Station) (*types.Transaction, error) {
	lastState := channel.lastState
	if !lastState.IsFinal {
		return &types.Transaction{}, ErrNotFinalState
	}

	signerFn := signTransaction(channel.c.ChainId, privateKey)
	adjudicator := channel.initProposal.Contract.Client.Adjudicator

	// Construct TransactionOpts based on options
	var transactOpts bind.TransactOpts
	if len(opts) == 0 {
		transactOpts = bind.TransactOpts{From: p.Address, Signer: signerFn}
	} else {
		transactOpts = bind.TransactOpts{GasPrice: opts[0].GasPrice, GasLimit: opts[0].GasLimit, From: p.Address, Signer: signerFn}
	}

	finalTurnNum := big.NewInt(int64(channel.lastState.TurnNum))
	concludeParams, err := buildConcludeParams(lastState, participantSignatures)
	if err != nil {
		return &types.Transaction{}, err
	}

	concludeTransaction, err := adjudicator.ConcludeAndTransferAllAssets(&transactOpts,
		finalTurnNum,
		concludeParams.FixedPart,
		concludeParams.AppData,
		concludeParams.OutcomeState,
		concludeParams.NumStates,
		concludeParams.WhoSignedWhat,
		concludeParams.Signatures,
	)

	if err != nil {
		return &types.Transaction{}, err
	}

	return concludeTransaction, nil
}

// CheckSignature returns true if signature is valid, existing in state channel participant list and
// connected to the specific state, false otherwise.
func (channel *Channel) CheckSignature(signature state.Signature, s *state.State) (bool, error) {
	address, err := s.RecoverSigner(signature)
	if err != nil {
		return false, err
	}

	var result bool = false
	for _, addr := range s.Participants {
		if addr == address {
			result = true
			break
		}
	}

	if !result {
		return false, ErrSignatureIsNotInList
	}

	return true, nil
}

// CurrentState returns information about current state.
func (channel *Channel) CurrentState() state.State {
	return *channel.lastState
}

// CheckHoldings returns current holdings for already opened state channel per asset.
func (channel *Channel) CheckHoldings() (*big.Int, error) {
	channelID := channel.c.Id
	contract := channel.initProposal.Contract
	adjudicator := contract.Client.Adjudicator

	holdings, err := adjudicator.Holdings(&bind.CallOpts{}, channel.initProposal.Contract.AssetAddress, channelID)
	if err != nil {
		return nil, err
	}

	return holdings, err
}

// StateIsFinal returns true if current state is final, false otherwise.
func (channel *Channel) StateIsFinal() bool {
	return channel.lastState.IsFinal
}

// EncodeToBytes tranform Channel struct to bytes.
func (channel *Channel) EncodeToBytes() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(channel)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// DecodeChannelFromBytes tranform bytes to Channel struct.
func DecodeChannelFromBytes(channelData []byte) (Channel, error) {
	if len(channelData) == 0 {
		return Channel{}, ErrEmptyByteArray
	}

	buf := bytes.NewBuffer(channelData)
	dec := gob.NewDecoder(buf)
	var channel Channel

	if err := dec.Decode(&channelData); err != nil {
		return Channel{}, err
	}

	return channel, nil
}

// signState adds a participant's signature to the newState.
// An error is thrown if the signature is invalid.
func (channel *Channel) signState(newState *state.State, privateKey []byte) (state.Signature, error) {
	signature, err := newState.Sign(privateKey)
	if err != nil {
		return state.Signature{}, err
	}

	ok := channel.c.AddStateWithSignature(*newState, signature)
	if !ok {
		return state.Signature{}, ErrInvalidSignature
	}

	return signature, nil
}
