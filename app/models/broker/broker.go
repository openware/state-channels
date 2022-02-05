package broker

import (
	"github.com/statechannels/go-nitro/types"
)

type Broker struct {
	Address     types.Address
	Destination types.Destination
	PrivateKey  []byte
	Role        uint
}

func New(address types.Address, destination types.Destination, privateKey []byte, role uint) *Broker {
	return &Broker{Address: address, Destination: destination, PrivateKey: privateKey, Role: role}
}
