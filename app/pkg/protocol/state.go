package protocol

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	st "github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
)

// buildState constructs state from input params.
func buildState(
	contract *Contract, participants []*Participant,
	channelNonce *big.Int, appData []byte, turnNum uint64, isFinal bool) st.State {
	addrs := addresses(participants)
	assetExit := singleAssetExit(contract.AssetAddress, participants)

	state := st.State{
		ChainId:           contract.Client.ChainID,
		Participants:      addrs,
		ChannelNonce:      channelNonce,
		ChallengeDuration: big.NewInt(60),
		AppData:           appData,
		Outcome:           outcome.Exit{assetExit},
		TurnNum:           turnNum,
		IsFinal:           isFinal,
	}
	return state
}

// singleAssetExit returns singleAssetExit struct formed from allocations
func singleAssetExit(assetAddress common.Address, participants []*Participant) outcome.SingleAssetExit {
	var allocations []outcome.Allocation

	for _, p := range participants {
		allocations = append(allocations, outcome.Allocation{
			Destination: p.Destination,
			Amount:      p.LockedAmount,
		})
	}

	return outcome.SingleAssetExit{
		Asset:       assetAddress,
		Allocations: allocations,
	}
}

// addresses returns participant's addresses
func addresses(participants []*Participant) []common.Address {
	var addresses []common.Address

	for _, p := range participants {
		addresses = append(addresses, p.Address)
	}

	return addresses
}
