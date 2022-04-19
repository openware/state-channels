package protocol

import (
	"app/pkg/nitro"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/crypto"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/assert"
)

var (
	participant1 = NewParticipant(common.HexToAddress("0xdd2fd4581271e230360230f9337d5c0430bf44c0"), types.Destination(common.HexToHash("0xdd2fd4581271e230360230f9337d5c0430bf44c0")), uint(0), big.NewInt(2))
	participant2 = NewParticipant(common.HexToAddress("0x8626f6940e2eb28930efb4cef49b2d1f2c9c1199"), types.Destination(common.HexToHash("0x8626f6940e2eb28930efb4cef49b2d1f2c9c1199")), uint(1), big.NewInt(2))
)

func getChannel() (*Channel, error) {
	contract := NewContract(nitro.Client{ChainID: big.NewInt(2)}, common.HexToAddress("0x"))
	proposal := NewInitProposal(*participant1, *contract)
	proposal.AddParticipant(*participant2)

	ch, err := InitChannel(proposal, 0)
	return ch, err
}
func TestInitChannel(t *testing.T) {
	t.Run("successful channel initialization", func(t *testing.T) {
		participant := NewParticipant(common.HexToAddress("0x01"), types.Destination(common.HexToHash("0x01")), uint(1), big.NewInt(2))
		contract := NewContract(nitro.Client{ChainID: big.NewInt(2)}, common.HexToAddress("0x"))
		proposal := NewInitProposal(*participant, *contract)

		ch, err := InitChannel(proposal, 0)

		assert.NoError(t, err)
		assert.Equal(t, proposal, ch.initProposal)
		assert.Equal(t, proposal.State, ch.lastState)
	})

	t.Run("unsuccessful channel initialization", func(t *testing.T) {
		participant := NewParticipant(common.HexToAddress("0x01"), types.Destination(common.HexToHash("0x01")), uint(1), big.NewInt(2))
		contract := NewContract(nitro.Client{}, common.HexToAddress("0x"))
		proposal := NewInitProposal(*participant, *contract)
		proposal.State.TurnNum = uint64(5)

		_, err := InitChannel(proposal, 0)
		assert.Error(t, err)
	})
}

