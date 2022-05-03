package grpc

import (
	"context"

	"app/internal/proto"
	"app/pkg/nitro"
	"app/pkg/protocol"

	"github.com/ethereum/go-ethereum/common"
)

type InitProposalService struct {
	proto.UnimplementedInitProposalServiceServer
}

func NewInitProposalService() *InitProposalService {
	svc := InitProposalService{}

	return &svc
}

func (svc *InitProposalService) Create(ctx context.Context, req *proto.CreateProposalRequest) (*proto.CreateProposalResponse, error) {
	participant := fromProtoParticipant(req.Participant)
	nitroClient, err := nitro.NewClient(req.ContractAddress, req.RpcUrl)
	if err != nil {
		return &proto.CreateProposalResponse{}, err
	}

	contract := protocol.NewContract(nitroClient, common.HexToAddress(req.AssetAddress))
	proposal := protocol.NewInitProposal(participant, contract)

	initProposal, err := toProtoInitialProposal(proposal)
	if err != nil {
		return &proto.CreateProposalResponse{}, err
	}

	return &proto.CreateProposalResponse{
		InitialProposal: initProposal,
	}, nil
}

// bytes as initial proposal
func (svc *InitProposalService) AddParticipant(ctx context.Context, req *proto.AddParticipantRequest) error {
	participant := fromProtoParticipant(req.Participant)

	initialProposal, err := fromProtoInitialProposal(req.InitialProposal)
	if err != nil {
		return err
	}

	initialProposal.AddParticipant(participant)

	return nil
}
