package broker

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/channel"

	"github.com/statechannels/go-nitro/channel/state"
	typ "github.com/statechannels/go-nitro/types"
)

var ChainId = big.NewInt(43112)

type Broker struct {
	Address     typ.Address
	Destination typ.Destination
	PrivateKey  []byte
	Role        uint
}

func New(address typ.Address, destination typ.Destination, privateKey []byte, role uint) *Broker {
	return &Broker{Address: address, Destination: destination, PrivateKey: privateKey, Role: role}
}

func (b *Broker) SignTransaction(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	signer := types.LatestSignerForChainID(ChainId)
	hash := signer.Hash(tx)

	prv, err := crypto.ToECDSA(b.PrivateKey)
	if err != nil {
		return nil, err
	}

	signature, err := crypto.Sign(hash.Bytes(), prv)
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(signer, signature)
}

func (b *Broker) SignState(c *channel.Channel, newState state.State) (state.Signature, error) {
	signature, err := newState.Sign(b.PrivateKey)
	ok := c.AddStateWithSignature(newState, signature)
	if err != nil && !ok {
		return signature, err
	}
	return signature, nil
}
