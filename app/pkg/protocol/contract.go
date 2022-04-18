package protocol

import (
	"app/pkg/nitro"

	"github.com/ethereum/go-ethereum/common"
)

// Contract stores information about SC client and asset.
type Contract struct {
	Client       nitro.Client
	AssetAddress common.Address
}

// NewContract returns a new Contract from supplied params.
func NewContract(client nitro.Client, assetAddress common.Address) *Contract {
	return &Contract{
		Client:       client,
		AssetAddress: assetAddress,
	}
}
