package examples

import (
	"app/pkg/protocol"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/shopspring/decimal"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/crypto"
)

func Demo(participants []protocol.Participant, privKeys map[protocol.Participant][]byte, contract *protocol.Contract) error {
	estimatedGasPrice, err := protocol.CalculateGasPrice(contract.Client.Eth)
	if err != nil {
		return err
	}

	gasStation := protocol.GasStation{GasPrice: estimatedGasPrice}

	ch, err := initChannel(participants, privKeys, contract)
	if err != nil {
		return nil
	}

	err = fundChannel(ch, participants, privKeys, gasStation)
	if err != nil {
		return nil
	}

	err = confirmChannelFund(ch, privKeys)
	if err != nil {
		return nil
	}

	err = proposeState(ch, participants, privKeys)
	if err != nil {
		return nil
	}

	err = concludeChannel(ch, participants, privKeys, gasStation)
	if err != nil {
		return nil
	}

	return nil
}

func initChannel(
	participants []protocol.Participant,
	privKeys map[protocol.Participant][]byte,
	contract *protocol.Contract) (*protocol.Channel, error) {

	prop, err := initialProposal(participants, contract)
	if err != nil {
		return &protocol.Channel{}, err
	}

	err = confirmPrompt("Initialize channel")
	if err != nil {
		return &protocol.Channel{}, err
	}

	fmt.Println("\nChannel initialization...")
	ch, err := protocol.InitChannel(*prop, participants[0].Index)
	if err != nil {
		return &protocol.Channel{}, err
	}

	fmt.Println("Channel initialized")

	for p, pKey := range privKeys {
		_, err := ch.ApproveInitChannel(pKey)
		if err != nil {
			return &protocol.Channel{}, err
		}
		fmt.Printf("Sign PreFund state by participant with index [%d]\n", p.Index)

	}

	fmt.Println("PreFund state has been completed\n")

	return ch, nil
}

func initialProposal(
	participants []protocol.Participant,
	contract *protocol.Contract) (*protocol.InitProposal, error) {

	prop := protocol.NewInitProposal(participants[0], *contract)
	for _, p := range participants[1:] {
		prop.AddParticipant(p)
	}

	fmt.Printf("%v\n\n", color.HiWhiteString("Initial data:"))
	fmt.Printf("%v\n", color.HiYellowString("Participants:"))

	for i, p := range prop.Participants {
		fmt.Printf("%v\n", color.GreenString("Participant index [%d], address [%s]", i, p.Address))
	}

	fmt.Printf("\n%v\n", color.HiYellowString("Initial State (PreFund state):"))
	fmt.Println(color.GreenString("Chain ID: %v", prop.State.ChainId))
	fmt.Println(color.GreenString("Channel Nonce: %v", time.UnixMilli(prop.State.ChannelNonce.Int64())))
	fmt.Println(color.GreenString("App Data: %v", prop.State.AppData))
	fmt.Println(color.GreenString("Asset: %v", prop.State.Outcome[0].Asset))
	fmt.Println(color.GreenString("Outcome Allocations: "))

	for _, a := range prop.State.Outcome[0].Allocations {
		destination, err := a.Destination.ToAddress()
		if err != nil {
			return &protocol.InitProposal{}, err
		}
		fmt.Println(color.GreenString("Destination Address [%v] Amount [%v]", destination, a.Amount))
	}
	fmt.Println(color.GreenString("Turn Number: %v", prop.State.TurnNum))
	fmt.Println(color.GreenString("Is Final: %v\n", prop.State.IsFinal))

	return prop, nil
}

func fundChannel(
	ch *protocol.Channel,
	participants []protocol.Participant,
	privKeys map[protocol.Participant][]byte,
	gasStation protocol.GasStation) error {

	err := confirmPrompt("Fund channel")
	if err != nil {
		return err
	}
	fmt.Println()

	for _, p := range participants {
		transaction, err := ch.FundChannel(p, privKeys[p], gasStation)
		if err != nil {
			return err
		}
		fmt.Printf("Funding channel by participant [%d] with amount [%d], transaction hash [%s] \n", p.Index, p.LockedAmount, transaction.Hash())
		time.Sleep(time.Second * 10)
	}
	fmt.Println("Channel funding has been completed\n")

	return nil
}

func confirmChannelFund(
	ch *protocol.Channel,
	privKeys map[protocol.Participant][]byte) error {

	err := confirmPrompt("Sign PostFund state")
	if err != nil {
		return err
	}
	fmt.Println()

	for p, pKey := range privKeys {
		_, err := ch.ApproveChannelFunding(pKey)
		if err != nil {
			return err
		}

		fmt.Printf("Sign PostFund state by participant [%d]\n", p.Index)
	}
	fmt.Println(color.HiYellowString("\nPostFund state: "))
	fmt.Println(color.GreenString("App Data: %v", ch.LastState.AppData))
	fmt.Println(color.GreenString("Turn Number: %d", ch.LastState.TurnNum))
	fmt.Println(color.GreenString("Is Final:  %v\n", ch.LastState.IsFinal))
	fmt.Println("PostFund state has been completed\n")

	return nil
}

