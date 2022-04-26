package init_proposal

import (
	"context"

	"app/internal/proto"
)

type InitProposalService struct {
	proto.UnimplementedInitProposalServiceServer
}

func NewInitProposalService() *InitProposalService {
	svc := InitProposalService{}

	// register("/InitProposal/Create", []string{})
	// register("/InitProposal/AddParticipant", []string{})

	return &svc
}

func (svc *InitProposalService) Create(ctx context.Context, req *proto.CreateProposalRequest) *proto.CreateProposalResponse {
	// res := protocol.NewInitProposal(req.Participant, req.Contract)

	return &proto.CreateProposalResponse{
		InitialProposal: &proto.InitialProposal{},
	}
}

func (svc *InitProposalService) AddParticipant(ctx context.Context, req *proto.AddParticipantRequest) {

}
