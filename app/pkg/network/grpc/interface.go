package grpc

import (
	"app/internal/proto"
	"app/pkg/protocol"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	st "github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

func fromProtoParticipant(participant *proto.Participant) *protocol.Participant {
	return &protocol.Participant{
		Address:      common.HexToAddress(participant.Address),
		Destination:  types.Destination(common.HexToHash(participant.Destination)),
		LockedAmount: big.NewInt(participant.LockedAmount),
		Index:        uint(participant.Index),
	}
}

func toProtoParticipant(participant *protocol.Participant) *proto.Participant {
	return &proto.Participant{
		Address:      participant.Address.String(),
		Destination:  participant.Destination.String(),
		LockedAmount: participant.LockedAmount.Int64(),
		Index:        uint64(participant.Index),
	}
}

func fromProtoStateProposal(data []byte) (*protocol.StateProposal, error) {
	sp, err := protocol.DecodeStateProposalFromBytes(data)
	if err != nil {
		return &protocol.StateProposal{}, nil
	}

	return &sp, nil
}

func fromProtoState(state *proto.State) *st.State {
	participants := []types.Address{}
	for _, p := range state.Participants {
		participants = append(participants, common.HexToAddress(p))
	}

	s := st.State{
		ChainId:           big.NewInt(int64(state.ChainId)),
		ChannelNonce:      big.NewInt(state.ChannelNonce),
		Participants:      participants,
		IsFinal:           state.IsFinal,
		TurnNum:           state.TurnNum,
		ChallengeDuration: big.NewInt(int64(state.ChallengeDuration)),
		AppData:           state.AppData,
		AppDefinition:     common.HexToAddress(state.AppDefinition),
	}

	return &s
}

func toProtoState(state *st.State) (*proto.State, error) {
	outcome, err := state.Outcome.Encode()
	if err != nil {
		return &proto.State{}, err
	}

	participants := []string{}
	for _, p := range state.Participants {
		participants = append(participants, p.String())
	}

	s := proto.State{
		ChainId:           state.ChainId.Uint64(),
		ChannelNonce:      state.ChannelNonce.Int64(),
		AppDefinition:     state.AppDefinition.String(),
		ChallengeDuration: state.ChallengeDuration.Uint64(),
		AppData:           state.AppData,
		IsFinal:           state.IsFinal,
		TurnNum:           state.TurnNum,
		Outcome:           outcome,
		Participants:      participants,
	}

	return &s, nil
}

func fromProtoInitialProposal(data []byte) (*protocol.InitProposal, error) {
	ip, err := protocol.DecodeInitProposalFromBytes(data)
	if err != nil {
		return &protocol.InitProposal{}, err
	}

	return &ip, nil
}

func toProtoSignature(signature state.Signature) *proto.Signature {
	sig := proto.Signature{
		R: signature.R,
		S: signature.S,
		V: []byte{signature.V},
	}

	return &sig
}

func fromProtoSignature(signature *proto.Signature) state.Signature {
	sig := state.Signature{
		R: signature.R,
		S: signature.S,
		V: signature.V[0],
	}

	return sig
}

func fromProtoChannel(data []byte) (*protocol.Channel, error) {
	ch, err := protocol.DecodeChannelFromBytes(data)
	if err != nil {
		return &protocol.Channel{}, err
	}

	return &ch, nil
}

func fromProtoSignatures(signatures map[string]*proto.Signature) map[common.Address]state.Signature {
	sigs := make(map[common.Address]state.Signature)

	for hex, sig := range signatures {
		sigs[common.HexToAddress(hex)] = fromProtoSignature(sig)
	}

	return sigs
}
