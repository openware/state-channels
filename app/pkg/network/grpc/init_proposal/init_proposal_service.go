package init_proposal

import (
	"context"
	"errors"
	"math/big"

	"app/internal/proto"
	"app/pkg/nitro"
	"app/pkg/protocol"

	"github.com/ethereum/go-ethereum/common"
	st "github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

var (
	ErrEmptyParticipant = errors.New("grpc initial proposal: empty participant")
	ErrEmptyContract    = errors.New("grpc initial proposal: empty contract")
)

type InitProposalService struct {
	proto.UnimplementedInitProposalServiceServer
}

func NewInitProposalService() *InitProposalService {
	svc := InitProposalService{}

	return &svc
}

func fromProtoParticipant(participant *proto.Participant) (*protocol.Participant, error) {
	if participant == nil {
		return &protocol.Participant{}, ErrEmptyParticipant
	}

	p := protocol.Participant{
		Address:      common.HexToAddress(participant.Address),
		Destination:  types.Destination(common.HexToHash(participant.Destination)),
		LockedAmount: big.NewInt(participant.LockedAmount),
		Index:        uint(participant.Index),
	}

	return &p, nil
}

func fromProtoContract(contract *proto.Contract) (*protocol.Contract, error) {
	if contract == nil {
		return &protocol.Contract{}, ErrEmptyContract
	}

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

func fromProtoState(state *proto.State) (*st.State, error) {
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

	return &s, nil
}

func fromProtoInitialProposal(ip *proto.InitialProposal) (*protocol.InitProposal, error) {
	contract, err := fromProtoContract(ip.Contract)
	if err != nil {
		return &protocol.InitProposal{}, err
	}

	var participants []*protocol.Participant
	for _, p := range ip.Participants {
		protocolParticipant, err := fromProtoParticipant(p)
		if err != nil {
			return &protocol.InitProposal{}, err
		}

		participants = append(participants, protocolParticipant)
	}

	state, err := fromProtoState(ip.State)
	if err != nil {
		return &protocol.InitProposal{}, err
	}

	proposal := protocol.InitProposal{
		ChannelNonce: big.NewInt(ip.ChannelNonce),
		Contract:     contract,
		Participants: participants,
		State:        state,
	}

	return &proposal, nil
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

func toProtoParticipant(participant *protocol.Participant) (*proto.Participant, error) {
	p := proto.Participant{
		Address:      participant.Address.String(),
		LockedAmount: participant.LockedAmount.Int64(),
		Destination:  participant.Destination.String(),
		Index:        uint64(participant.Index),
	}

	return &p, nil
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

func toProtoInitialProposal(ip *protocol.InitProposal) (*proto.InitialProposal, error) {
	contract, err := toProtoContract(ip.Contract)
	if err != nil {
		return &proto.InitialProposal{}, err
	}

	participants := []*proto.Participant{}
	for _, p := range ip.Participants {
		protoParticipant, err := toProtoParticipant(p)
		if err != nil {
			return &proto.InitialProposal{}, err
		}

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

func (svc *InitProposalService) Create(ctx context.Context, req *proto.CreateProposalRequest) (*proto.CreateProposalResponse, error) {
	participant, err := fromProtoParticipant(req.Participant)
	if err != nil {
		return &proto.CreateProposalResponse{}, err
	}

	contract, err := fromProtoContract(req.Contract)
	if err != nil {
		return &proto.CreateProposalResponse{}, err
	}

	proposal := protocol.NewInitProposal(participant, contract)

	initProposal, err := toProtoInitialProposal(proposal)
	if err != nil {
		return &proto.CreateProposalResponse{}, err
	}

	return &proto.CreateProposalResponse{
		InitialProposal: initProposal,
	}, nil
}

func (svc *InitProposalService) AddParticipant(ctx context.Context, req *proto.AddParticipantRequest) error {
	participant, err := fromProtoParticipant(req.Participant)
	if err != nil {
		return err
	}

	initialProposal, err := fromProtoInitialProposal(req.InitialProposal)
	if err != nil {
		return err
	}

	initialProposal.AddParticipant(participant)

	return nil
}
