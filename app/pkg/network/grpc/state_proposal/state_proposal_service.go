package state_proposal

import (
	"app/internal/proto"
	"app/pkg/protocol"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	st "github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

type StateProposalService struct {
	proto.UnimplementedStateProposalServiceServer
}

func NewStateProposalService() *StateProposalService {
	svc := StateProposalService{}

	return &svc
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

func fromProtoStateProposal(stateProposal *proto.StateProposal) (*protocol.StateProposal, error) {
	sp := protocol.StateProposal{}

	return &sp, nil
}

func (svc *StateProposalService) Create(ctx context.Context, req proto.CreateRequest) (proto.CreateResponse, error) {
	state, err := fromProtoState(req.State)
	if err != nil {
		return proto.CreateResponse{}, nil
	}

	stateProposal, err := protocol.NewStateProposal(state)
	if err != nil {
		return proto.CreateResponse{}, nil
	}

	protoStateProposal, err := toProtoStateProposal(stateProposal)
	if err != nil {
		return proto.CreateResponse{}, nil
	}

	return proto.CreateResponse{StateProposal: protoStateProposal}, nil
}

func (svc *StateProposalService) SetFinal(ctx context.Context, req proto.StateProposalRequest) error {
	stateProposal, err := fromProtoStateProposal(req.StateProposal)
	if err != nil {
		return err
	}

	stateProposal.SetFinal()
	return nil
}

func (svc *StateProposalService) PendingLiability(ctx context.Context, req proto.LiabilityRequest) error {
	stateProposal, err := fromProtoStateProposal(req.StateProposal)
	if err != nil {
		return err
	}

	return nil
}

func (svc *StateProposalService) ExecutedLiability(ctx context.Context, req proto.LiabilityRequest) error {
	return nil
}

func (svc *StateProposalService) RevertLiability(ctx context.Context, req proto.LiabilityRequest) error {
	return nil
}

func (svc *StateProposalService) ApproveLiabilities(ctx context.Context, req proto.StateProposalRequest) error {
	return nil
}
