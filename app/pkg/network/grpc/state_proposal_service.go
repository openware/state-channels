package grpc

import (
	"app/internal/liability"
	"app/internal/proto"
	"context"

	"github.com/shopspring/decimal"
)

type StateProposalService struct {
	proto.UnimplementedStateProposalServiceServer
}

func NewStateProposalService() *StateProposalService {
	svc := StateProposalService{}

	return &svc
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

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return err
	}

	stateProposal.PendingLiability(uint(req.From), uint(req.To), liability.Asset(req.Asset), amount)

	return nil
}

func (svc *StateProposalService) ExecutedLiability(ctx context.Context, req proto.LiabilityRequest) error {
	stateProposal, err := fromProtoStateProposal(req.StateProposal)
	if err != nil {
		return err
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return err
	}

	return stateProposal.ExecutedLiability(uint(req.From), uint(req.To), liability.Asset(req.Asset), amount)
}

func (svc *StateProposalService) RevertLiability(ctx context.Context, req proto.LiabilityRequest) error {
	stateProposal, err := fromProtoStateProposal(req.StateProposal)
	if err != nil {
		return err
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return err
	}

	return stateProposal.RevertLiability(uint(req.From), uint(req.To), liability.Asset(req.Asset), amount)
}

func (svc *StateProposalService) ApproveLiabilities(ctx context.Context, req proto.StateProposalRequest) error {
	stateProposal, err := fromProtoStateProposal(req.StateProposal)
	if err != nil {
		return err
	}

	return stateProposal.ApproveLiabilities()
}
