package grpc

import (
	"app/internal/proto"
	"app/pkg/protocol"
	"context"
)

type ChannelService struct {
	proto.UnimplementedChannelServiceServer
}

func NewChannelService() *ChannelService {
	svc := ChannelService{}

	return &svc
}

func Init(ctx context.Context, req proto.InitChannelRequest) (*proto.InitChannelResponse, error) {
	initialProposal, err := fromProtoInitialProposal(req.InitialProposal)
	if err != nil {
		return &proto.InitChannelResponse{}, nil
	}

	channel, err := protocol.InitChannel(initialProposal, uint(req.ParticipantIndex))
	if err != nil {
		return &proto.InitChannelResponse{}, nil
	}

	protoChannel, err := toProtoChannel(channel)
	if err != nil {
		return &proto.InitChannelResponse{}, nil
	}

	return &proto.InitChannelResponse{Channel: protoChannel}, nil
}

func ApproveInit(ctx context.Context, req proto.ApproveRequest) (*proto.ApproveResponse, error) {
	channel, err := fromProtoChannel(req.Channel)
	if err != nil {
		return &proto.ApproveResponse{}, nil
	}

	signature, err := channel.ApproveInitChannel(req.PrivateKey)
	if err != nil {
		return &proto.ApproveResponse{}, nil
	}

	protoSignature := toProtoSignature(signature)

	return &proto.ApproveResponse{Signature: protoSignature}, nil
}

func Fund(ctx context.Context, req proto.FundChannelRequest) (*proto.FundChannelResponse, error) {
	channel, err := fromProtoChannel(req.Channel)
	if err != nil {
		return &proto.FundChannelResponse{}, nil
	}

	participant := fromProtoParticipant(req.Participant)

	// TODO
	// gas station

	transaction, err := channel.FundChannel(participant, req.PrivateKey)
	if err != nil {
		return &proto.FundChannelResponse{}, err
	}

	return &proto.FundChannelResponse{TxId: transaction.Hash().String()}, nil
}

func ApproveFunding(ctx context.Context, req proto.ApproveRequest) (*proto.ApproveResponse, error) {
	channel, err := fromProtoChannel(req.Channel)
	if err != nil {
		return &proto.ApproveResponse{}, nil
	}

	signature, err := channel.ApproveChannelFunding(req.PrivateKey)
	protoSignature := toProtoSignature(signature)

	return &proto.ApproveResponse{Signature: protoSignature}, nil
}

func ProposeState(ctx context.Context, req proto.ChannelRequest) (*proto.ProposeResponse, error) {
	channel, err := fromProtoChannel(req.Channel)
	if err != nil {
		return &proto.ProposeResponse{}, nil
	}

	stateProposal, err := channel.ProposeState()
	if err != nil {
		return &proto.ProposeResponse{}, nil
	}

	protoStateProposal, err := toProtoStateProposal(stateProposal)
	if err != nil {
		return &proto.ProposeResponse{}, nil
	}

	return &proto.ProposeResponse{StateProposal: protoStateProposal}, nil
}

func SignState(ctx context.Context, req proto.SignStateRequest) (*proto.SignStateResponse, error) {
	channel, err := fromProtoChannel(req.Channel)
	if err != nil {
		return &proto.SignStateResponse{}, nil
	}

	stateProposal, err := fromProtoStateProposal(req.StateProposal)
	if err != nil {
		return &proto.SignStateResponse{}, err
	}

	signature, err := channel.SignState(stateProposal, req.PrivateKey)
	if err != nil {
		return &proto.SignStateResponse{}, err
	}

	protoSignature := toProtoSignature(signature)

	return &proto.SignStateResponse{Signature: protoSignature}, nil
}

func Conclude(ctx context.Context, req proto.ConcludeRequest) (*proto.ConcludeResponse, error) {
	channel, err := fromProtoChannel(req.Channel)
	if err != nil {
		return &proto.ConcludeResponse{}, nil
	}

	participant := fromProtoParticipant(req.Participant)
	signatures := fromProtoSignatures(req.Signatures)

	transaction, err := channel.Conclude(participant, req.PrivateKey, signatures)
	if err != nil {
		return &proto.ConcludeResponse{}, err
	}

	return &proto.ConcludeResponse{TxId: transaction.Hash().String()}, nil
}

func CheckSignature(ctx context.Context, req proto.CheckSignatureRequest) (*proto.BoolResponse, error) {
	channel, err := fromProtoChannel(req.Channel)
	if err != nil {
		return &proto.BoolResponse{}, nil
	}

	signature := fromProtoSignature(req.Signature)
	state := fromProtoState(req.State)

	ok, err := channel.CheckSignature(signature, state)
	if err != nil {
		return &proto.BoolResponse{}, nil
	}

	return &proto.BoolResponse{Ok: ok}, nil
}

func CurrentState(ctx context.Context, req proto.ChannelRequest) (*proto.CurrentStateResponse, error) {
	channel, err := fromProtoChannel(req.Channel)
	if err != nil {
		return &proto.CurrentStateResponse{}, nil
	}

	state := channel.CurrentState()

	protoState, err := toProtoState(&state)
	if err != nil {
		return &proto.CurrentStateResponse{}, nil
	}

	return &proto.CurrentStateResponse{State: protoState}, nil
}

func CheckHoldings(ctx context.Context, req proto.ChannelRequest) (*proto.CheckHoldingsResponse, error) {
	channel, err := fromProtoChannel(req.Channel)
	if err != nil {
		return &proto.CheckHoldingsResponse{}, nil
	}

	amount, err := channel.CheckHoldings()
	if err != nil {
		return &proto.CheckHoldingsResponse{}, nil
	}

	return &proto.CheckHoldingsResponse{Amount: amount.Int64()}, nil
}

func StateIsFinal(ctx context.Context, req proto.ChannelRequest) (*proto.BoolResponse, error) {
	channel, err := fromProtoChannel(req.Channel)
	if err != nil {
		return &proto.BoolResponse{}, nil
	}

	isFinal := channel.StateIsFinal()

	return &proto.BoolResponse{Ok: isFinal}, nil
}
