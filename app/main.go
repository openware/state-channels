package main

import (
	"app/models/broker"
	"app/pkg/contract"
	"app/pkg/parser"
	"app/pkg/state"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	states "github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
)

var Broker1, Broker2 *broker.Broker
var ChainId *big.Int
var ContractAddress string

var (
	RpcUrl                  = "http://127.0.0.1:8545"
	AssetAddress            = common.HexToAddress("0x00")
	MaxTurnNum              = 5
	GasLimit                = uint64(6721975)
	GasPrice                = big.NewInt(20000000000)
	initialAmountForBroker1 = big.NewInt(100)
	initialAmountForBroker2 = big.NewInt(200)
)

// Simple state channel example
func main() {
	// STEP 1 - Initialize participants (brokers), deploy smart-contract
	mydir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vaultAccount, err := parser.FromFileToAccounts(mydir + "/../contracts/accounts.json")
	if err != nil {
		panic(err)
	}

	contractObj, err := parser.FromFileToContract(mydir + "/../contracts/addresses.json")
	if err != nil {
		panic(err)
	}

	ContractAddress = contractObj.ChainIds[0].SC.NitroAdj.Address

	// Initialize Broker 1
	broker1 := vaultAccount.Accounts[0]
	privateKey1 := strings.TrimPrefix(broker1.PrivateKey, "0x")
	Broker1 = broker.New(common.HexToAddress(broker1.Address), types.Destination(common.HexToHash(broker1.Address)), common.Hex2Bytes(privateKey1), 0)

	// Initialize Broker 2
	broker2 := vaultAccount.Accounts[1]
	privateKey2 := strings.TrimPrefix(broker2.PrivateKey, "0x")
	Broker2 = broker.New(common.HexToAddress(broker2.Address), types.Destination(common.HexToHash(broker2.Address)), common.Hex2Bytes(privateKey2), 1)

	// STEP 2 - Initialize SC client
	client, err := contract.NewClient(ContractAddress, RpcUrl)
	if err != nil {
		panic(err)
	}

	// STEP 3 - Get Chain ID from the client
	ChainId = client.ChainId
	fmt.Printf("Chain Id: %v\n", ChainId)

	// STEP 4 - Initialize pre fund state to sign
	preFundState := state.Build(
		ChainId, AssetAddress,
		initialAmountForBroker1, initialAmountForBroker2,
		Broker1, Broker2, false, 0,
	)

	// STEP 5 - Open a channel between brokers
	c, err := channel.New(preFundState, Broker1.Role)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Channel ID: %+v \n", c.Id)

	// STEP 6 - Sign prefund state by broker 1
	signature1, err := state.Sign(&c, preFundState, Broker1)
	if err != nil {
		panic(err)
	}

	// STEP 7 - Sign prefund state by broker 2
	signature2, err := state.Sign(&c, preFundState, Broker2)
	if err != nil {
		panic(err)
	}

	// STEP 8 - Deposit process. Deposit funds only when prefund state has been completed
	if c.PreFundComplete() {
		// STEP 8.1 - Broker 1 deposit funds to SC
		signerFn := Broker1.SignTransaction
		transactionOpts1 := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: Broker1.Address, Signer: signerFn, Value: initialAmountForBroker1}
		transaction1, err := client.Contract.Deposit(&transactionOpts1, AssetAddress, c.Id, big.NewInt(0), initialAmountForBroker1)
		if err != nil {
			fmt.Printf("Deposit 1: %+v\n", err)
		} else {
			fmt.Printf("Deposit 1:  %+v\n", transaction1)
		}

		// STEP 8.2 - Broker 2 deposit funds to SC
		// Expected held should be the same as Broker2's initial amount of funds
		// If there is no funds, SC'll revert transaction
		signerFn = Broker2.SignTransaction
		transactionOpts2 := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: Broker2.Address, Signer: signerFn, Value: initialAmountForBroker2}
		transaction2, err := client.Contract.Deposit(&transactionOpts2, AssetAddress, c.Id, initialAmountForBroker1, initialAmountForBroker2)
		if err != nil {
			fmt.Printf("Deposit 2: %+v\n", err)
		} else {
			fmt.Printf("Deposit 2:  %+v\n", transaction2)
		}
	}

	// STEP 9 - Sign post fund state by broker 1
	signature1, err = state.Sign(&c, c.PostFundState(), Broker1)
	if err != nil {
		panic(err)
	}

	// STEP 10 - Sign post fund state by broker 2
	signature2, err = state.Sign(&c, c.PostFundState(), Broker2)
	if err != nil {
		panic(err)
	}

	str, err := client.Contract.UnpackStatus(nil, c.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contract On-Chain Status: %+v\n", str)

	// STEP 11 - Interact between participants off-chain
	if c.PostFundComplete() {
		initialTurnNum := len(c.SignedStateForTurnNum)
		for i := initialTurnNum; i < MaxTurnNum; i++ {
			fmt.Printf("\n\nTurn Number: %d\n", i)

			// Build new state
			newState := state.Build(
				ChainId, AssetAddress,
				initialAmountForBroker1, initialAmountForBroker2,
				Broker1, Broker2, false, uint64(i),
			)

			c.SignedStateForTurnNum[uint64(i)] = states.NewSignedState(newState)

			// Sign state by broker1
			signature1, err = state.Sign(&c, newState, Broker1)
			if err != nil {
				panic(err)
			}

			fmt.Printf("State completion after 1 participant: %v\n", c.CurrentStateComplete(newState))
			fmt.Printf("State signed by Broker1: %v\n", c.CurrentStateSignedByMe(newState))

			// Sign state by broker2
			signature2, err = state.Sign(&c, newState, Broker2)
			if err != nil {
				panic(err)
			}

			fmt.Printf("State completion after 2 participant: %v\n", c.CurrentStateComplete(newState))
		}
	}

	// STEP 12 - Finalize a channel
	finalTurnNum := uint64(MaxTurnNum)

	finalState := state.Build(
		ChainId, AssetAddress,
		initialAmountForBroker1, initialAmountForBroker2,
		Broker1, Broker2, true, finalTurnNum,
	)

	c.SignedStateForTurnNum[finalTurnNum] = states.NewSignedState(finalState)

	signature1, err = state.Sign(&c, finalState, Broker1)
	if err != nil {
		panic(err)
	}

	signature2, err = state.Sign(&c, finalState, Broker2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n Channel signatures: %+v\n\n", c.SignedStateForTurnNum)

	concludeParams, err := state.BuildConcludeParams(finalState, signature1, signature2)
	if err != nil {
		panic(err)
	}

	signerFn := Broker1.SignTransaction
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: Broker1.Address, Signer: signerFn}
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
		panic(err)
	}

	fmt.Printf("Conclude transaction: %+v\n", concludeTransaction)
}
