package examples

import (
	"app/pkg/eth/gasprice"
	"app/pkg/protocol"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/crypto"
)

var MaxTurnNum = uint64(5)

func Simple(participants []protocol.Participant, privKeys map[protocol.Participant][]byte, contract *protocol.Contract) error {
	prop := protocol.NewInitProposal(participants[0], *contract)
	for _, p := range participants[1:] {
		prop.AddParticipant(p)
	}

	ch, err := protocol.InitChannel(prop, participants[0].Index)
	if err != nil {
		return err
	}

	for _, pKey := range privKeys {
		_, err := ch.ApproveInitChannel(pKey)
		if err != nil {
			return err
		}
	}

	estimatedGasPrice, err := gasprice.Calculate(contract.Client.Eth)
	if err != nil {
		return err
	}

	gasStation := gasprice.Station{GasPrice: estimatedGasPrice}

	for _, p := range participants {
		_, err = ch.FundChannel(p, privKeys[p], gasStation)
		if err != nil {
			return err
		}
	}

	for _, pKey := range privKeys {
		_, err := ch.ApproveChannelFunding(pKey)
		if err != nil {
			return err
		}
	}

	initialTurnNum := channel.PostFundTurnNum
	for i := initialTurnNum; i < MaxTurnNum; i++ {
		st, err := ch.ProposeState()
		if err != nil {
			return err
		}

		for _, pKey := range privKeys {
			_, err := ch.SignState(st, pKey)
			if err != nil {
				return err
			}
		}
	}

	finalState, err := ch.ProposeState()
	if err != nil {
		return err
	}
	finalState.SetFinal()

	participantSignatures := make(map[common.Address]crypto.Signature)
	for p, pKey := range privKeys {
		signature, err := ch.SignState(finalState, pKey)
		if err != nil {
			return err
		}

		participantSignatures[p.Address] = signature
	}

	_, err = ch.Conclude(participants[0], privKeys[participants[0]], participantSignatures, gasStation)
	if err != nil {
		return err
	}

	return nil
}
