package main

import (
	"app/models/broker"
	"app/pkg/contract"
	"fmt"
	"math/big"
	"time"

	ethAbi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/statechannels/go-nitro/abi"
	"github.com/statechannels/go-nitro/channel"
	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	typ "github.com/statechannels/go-nitro/types"
)

var (
	DefaultTimeout = 3 * time.Second
	AssetAddress   = common.HexToAddress("0x00")
	MaxTurnNum     = 5
	GasLimit       = uint64(6721975)

	// TODO Read from file
	RpcUrl          = "http://127.0.0.1:8545"
	ContractAddress = "0xCc388ae2496E15ff8C6df70566171c750B5118E2"
	Broker1         = broker.New(
		common.HexToAddress(`0x7d6fe92F348B6F2216A7AA2c2F0Dd9b8c830e490`),
		typ.AddressToDestination(common.HexToAddress(`0x566D32D5b8F3DC45851f6DC43533EE28DD3C43d5`)),
		common.Hex2Bytes(`9a14bf0eb618a3407a12a83a74dfe7bbed098ccc6347985b92ab08e81996cfc9`),
		0,
	)

	Broker2 = broker.New(
		common.HexToAddress(`0xb1239c28162bf9b3e2aa6Dc2c78066B26D5423F7`),
		typ.AddressToDestination(common.HexToAddress(`0xb1239c28162bf9b3e2aa6Dc2c78066B26D5423F7`)),
		common.Hex2Bytes(`2305f34d1dcab90a5143856446d5213c0b29ae353a25845445c8050d4bca38d9`),
		1,
	)
)

func BuildState(
	chainID, channelNonce *big.Int,
	amountForBroker1, amountForBroker2 *big.Int,
	isFinal bool, turnNum uint64) state.State {

	state := state.State{
		ChainId:           chainID,
		Participants:      []typ.Address{Broker1.Address, Broker2.Address},
		ChannelNonce:      channelNonce,
		ChallengeDuration: big.NewInt(60),
		AppData:           []byte{},
		Outcome: outcome.Exit{
			outcome.SingleAssetExit{
				Asset: AssetAddress,
				Allocations: outcome.Allocations{
					outcome.Allocation{
						Destination: Broker1.Destination,
						Amount:      amountForBroker1,
					},
					outcome.Allocation{
						Destination: Broker2.Destination,
						Amount:      amountForBroker2,
					},
				},
			},
		},
		TurnNum: turnNum,
		IsFinal: isFinal,
	}

	return state
}

func SignState(c *channel.Channel, newState state.State, broker broker.Broker) (state.Signature, error) {
	signature, err := newState.Sign(broker.PrivateKey)
	c.AddStateWithSignature(newState, signature)
	if err != nil {
		return signature, err
	}
	return signature, nil
}

func SignTransaction(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	signer := types.LatestSignerForChainID(big.NewInt(1))
	hash := signer.Hash(tx)

	var privateKey []byte
	if address == Broker1.Address {
		privateKey = Broker1.PrivateKey
	} else {
		privateKey = Broker2.PrivateKey
	}

	prv, err := crypto.ToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	signature, err := crypto.Sign(hash.Bytes(), prv)
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(signer, signature)
}

