package protocol

import (
	"app/pkg/nitro"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// ConcludeParams represents information for state channel finalization.
type ConcludeParams struct {
	OutcomeState  types.Bytes
	AppData       types.Bytes
	Signatures    []nitro.IForceMoveSignature
	FixedPart     nitro.IForceMoveFixedPart
	NumStates     uint8
	WhoSignedWhat []uint8
}

// BuildConcludeParams builds conclude params for state channel finalization.
func BuildConcludeParams(s state.State, participantSignatures map[common.Address]state.Signature) (ConcludeParams, error) {
	outcomeState, err := s.Outcome.Encode()
	if err != nil {
		return ConcludeParams{}, err
	}

	appData := s.VariablePart().AppData
	moveSignatures := forceMoveSignatures(s, participantSignatures)

	fixedPart := nitro.IForceMoveFixedPart{
		ChainId:           s.ChainId,
		Participants:      s.Participants,
		ChannelNonce:      s.ChannelNonce,
		AppDefinition:     s.AppDefinition,
		ChallengeDuration: s.ChallengeDuration,
	}

	// TODO
	// Now system supports only positive case, if all participants send only one state to conclude channel
	// To support later
	// States       | S1  |  S2     | S3  |
	// Participants |  -  |  P1, P2 |  P2 |
	// whoSignedWhat array should be [0, 0, 1]
	// where S2 will be counted as 0
	whoSignedWhat := []uint8{}
	for i := 0; i < len(s.Participants); i++ {
		whoSignedWhat = append(whoSignedWhat, uint8(0))
	}

	params := ConcludeParams{
		OutcomeState:  outcomeState,
		AppData:       appData,
		Signatures:    moveSignatures,
		FixedPart:     fixedPart,
		NumStates:     uint8(1),
		WhoSignedWhat: whoSignedWhat,
	}

	return params, nil
}

// forceMoveSignatures forms signatures as IForceMoveSignature type.
func forceMoveSignatures(s state.State, participantSignatures map[common.Address]state.Signature) []nitro.IForceMoveSignature {
	var moveSignatures []nitro.IForceMoveSignature
	var signatureR, signatureS [32]byte

	for _, a := range s.Participants {
		signature := participantSignatures[a]
		copy(signatureR[:], signature.R)
		copy(signatureS[:], signature.S)

		moveSignatures = append(moveSignatures, nitro.IForceMoveSignature{
			V: signature.V, R: signatureR, S: signatureS,
		})

	}

	return moveSignatures
}
