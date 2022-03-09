package examples

import (
	"app/models/broker"
	"app/pkg/contract"
	"app/pkg/state"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	states "github.com/statechannels/go-nitro/channel/state"
)

var (
	MaxTurnNum              = 5
	initialAmountForBroker1 = big.NewInt(100)
	initialAmountForBroker2 = big.NewInt(200)
)

func Simple(broker1 *broker.Broker,
	broker2 *broker.Broker,
	client contract.Client,
	assetAddress common.Address,
	chainID *big.Int) error {
	// STEP 1 - Initialize pre fund state to sign
	// This should be change later
	// System can increase by 1 on every channel creation
	channelNonce := big.NewInt(0)
	preFundState := state.Build(
		chainID, channelNonce, assetAddress,
		initialAmountForBroker1, initialAmountForBroker2,
		broker1, broker2, false, 0, []byte{},
	)

	// STEP 2 - Open a channel between brokers
	c, err := channel.New(preFundState, broker1.Role)
	if err != nil {
		return err
	}
	fmt.Printf("Channel ID: %+v \n", c.Id)

	// STEP 2 - Sign prefund state by broker 1
	signature1, err := broker1.SignState(&c, preFundState)
	if err != nil {
		return err
	}

	// STEP 3 - Sign prefund state by broker 2
	signature2, err := broker2.SignState(&c, preFundState)
	if err != nil {
		return err
	}

	// STEP 4 - Deposit process. Deposit funds only when prefund state has been completed
	if c.PreFundComplete() {
		// STEP 4.1 - Broker 1 deposit funds to SC
		signerFn := broker1.SignTransaction
		transactionOpts1 := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: broker1.Address, Signer: signerFn, Value: initialAmountForBroker1}
		transaction1, err := client.Contract.Deposit(&transactionOpts1, assetAddress, c.Id, big.NewInt(0), initialAmountForBroker1)
		if err != nil {
			fmt.Printf("Deposit 1: %+v\n", err)
		} else {
			fmt.Printf("Deposit 1:  %+v\n", transaction1)
		}

		// STEP 4.2 - Broker 2 deposit funds to SC
		// Expected held should be the same as broker2's initial amount of funds
		// If there is no funds, SC'll revert transaction
		signerFn = broker2.SignTransaction
		transactionOpts2 := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: broker2.Address, Signer: signerFn, Value: initialAmountForBroker2}
		transaction2, err := client.Contract.Deposit(&transactionOpts2, assetAddress, c.Id, initialAmountForBroker1, initialAmountForBroker2)
		if err != nil {
			fmt.Printf("Deposit 2: %+v\n", err)
		} else {
			fmt.Printf("Deposit 2:  %+v\n", transaction2)
		}
	}

	// STEP 5 - Sign post fund state by broker 1
	signature1, err = broker1.SignState(&c, c.PostFundState())
	if err != nil {
		return err
	}

	// STEP 6 - Sign post fund state by broker 2
	signature2, err = broker2.SignState(&c, c.PostFundState())
	if err != nil {
		return err
	}

	str, err := client.Contract.UnpackStatus(nil, c.Id)
	if err != nil {
		return err
	}
	fmt.Printf("Contract On-Chain Status: %+v\n", str)

	// STEP 7 - Interact between participants off-chain
	if c.PostFundComplete() {
		initialTurnNum := len(c.SignedStateForTurnNum)
		for i := initialTurnNum; i < MaxTurnNum; i++ {
			fmt.Printf("\n\nTurn Number: %d\n", i)

			// Build new state
			newState := state.Build(
				chainID, channelNonce, assetAddress,
				initialAmountForBroker1, initialAmountForBroker2,
				broker1, broker2, false, uint64(i), []byte{},
			)

			c.SignedStateForTurnNum[uint64(i)] = states.NewSignedState(newState)

			// Sign state by broker1
			signature1, err = broker1.SignState(&c, newState)
			if err != nil {
				return err
			}

			fmt.Printf("State completion after 1 participant: %v\n", c.CurrentStateComplete(newState))
			fmt.Printf("State signed by broker1: %v\n", c.CurrentStateSignedByMe(newState))

			// Sign state by broker2
			signature2, err = broker2.SignState(&c, newState)
			if err != nil {
				return err
			}

			fmt.Printf("State completion after 2 participant: %v\n", c.CurrentStateComplete(newState))
		}
	}

	// STEP 8 - Finalize a channel
	finalTurnNum := uint64(MaxTurnNum)

	finalState := state.Build(
		chainID, channelNonce, assetAddress,
		initialAmountForBroker1, initialAmountForBroker2,
		broker1, broker2, true, finalTurnNum, []byte{},
	)

	c.SignedStateForTurnNum[finalTurnNum] = states.NewSignedState(finalState)

	signature1, err = broker1.SignState(&c, finalState)
	if err != nil {
		return err
	}

	signature2, err = broker2.SignState(&c, finalState)
	if err != nil {
		return err
	}
	fmt.Printf("\n\n Channel signatures: %+v\n\n", c.SignedStateForTurnNum)

	concludeParams, err := state.BuildConcludeParams(finalState, signature1, signature2)
	if err != nil {
		return err
	}

	signerFn := broker1.SignTransaction
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: broker1.Address, Signer: signerFn}
	concludeTransaction, err := client.Contract.ConcludeAndTransferAllAssets(&transactionOpts,
		big.NewInt(int64(finalTurnNum)),
		concludeParams.FixedPart,
		concludeParams.AppPart,
		concludeParams.OutcomeState,
		concludeParams.NumStates,
		concludeParams.WhoSignedWhat,
		concludeParams.Signatures,
	)

	if err != nil {
		return err
	}

	fmt.Printf("Conclude transaction: %+v\n\n\n", concludeTransaction)

	return nil
}
