package examples

import (
	"app/models/broker"
	"app/pkg/contract"
	"app/pkg/liability"
	"app/pkg/state"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/statechannels/go-nitro/channel"
	states "github.com/statechannels/go-nitro/channel/state"
)

var (
	GasLimit = uint64(6721975)
	GasPrice = big.NewInt(20000000000)
)

func SimpleTrade(broker1 *broker.Broker, broker2 *broker.Broker, client contract.Client, assetAddress common.Address, chainID *big.Int) error {
	channelNonce := big.NewInt(0)
	initialAmountForBroker1 := big.NewInt(10000)
	initialAmountForBroker2 := big.NewInt(20000)

	preFundState := state.Build(
		chainID, channelNonce, assetAddress,
		initialAmountForBroker1, initialAmountForBroker2,
		broker1, broker2, false, 0, []byte{},
	)

	c, err := channel.New(preFundState, broker1.Role)
	if err != nil {
		return err
	}

	_, err = broker1.SignState(&c, preFundState)
	if err != nil {
		return err
	}

	_, err = broker2.SignState(&c, preFundState)
	if err != nil {
		return err
	}

	if c.PreFundComplete() {
		signerFn := broker1.SignTransaction
		transactionOpts1 := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: broker1.Address, Signer: signerFn, Value: initialAmountForBroker1}
		_, err = client.Contract.Deposit(&transactionOpts1, assetAddress, c.Id, big.NewInt(0), initialAmountForBroker1)
		if err != nil {
			return err
		}

		signerFn = broker2.SignTransaction
		transactionOpts2 := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: broker2.Address, Signer: signerFn, Value: initialAmountForBroker2}
		_, err = client.Contract.Deposit(&transactionOpts2, assetAddress, c.Id, initialAmountForBroker1, initialAmountForBroker2)
		if err != nil {
			return err
		}
	}

	_, err = broker1.SignState(&c, c.PostFundState())
	if err != nil {
		return err
	}

	_, err = broker2.SignState(&c, c.PostFundState())
	if err != nil {
		return err
	}

	appData1 := liability.Liability{
		From:            broker1,
		To:              broker2,
		Type:            liability.REQ,
		ToCurrencyArray: []string{"ETH", "BTC"},
		ToAmountArray:   []decimal.Decimal{decimal.NewFromInt(12), decimal.NewFromInt(-2)},
	}

	appDataBytes1, err := appData1.EncodeToBytes()
	if err != nil {
		return err
	}

	state2 := state.Build(
		chainID, channelNonce, assetAddress,
		initialAmountForBroker1, initialAmountForBroker2,
		broker1, broker2, false, 2, appDataBytes1,
	)

	fmt.Printf("State 2: %v\n\n", state2)

	_, err = broker1.SignState(&c, state2)
	if err != nil {
		return err
	}

	_, err = broker2.SignState(&c, state2)
	if err != nil {
		return err
	}

	appData2 := liability.Liability{
		From:            broker1,
		To:              broker2,
		Type:            liability.ACK,
		ToCurrencyArray: []string{"ETH", "BTC"},
		ToAmountArray:   []decimal.Decimal{decimal.NewFromInt(12), decimal.NewFromInt(-2)},
	}

	appDataBytes2, err := appData2.EncodeToBytes()
	if err != nil {
		return err
	}

	state3 := state.Build(
		chainID, channelNonce, assetAddress,
		initialAmountForBroker1, initialAmountForBroker2,
		broker1, broker2, false, 3, appDataBytes2,
	)

	fmt.Printf("State 2: %v\n\n", state3)

	_, err = broker1.SignState(&c, state3)
	if err != nil {
		return err
	}

	_, err = broker2.SignState(&c, state3)
	if err != nil {
		return err
	}

	// In this example we need to conclude a channel each time, by default state channels will be opened as long as possible
	// STEP 8 - Finalize a channel
	finalTurnNum := uint64(4)

	finalState := state.Build(
		chainID, channelNonce, assetAddress,
		initialAmountForBroker1, initialAmountForBroker2,
		broker1, broker2, true, finalTurnNum, []byte{},
	)

	c.SignedStateForTurnNum[finalTurnNum] = states.NewSignedState(finalState)

	signature1, err := broker1.SignState(&c, finalState)
	if err != nil {
		return err
	}

	signature2, err := broker2.SignState(&c, finalState)
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