// Simple state channel example
func main() {
	// STEP 1 - initialize client
	client, err := contract.NewClient(ContractAddress, RpcUrl)
	if err != nil {
		panic(err)
	}

	// STEP 2 - Get Chain ID from the client
	ChainId := client.ChainId
	fmt.Printf("Chain Id: %v\n", ChainId)

	// STEP 3 - Initialize pre fund state to sign
	channelNonce := big.NewInt(1) //big.NewInt(time.Now().UnixNano())
	initialAmountForBroker1 := big.NewInt(100)
	initialAmountForBroker2 := big.NewInt(200)
	preFundState := BuildState(ChainId, channelNonce, initialAmountForBroker1, initialAmountForBroker2, false, 0)

	// STEP 4 - Open a channel between brokers
	c, err := channel.New(preFundState, Broker1.Role)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Channel ID: %+v \n", c.Id)

	// STEP 5 - Sign prefund state by broker 1
	signature1, err := SignState(&c, c.PreFundState(), *Broker1)
	if err != nil {
		panic(err)
	}

	// STEP 6 - Sign prefund state by broker 2
	signature2, err := SignState(&c, c.PreFundState(), *Broker2)
	if err != nil {
		panic(err)
	}

	// STEP 7 - Deposit process. Deposit funds only when prefund state has been completed
	signerFn := SignTransaction
	if c.PreFundComplete() {
		// STEP 7.1 - Broker 1 deposit funds to smart-contract
		transactionOpts1 := bind.TransactOpts{From: Broker1.Address, Signer: signerFn, Value: initialAmountForBroker1}
		transaction1, err := client.Contract.Deposit(&transactionOpts1, AssetAddress, c.Id, big.NewInt(0), initialAmountForBroker1)
		if err != nil {
			fmt.Printf("Deposit 1: %+v\n", err)
		} else {
			fmt.Printf("Deposit 1:  %+v\n", transaction1)
		}

		// STEP 7.2 - Broker 2 deposit funds to smart-contract
		// Expected held should be the same as Broker2's initial amount of funds
		// If there is no funds, SC'll revert transaction
		transactionOpts2 := bind.TransactOpts{From: Broker2.Address, Signer: signerFn, Value: initialAmountForBroker2}
		transaction2, err := client.Contract.Deposit(&transactionOpts2, AssetAddress, c.Id, initialAmountForBroker1, initialAmountForBroker2)
		if err != nil {
			fmt.Printf("Deposit 2: %+v\n", err)
		} else {
			fmt.Printf("Deposit 2:  %+v\n", transaction2)
		}
	}

	// STEP 8 - Sign post fund state by broker 1
	signature1, err = SignState(&c, c.PostFundState(), *Broker1)
	if err != nil {
		panic(err)
	}

	// STEP 9 - Sign post fund state by broker 2
	signature2, err = SignState(&c, c.PostFundState(), *Broker2)
	if err != nil {
		panic(err)
	}

	str, err := client.Contract.UnpackStatus(nil, c.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contract Unpack Status after post fund state: %+v\n", str)

	// STEP 9 - Interact between participants off-chain
	if c.PostFundComplete() {
		initialTurnNum := len(c.SignedStateForTurnNum)
		for i := initialTurnNum; i < MaxTurnNum; i++ {
			fmt.Printf("\n\nTurn Number: %d\n", i)

			amountForBroker1 := initialAmountForBroker1
			amountForBroker2 := initialAmountForBroker2
			newState := BuildState(ChainId, channelNonce, amountForBroker1, amountForBroker2, false, uint64(i))

			c.SignedStateForTurnNum[uint64(i)] = state.NewSignedState(newState)

			signature1, err = SignState(&c, newState, *Broker1)
			if err != nil {
				panic(err)
			}

			fmt.Printf("State completion after 1 participant: %v\n", c.CurrentStateComplete(newState))
			fmt.Printf("State signed by Broker1: %v\n", c.CurrentStateSignedByMe(newState))

			signature2, err = SignState(&c, newState, *Broker2)
			if err != nil {
				panic(err)
			}

			fmt.Printf("State completion after 2 participant: %v\n", c.CurrentStateComplete(newState))
		}
	}

	// STEP 10 - Finalize a channel
	finalTurnNum := uint64(MaxTurnNum)
	finalState := BuildState(ChainId, channelNonce, initialAmountForBroker1, initialAmountForBroker2, true, finalTurnNum)
	c.SignedStateForTurnNum[finalTurnNum] = state.NewSignedState(finalState)

	signature1, err = SignState(&c, finalState, *Broker1)
	if err != nil {
		panic(err)
	}

	signature2, err = SignState(&c, finalState, *Broker2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n Channel signatures: %+v\n\n", c.SignedStateForTurnNum)

	outcomeState, err := finalState.Outcome.Hash()
	if err != nil {
		panic(err)
	}

	encoded, err := ethAbi.Arguments{
		{Type: abi.Uint256},
		{Type: abi.Address},
		{Type: abi.Bytes},
	}.Pack(c.ChallengeDuration, c.AppDefinition, []byte(finalState.AppData))
	if err != nil {
		panic(err)
	}
	appPartHash := crypto.Keccak256Hash(encoded)

	var signatureR1, signatureS1, signatureR2, signatureS2 [32]byte

	copy(signatureR1[:], signature1.R)
	copy(signatureS1[:], signature1.S)
	copy(signatureR2[:], signature2.R)
	copy(signatureS2[:], signature2.S)

	sigs := []contract.IForceMoveSignature{
		{V: signature1.V, R: signatureR1, S: signatureS1},
		{V: signature2.V, R: signatureR2, S: signatureS2},
	}

	fixedPart := contract.IForceMoveFixedPart{
		ChainId:           ChainId,
		Participants:      c.Participants,
		ChannelNonce:      c.ChannelNonce,
		AppDefinition:     c.AppDefinition,
		ChallengeDuration: c.ChallengeDuration,
	}
	numStates := uint8(len(c.Participants) - 1)
	whoSignedWhat := []uint8{uint8(0), uint8(0)}
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, From: Broker1.Address, Signer: signerFn}

	concludeTransaction, err := client.Contract.ConcludeAndTransferAllAssets(&transactionOpts,
		big.NewInt(int64(finalTurnNum)),
		fixedPart,
		appPartHash,
		outcomeState.Bytes(),
		numStates,
		whoSignedWhat,
		sigs,
	)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	fmt.Printf("Conclude transaction: %+v\n", concludeTransaction)
}