func TestApproveChannelInit(t *testing.T) {
	t.Run("invalid private keys", func(t *testing.T) {
		privKeys := make(map[Participant][]byte)
		ch, err := getChannel()
		assert.NoError(t, err)

		privKeys[*participant1] = []byte{}
		privKeys[*participant2] = []byte{}
		for _, key := range privKeys {
			_, err := ch.ApproveInitChannel(key)
			assert.Error(t, err)
		}
	})

	t.Run("prefund state has been completed", func(t *testing.T) {
		privKeys := make(map[Participant][]byte)
		ch, err := getChannel()
		assert.NoError(t, err)

		privKeys[*participant1] = common.Hex2Bytes("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
		privKeys[*participant2] = common.Hex2Bytes("df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e")

		// Approve channel initialization
		for _, key := range privKeys {
			_, err := ch.ApproveInitChannel(key)
			assert.NoError(t, err)
		}
		assert.Equal(t, uint64(0), ch.lastState.TurnNum)

		// Post fund state
		for _, key := range privKeys {
			_, err := ch.ApproveChannelFunding(key)
			assert.NoError(t, err)
		}

		assert.Equal(t, uint64(1), ch.lastState.TurnNum)

		// try to approve channel initialization
		for _, key := range privKeys {
			_, err := ch.ApproveInitChannel(key)
			assert.Error(t, err, ErrCompletedState)
		}
	})

	t.Run("valid private keys", func(t *testing.T) {
		privKeys := make(map[Participant][]byte)
		ch, err := getChannel()
		assert.NoError(t, err)

		privKeys[*participant1] = common.Hex2Bytes("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
		privKeys[*participant2] = common.Hex2Bytes("df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e")
		for _, key := range privKeys {
			_, err := ch.ApproveInitChannel(key)
			assert.NoError(t, err)
		}
	})
}

func TestFundChannel(t *testing.T) {
	t.Run("prefund state has not been completed", func(t *testing.T) {
		privKeys := make(map[Participant][]byte)
		ch, err := getChannel()
		assert.NoError(t, err)

		privKeys[*participant1] = common.Hex2Bytes("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
		privKeys[*participant2] = common.Hex2Bytes("df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e")

		for p, key := range privKeys {
			_, err := ch.FundChannel(p, key)
			assert.Error(t, err, ErrIncompleteState)
		}
	})
}

func TestApproveChannelFunding(t *testing.T) {
	t.Run("invalid private keys", func(t *testing.T) {
		privKeys := make(map[Participant][]byte)
		ch, err := getChannel()
		assert.NoError(t, err)

		privKeys[*participant1] = common.Hex2Bytes("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
		privKeys[*participant2] = common.Hex2Bytes("df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e")

		// approve channel initialization
		for _, key := range privKeys {
			_, err := ch.ApproveInitChannel(key)
			assert.NoError(t, err)
		}
		assert.Equal(t, uint64(0), ch.lastState.TurnNum)

		privKeys[*participant1] = []byte{}
		privKeys[*participant2] = []byte{}

		// Post fund state
		for _, key := range privKeys {
			_, err := ch.ApproveChannelFunding(key)
			assert.Error(t, err)
		}
	})

	t.Run("invalid state", func(t *testing.T) {
		privKeys := make(map[Participant][]byte)
		ch, err := getChannel()
		assert.NoError(t, err)

		privKeys[*participant1] = common.Hex2Bytes("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
		privKeys[*participant2] = common.Hex2Bytes("df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e")

		// approve channel initialization
		for _, key := range privKeys {
			_, err := ch.ApproveInitChannel(key)
			assert.NoError(t, err)
		}
		assert.Equal(t, uint64(0), ch.lastState.TurnNum)

		// Post fund state
		for _, key := range privKeys {
			_, err := ch.ApproveChannelFunding(key)
			assert.NoError(t, err)
		}
		assert.Equal(t, uint64(1), ch.lastState.TurnNum)

		// try to approve channel funding again
		for _, key := range privKeys {
			_, err := ch.ApproveChannelFunding(key)
			assert.Error(t, err, ErrCompletedState)
		}
	})

	t.Run("successful approval funding channel", func(t *testing.T) {
		privKeys := make(map[Participant][]byte)
		ch, err := getChannel()
		assert.NoError(t, err)
		privKeys[*participant1] = common.Hex2Bytes("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
		privKeys[*participant2] = common.Hex2Bytes("df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e")

		for _, key := range privKeys {
			_, err := ch.ApproveInitChannel(key)
			assert.NoError(t, err)
		}
		assert.Equal(t, uint64(0), ch.lastState.TurnNum)

		// Post fund state
		for _, key := range privKeys {
			_, err := ch.ApproveChannelFunding(key)
			assert.NoError(t, err)
		}
		assert.Equal(t, uint64(1), ch.lastState.TurnNum)
	})
}

func TestProposeState(t *testing.T) {
	privKeys := make(map[Participant][]byte)
	ch, err := getChannel()
	assert.NoError(t, err)

	privKeys[*participant1] = common.Hex2Bytes("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
	privKeys[*participant2] = common.Hex2Bytes("df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e")

	// approve channel initialization
	for _, key := range privKeys {
		_, err := ch.ApproveInitChannel(key)
		assert.NoError(t, err)
	}
	assert.Equal(t, uint64(0), ch.lastState.TurnNum)

	// Post fund state
	for _, key := range privKeys {
		_, err := ch.ApproveChannelFunding(key)
		assert.NoError(t, err)
	}
	assert.Equal(t, uint64(1), ch.lastState.TurnNum)

	t.Run("propose state", func(t *testing.T) {
		_, err = ch.ProposeState()
		assert.NoError(t, err)
		assert.Equal(t, uint64(2), ch.lastState.TurnNum)

		_, err = ch.ProposeState()
		assert.NoError(t, err)
		assert.Equal(t, uint64(3), ch.lastState.TurnNum)
	})
}

func TestSignState(t *testing.T) {
	privKeys := make(map[Participant][]byte)
	ch, err := getChannel()
	assert.NoError(t, err)

	privKeys[*participant1] = common.Hex2Bytes("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
	privKeys[*participant2] = common.Hex2Bytes("df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e")

	// approve channel initialization
	for _, key := range privKeys {
		_, err := ch.ApproveInitChannel(key)
		assert.NoError(t, err)
	}
	assert.Equal(t, uint64(0), ch.lastState.TurnNum)

	// Post fund state
	for _, key := range privKeys {
		_, err := ch.ApproveChannelFunding(key)
		assert.NoError(t, err)
	}
	assert.Equal(t, uint64(1), ch.lastState.TurnNum)

	stateProposal, err := ch.ProposeState()
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), ch.lastState.TurnNum)

	t.Run("invalid keys", func(t *testing.T) {
		privKeys[*participant1] = []byte{}
		privKeys[*participant2] = []byte{}

		for _, key := range privKeys {
			_, err := ch.SignState(stateProposal, key)
			assert.Error(t, err)
		}
	})

	t.Run("successful sstate proposal sign", func(t *testing.T) {
		privKeys[*participant1] = common.Hex2Bytes("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
		privKeys[*participant2] = common.Hex2Bytes("df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e")

		for _, key := range privKeys {
			_, err := ch.SignState(stateProposal, key)
			assert.NoError(t, err)
		}
		assert.Equal(t, uint64(2), ch.lastState.TurnNum)
	})
}

func TestConcludeChannel(t *testing.T) {
	privKeys := make(map[Participant][]byte)
	ch, err := getChannel()
	assert.NoError(t, err)

	privKeys[*participant1] = common.Hex2Bytes("de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0")
	privKeys[*participant2] = common.Hex2Bytes("df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e")

	// approve channel initialization
	for _, key := range privKeys {
		_, err := ch.ApproveInitChannel(key)
		assert.NoError(t, err)
	}
	assert.Equal(t, uint64(0), ch.lastState.TurnNum)

	// Post fund state
	for _, key := range privKeys {
		_, err := ch.ApproveChannelFunding(key)
		assert.NoError(t, err)
	}
	assert.Equal(t, uint64(1), ch.lastState.TurnNum)

	stateProposal, err := ch.ProposeState()
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), ch.lastState.TurnNum)

	participantSignatures := make(map[common.Address]crypto.Signature)
	for p, key := range privKeys {
		signature, err := ch.SignState(stateProposal, key)
		assert.NoError(t, err)
		participantSignatures[p.Address] = signature
	}
	assert.Equal(t, uint64(2), ch.lastState.TurnNum)

	t.Run("not final state", func(t *testing.T) {
		_, err = ch.Conclude(*participant1, privKeys[*participant1], participantSignatures)
		assert.Error(t, err, ErrNotFinalState)
	})
}
