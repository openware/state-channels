package liability

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	ErrEmptyByteArray         = errors.New("liability: empty byte array")
	ErrNonExistingLiabilities = errors.New("liability: liabilities for such participant don't exist")
	ErrNoPendingLiability     = errors.New("liability: liability with pending type doesn't exist")
	ErrInvalidOperation       = errors.New("liability: given amount is bigger than actual amount")
)

type Asset string

// Liabilities represents information about asset and amount of that asset
type Liabilities struct {
	Pending  map[Asset]decimal.Decimal
	Executed map[Asset]decimal.Decimal
}

// LiabilitiesMap represents information about participant index and appropriate Liability
type LiabilitiesMap map[uint]*Liabilities

// LiabilitiesState represents information about participant index and appropriate LiabilitiesMap
// Example:
// FROM: {
// 	TO: { Pending: { ETH: 1, BTC: 0.1 }, Executed: { ETH: 7 } }
// }
type LiabilitiesState map[uint]LiabilitiesMap

// NewLiabilities creates new Liabilities instance.
func NewLiabilities() *Liabilities {
	pending := make(map[Asset]decimal.Decimal)
	executed := make(map[Asset]decimal.Decimal)

	return &Liabilities{
		Pending:  pending,
		Executed: executed,
	}
}

// AddPendingLiability add pending liability.
func (l *Liabilities) AddPendingLiability(asset Asset, amount decimal.Decimal) {
	l.Pending[asset] = l.Pending[asset].Add(amount)
}

// AddExecutedLiability add executed liability.
func (l *Liabilities) AddExecutedLiability(asset Asset, amount decimal.Decimal) error {
	err := l.validate(asset, amount)
	if err != nil {
		return err
	}

	l.Executed[asset] = amount
	if l.Pending[asset].Equal(amount) {
		delete(l.Pending, asset)
	} else {
		l.Pending[asset] = l.Pending[asset].Sub(amount)
	}

	return nil
}

// AddRevertLiability reverts liability.
func (l *Liabilities) AddRevertLiability(asset Asset, amount decimal.Decimal) error {
	err := l.validate(asset, amount)
	if err != nil {
		return err
	}

	l.Pending[asset] = l.Pending[asset].Sub(amount)
	if l.Pending[asset].Equal(decimal.Zero) {
		delete(l.Pending, asset)
	}

	return nil
}

// validate validates given params for further operations.
func (l *Liabilities) validate(asset Asset, amount decimal.Decimal) error {
	if _, ok := l.Pending[asset]; !ok {
		return ErrNoPendingLiability
	}

	if amount.Cmp(l.Pending[asset]) == 1 {
		return ErrInvalidOperation
	}

	return nil
}

// AddPendingLiability adds new pending liability to liabilities state.
func (ls LiabilitiesState) AddPendingLiability(from, to uint, asset Asset, amount decimal.Decimal) {
	liabilitiesMap, found := ls[from]
	if !found {
		ls[from] = make(LiabilitiesMap)
	}

	_, found = liabilitiesMap[to]
	if !found {
		ls[from][to] = NewLiabilities()
	}

	ls[from][to].AddPendingLiability(asset, amount)
}

// AddExecutedLiability adds new executed liability to liabilities state.
func (ls LiabilitiesState) AddExecutedLiability(from, to uint, asset Asset, amount decimal.Decimal) error {
	_, found := ls[from][to]
	if !found {
		return ErrNonExistingLiabilities
	}

	err := ls[from][to].AddExecutedLiability(asset, amount)
	if err != nil {
		return err
	}

	return nil
}

// AddRevertLiability adds new revert liability to liabilities state.
func (ls LiabilitiesState) AddRevertLiability(from, to uint, asset Asset, amount decimal.Decimal) error {
	_, found := ls[from][to]
	if !found {
		return ErrNonExistingLiabilities
	}

	err := ls[from][to].AddRevertLiability(asset, amount)
	if err != nil {
		return err
	}

	return nil
}

// Print prints pretty LiabilitiesState.
func (ls LiabilitiesState) Print() {
	for index, liabilitiesMap := range ls {
		fmt.Printf("From Participant [%d] ", index)
		for pIndex, liabilities := range liabilitiesMap {
			fmt.Printf("To Participant [%d] %+v\n", pIndex, liabilities)
		}
	}
}

// EncodeToBytes tranform liabilitiesState struct to bytes.
func (ls LiabilitiesState) EncodeToBytes() ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(ls)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

// DecodeFromBytes tranform bytes to liabilitiesState struct.
func DecodeFromBytes(liabilityData []byte) (LiabilitiesState, error) {
	if len(liabilityData) == 0 {
		return LiabilitiesState{}, ErrEmptyByteArray
	}

	buf := bytes.NewBuffer(liabilityData)
	dec := gob.NewDecoder(buf)
	var liabilitiesState LiabilitiesState

	if err := dec.Decode(&liabilitiesState); err != nil {
		return LiabilitiesState{}, err
	}

	return liabilitiesState, nil
}
