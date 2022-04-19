package liability

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	ErrEmptyByteArray       = errors.New("liability: empty byte array")
	ErrNonExistingLiability = errors.New("liability: there is no such asset in REQ type")
	ErrReverseOperation     = errors.New("liability: can't revert operation")
	ErrAcknowledgeOperation = errors.New("liability: acknowledge amount is bigger than actual amount")
)

type Asset string

// Liability represents information about asset and amount of that asset
type Liability struct {
	REQ map[Asset]decimal.Decimal
	ACK map[Asset]decimal.Decimal
}

// LiabilitiesMap represents information about participant index and appropriate Liability
type LiabilitiesMap struct {
	Liabilities map[uint]*Liability
}

// LiabilityState represents information about participant index and appropriate LiabilitiesMap
// Example:
// FROM: {
// 	TO: { REQ: { ETH: 1, BTC: 0.1 }, ACK: { ETH: 7 } }
// }
type LiabilityState struct {
	State map[uint]*LiabilitiesMap
}

// NewLiability creates new Liability instance.
func NewLiability() *Liability {
	req := make(map[Asset]decimal.Decimal)
	ack := make(map[Asset]decimal.Decimal)

	return &Liability{
		REQ: req,
		ACK: ack,
	}
}

// AddRequestLiability adds liability with REQ type.
func (l *Liability) AddRequestLiability(asset Asset, amount decimal.Decimal) {
	l.REQ[asset] = l.REQ[asset].Add(amount)
}

// AddAcknowledgeLiability adds liability with ACK type.
func (l *Liability) AddAcknowledgeLiability(asset Asset, amount decimal.Decimal) error {
	if _, ok := l.REQ[asset]; !ok {
		return ErrNonExistingLiability
	}

	if amount.Cmp(l.REQ[asset]) == 1 {
		return ErrAcknowledgeOperation
	}

	if l.REQ[asset].Equal(amount) {
		l.ACK[asset] = l.REQ[asset]
		delete(l.REQ, asset)
	} else {
		l.REQ[asset] = l.REQ[asset].Sub(amount)
		l.ACK[asset] = amount
	}

	return nil
}

// AddRevertLiability reverts liability with REQ type only.
func (l *Liability) AddRevertLiability(asset Asset, amount decimal.Decimal) error {
	if _, ok := l.REQ[asset]; !ok {
		return ErrNonExistingLiability
	}

	if amount.Cmp(l.REQ[asset]) == 1 {
		return ErrReverseOperation
	}

	l.REQ[asset] = l.REQ[asset].Sub(amount)
	if l.REQ[asset].Equal(decimal.Zero) {
		delete(l.REQ, asset)
	}

	return nil
}

// Add combines 2 liabilities.
// 1: REQ: { BTC: 12, ETH: 14 }, ACK: {RTC: 2}
// 2: REQ: { BTC: 23, GOLD: 25 }, ACK: {USDT: 200}
// Result: REQ: { BTC: 35, ETH: 14, GOLD: 25 }, ACK: {RTC: 2, USDT: 200}}
func (l *Liability) Add(newLiability Liability) {
	for lAsset, lAmount := range newLiability.ACK {
		if val, ok := l.ACK[lAsset]; ok {
			l.ACK[lAsset] = val.Add(lAmount)
		} else {
			l.ACK[lAsset] = lAmount
		}
	}

	for lAsset, lAmount := range newLiability.REQ {
		if val, ok := l.REQ[lAsset]; ok {
			l.REQ[lAsset] = val.Add(lAmount)
		} else {
			l.REQ[lAsset] = lAmount
		}
	}
}

// NewLiabilitiesMap creates new liabilityMap from existing liability.
func NewLiabilitiesMap(index uint, liability Liability) *LiabilitiesMap {
	liabilitiesMap := make(map[uint]*Liability)
	liabilitiesMap[index] = &liability

	return &LiabilitiesMap{
		Liabilities: liabilitiesMap,
	}
}

// Add combines liabilityMap with already existing one.
func (lm *LiabilitiesMap) Add(liabilitiesMap LiabilitiesMap) {
	for index, liability := range liabilitiesMap.Liabilities {
		if lm.Liabilities[index] == nil {
			lm.Liabilities[index] = liability
		} else {
			lm.Liabilities[index].Add(*liability)
		}
	}
}

// NewLiabilityState creates new liabilityState.
func NewLiabilityState() *LiabilityState {
	state := make(map[uint]*LiabilitiesMap)

	return &LiabilityState{
		State: state,
	}
}

// AddLiability combines existing liabilities map with input liability.
func (ls *LiabilityState) AddLiability(from, to uint, liability Liability) {
	lMap := NewLiabilitiesMap(to, liability)
	if ls.State[from] == nil {
		ls.State[from] = lMap
	} else {
		ls.State[from].Add(*lMap)
	}
}

// AddRequestLiability adds new request liability to liability state.
func (ls *LiabilityState) AddRequestLiability(from, to uint, asset Asset, amount decimal.Decimal) error {
	liability := NewLiability()
	liability.AddRequestLiability(asset, amount)
	ls.AddLiability(from, to, *liability)

	return nil
}

// AddAcknowledgeLiability adds new acknowledge liability to liability state.
func (ls *LiabilityState) AddAcknowledgeLiability(from, to uint, asset Asset, amount decimal.Decimal) error {
	if ls.State[from] != nil && ls.State[from].Liabilities[to] != nil {
		liability := ls.State[from].Liabilities[to]
		err := liability.AddAcknowledgeLiability(asset, amount)
		if err != nil {
			return nil
		}
	} else {
		return ErrNonExistingLiability
	}

	return nil
}

// AddRevertLiability adds new revert liability to liability state.
func (ls *LiabilityState) AddRevertLiability(from, to uint, asset Asset, amount decimal.Decimal) error {
	if ls.State[from] != nil && ls.State[from].Liabilities[to] != nil {
		liability := ls.State[from].Liabilities[to]
		err := liability.AddRevertLiability(asset, amount)
		if err != nil {
			return nil
		}
	} else {
		return ErrNonExistingLiability
	}

	return nil
}

// MergeLiabilityState adds liabilityMap to already existing one.
func (ls *LiabilityState) MergeLiabilityState(newLiabilitiesState LiabilityState) *LiabilityState {
	for index, liability := range newLiabilitiesState.State {
		if ls.State[index] == nil {
			ls.State[index] = liability
		} else {
			ls.State[index].Add(*liability)
		}
	}

	return ls
}

// Print prints pretty LiabilityState.
func (ls *LiabilityState) Print() {
	for index, liabilityMap := range ls.State {
		fmt.Printf("From Participant [%d] ", index)
		for pIndex, liability := range liabilityMap.Liabilities {
			fmt.Printf("To Participant [%d] %+v\n", pIndex, liability)
		}
	}
}

// EncodeToBytes tranform liabilityState struct to bytes.
func (ls *LiabilityState) EncodeToBytes() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(ls)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// DecodeFromBytes tranform bytes to liabilityState struct.
func DecodeFromBytes(liabilityData []byte) (*LiabilityState, error) {
	if len(liabilityData) == 0 {
		return &LiabilityState{}, ErrEmptyByteArray
	}

	buf := bytes.NewBuffer(liabilityData)
	dec := gob.NewDecoder(buf)
	var liabilityState LiabilityState

	if err := dec.Decode(&liabilityState); err != nil {
		return &LiabilityState{}, err
	}

	return &liabilityState, nil
}
