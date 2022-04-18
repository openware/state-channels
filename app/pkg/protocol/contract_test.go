package protocol

import (
	"app/pkg/nitro"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewContract(t *testing.T) {
	contract := NewContract(nitro.Client{}, common.HexToAddress("0x"))

	assert.Equal(t, contract.AssetAddress, common.HexToAddress("0x"))
	assert.Equal(t, contract.Client, nitro.Client{})
}
