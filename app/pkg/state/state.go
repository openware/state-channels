package state

import (
	"app/models/broker"
	"app/pkg/contract"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

type ConcludeParams struct {
	OutcomeState  types.Bytes
	AppData       types.Bytes
	Signatures    []contract.IForceMoveSignature
	FixedPart     contract.IForceMoveFixedPart
	NumStates     uint8
	WhoSignedWhat []uint8
}

func Build(
	chainID, channelNonce *big.Int,
	assetAddress common.Address,
	amountForBroker1, amountForBroker2 *big.Int,
	broker1, broker2 *broker.Broker,
	isFinal bool, turnNum uint64) state.State {

	state := state.State{
		ChainId:           chainID,
		Participants:      []types.Address{broker1.Address, broker2.Address},
		ChannelNonce:      channelNonce,
		ChallengeDuration: big.NewInt(60),
		AppData:           []byte{},
		Outcome: outcome.Exit{
			outcome.SingleAssetExit{
				Asset: assetAddress,
				Allocations: outcome.Allocations{
					outcome.Allocation{
						Destination: broker2.Destination,
						Amount:      amountForBroker1,
					},
					outcome.Allocation{
						Destination: broker2.Destination,
						Amount:      amountForBroker2,
					},
				},
			},
		},
		TurnNum: turnNum,
		IsFinal: isFinal,
	}

	return state
}

func BuildConcludeParams(s state.State, signature1, signature2 state.Signature) (ConcludeParams, error) {
	outcomeState, err := s.Outcome.Encode()
	if err != nil {
		return ConcludeParams{}, err
	}

	appData := s.VariablePart().AppData

	var signatureR1, signatureS1, signatureR2, signatureS2 [32]byte

	copy(signatureR1[:], signature1.R)
	copy(signatureS1[:], signature1.S)
	copy(signatureR2[:], signature2.R)
	copy(signatureS2[:], signature2.S)

	sigs := []contract.IForceMoveSignature{
		{V: signature1.V, R: signatureR1, S: signatureS1},
		{V: signature2.V, R: signatureR2, S: signatureS2},
	}

	fixedPart := contract.IForceMoveFixedPart{
		ChainId:           s.ChainId,
		Participants:      s.Participants,
		ChannelNonce:      s.ChannelNonce,
		AppDefinition:     s.AppDefinition,
		ChallengeDuration: s.ChallengeDuration,
	}

	numStates := uint8(len(s.Participants) - 1)
	whoSignedWhat := []uint8{uint8(0), uint8(0)}

	params := ConcludeParams{
		OutcomeState:  outcomeState,
		AppData:       appData,
		Signatures:    sigs,
		FixedPart:     fixedPart,
		NumStates:     numStates,
		WhoSignedWhat: whoSignedWhat,
	}

	return params, nil
}
