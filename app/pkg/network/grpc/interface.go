package grpc

import (
	"app/internal/liability"
	"app/internal/proto"
	"app/pkg/nitro"
	"app/pkg/protocol"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel/state"
	st "github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

// TODO
// return errors where needed

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

func fromProtoContract(contract *proto.Contract) (*protocol.Contract, error) {
	client, err := nitro.DecodeFromBytes(contract.NitroClient)
	if err != nil {
		return &protocol.Contract{}, err
	}

	c := protocol.Contract{
		Client:       client,
		AssetAddress: common.HexToAddress(contract.AssetAddress),
	}

	return &c, nil
}

func toProtoContract(contract *protocol.Contract) (*proto.Contract, error) {
	client, err := contract.Client.EncodeToBytes()
	if err != nil {
		return &proto.Contract{}, nil
	}

	c := proto.Contract{
		AssetAddress: contract.AssetAddress.String(),
		NitroClient:  client,
	}

	return &c, nil
}

func fromProtoStateProposal(stateProposal *proto.StateProposal) (*protocol.StateProposal, error) {
	state := fromProtoState(stateProposal.State)
	liabilityStateBytes := stateProposal.LiabilityState
	liabilityState, err := liability.DecodeFromBytes(liabilityStateBytes)
	if err != nil {
		return &protocol.StateProposal{}, nil
	}

	sp := protocol.StateProposal{}
	sp.SetState(state)
	sp.SetLiabilitiesState(liabilityState)

	return &sp, nil
}

func toProtoStateProposal(stateProposal *protocol.StateProposal) (*proto.StateProposal, error) {
	state := stateProposal.State()
	protoState, err := toProtoState(&state)
	if err != nil {
		return &proto.StateProposal{}, nil
	}

	liabilityState, err := stateProposal.LiabilityState().EncodeToBytes()
	if err != nil {
		return &proto.StateProposal{}, nil
	}

	sp := proto.StateProposal{
		State:          protoState,
		LiabilityState: liabilityState,
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

func fromProtoInitialProposal(ip *proto.InitialProposal) (*protocol.InitProposal, error) {
	contract, err := fromProtoContract(ip.Contract)
	if err != nil {
		return &protocol.InitProposal{}, err
	}

	var participants []*protocol.Participant
	for _, p := range ip.Participants {
		protocolParticipant := fromProtoParticipant(p)
		participants = append(participants, protocolParticipant)
	}

	state := fromProtoState(ip.State)

	proposal := protocol.InitProposal{
		ChannelNonce: big.NewInt(ip.ChannelNonce),
		Contract:     contract,
		Participants: participants,
		State:        state,
	}

	return &proposal, nil
}

func toProtoInitialProposal(ip *protocol.InitProposal) (*proto.InitialProposal, error) {
	contract, err := toProtoContract(ip.Contract)
	if err != nil {
		return &proto.InitialProposal{}, err
	}

	participants := []*proto.Participant{}
	for _, p := range ip.Participants {
		protoParticipant := toProtoParticipant(p)
		participants = append(participants, protoParticipant)
	}

	state, err := toProtoState(ip.State)
	if err != nil {
		return &proto.InitialProposal{}, err
	}

	prop := proto.InitialProposal{
		ChannelNonce: ip.ChannelNonce.Int64(),
		Contract:     contract,
		Participants: participants,
		State:        state,
	}

	return &prop, nil
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

func fromProtoChannel(channel *proto.Channel) (*protocol.Channel, error) {
	ch := protocol.Channel{}

	return &ch, nil
}

func toProtoChannel(channel *protocol.Channel) (*proto.Channel, error) {
	ch := proto.Channel{}

	return &ch, nil
}

func fromProtoSignatures(signatures map[string]*proto.Signature) map[common.Address]state.Signature {
	sigs := make(map[common.Address]state.Signature)

	for hex, sig := range signatures {
		sigs[common.HexToAddress(hex)] = fromProtoSignature(sig)
	}

	return sigs
}
