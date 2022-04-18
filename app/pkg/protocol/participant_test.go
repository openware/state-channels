package protocol

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

func TestNewParticipant(t *testing.T) {
	participant := NewParticipant(common.HexToAddress("0x01"), types.Destination(common.HexToHash("0x01")), uint(1), big.NewInt(2))

	assert.NotEmpty(t, participant)
}