func proposeState(
	ch *protocol.Channel,
	participants []protocol.Participant,
	privKeys map[protocol.Participant][]byte,
) error {

	for {
		err := confirmPrompt("Propose new state")
		if err != nil {
			break
		}

		appData, ok, err := proposeLiability(*ch)
		if err != nil {
			return err
		}

		if len(appData) > 0 && ok {
			st, err := ch.ProposeState()
			if err != nil {
				return err
			}

			st.State.AppData = appData

			liabilityState, err := protocol.DecodeLiabilityFromBytes(st.State.AppData)
			if err != nil {
				return err
			}

			fmt.Println(color.HiYellowString("\nProposed state: "))
			fmt.Println(color.GreenString("App Data:"))
			color.Set(color.FgGreen)
			liabilityState.Print()
			color.Unset()
			fmt.Println(color.GreenString("Turn Number: %d", st.State.TurnNum))
			fmt.Println(color.GreenString("Is Final:  %v\n", st.State.IsFinal))

			for p, pKey := range privKeys {
				_, err := ch.SignState(st, pKey)
				if err != nil {
					return err
				}
				fmt.Printf("Sign proposed state by participant with index [%d]\n", p.Index)
			}
		}
	}

	return nil
}

func proposeLiability(ch protocol.Channel) ([]byte, bool, error) {
	ok := false
	lastState := ch.LastState
	sp, err := protocol.NewStateProposal(&state.State{AppData: lastState.AppData})
	if err != nil {
		return []byte{}, ok, err
	}

	for {
		err := confirmPrompt("Add New Liability")
		if err != nil {
			break
		}

		from, err := inputPrompt("From")
		if err != nil {
			break
		}

		to, err := inputPrompt("To")
		if err != nil {
			break
		}

		req, err := inputPrompt("Type (REQ/ACK/REVERT)")
		if err != nil {
			break
		}

		asset, err := inputPrompt("Asset")
		if err != nil {
			break
		}

		amount, err := inputPrompt("Amount")
		if err != nil {
			break
		}

		fromNumber, err := strconv.ParseUint(from, 10, 64)
		if err != nil {
			break
		}

		toNumber, err := strconv.ParseUint(to, 10, 64)
		if err != nil {
			break
		}

		amountNumber, err := decimal.NewFromString(amount)
		if err != nil {
			break
		}

		if req == "REQ" {
			err = sp.RequestLiability(uint(fromNumber), uint(toNumber), protocol.Asset(asset), amountNumber)
		} else if req == "ACK" {
			err = sp.AcknowledgeLiability(uint(fromNumber), uint(toNumber), protocol.Asset(asset), amountNumber)
		} else if req == "REVERT" {
			err = sp.RevertLiability(uint(fromNumber), uint(toNumber), protocol.Asset(asset), amountNumber)
		} else {
			break
		}

		if err != nil {
			return []byte{}, ok, err
		}

		err = sp.ApproveLiabilities()
		if err != nil {
			return []byte{}, ok, err
		}

		ok = true
	}

	return sp.State.AppData, ok, nil
}

func concludeChannel(
	ch *protocol.Channel,
	participants []protocol.Participant,
	privKeys map[protocol.Participant][]byte,
	gasStation protocol.GasStation) error {

	err := confirmPrompt("Finalize channel")
	if err != nil {
		return err
	}

	fmt.Println("\nChannel finalization...")
	finalState, err := ch.ProposeState()
	if err != nil {
		return err
	}
	finalState.SetFinal()

	liabilityState, err := protocol.DecodeLiabilityFromBytes(finalState.State.AppData)

	fmt.Println(color.HiYellowString("\nFinal state: "))
	if err == nil {
		fmt.Println(color.GreenString("App Data:"))
		color.Set(color.FgGreen)
		liabilityState.Print()
		color.Unset()
	} else {
		fmt.Println(color.GreenString("App Data: %v", ch.LastState.AppData))
	}

	fmt.Println(color.GreenString("Turn Number: %d", ch.LastState.TurnNum))
	fmt.Println(color.GreenString("Is Final:  %v\n", ch.LastState.IsFinal))

	participantSignatures := make(map[common.Address]crypto.Signature)
	for p, pKey := range privKeys {
		signature, err := ch.SignState(finalState, pKey)
		if err != nil {
			return err
		}

		fmt.Printf("Sign final state by participant with index [%d]\n", p.Index)
		participantSignatures[p.Address] = signature
	}

	transaction, err := ch.Conclude(participants[0], privKeys[participants[0]], participantSignatures, gasStation)
	if err != nil {
		return err
	}

	fmt.Printf("\nConclude transaction hash [%s]\n", transaction.Hash())

	return nil
}

func confirmPrompt(label string) error {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		return err
	}

	return nil
}

func inputPrompt(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}

	res, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return res, nil
}
