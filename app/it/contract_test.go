package it

import (
	"app/models/broker"
	"app/pkg/contract"
	"app/pkg/parser"
	"app/pkg/state"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/channel"
	states "github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/types"
	"github.com/stretchr/testify/suite"
)

const MaxRandNum = 100000000000

var (
	AssetAddress = common.HexToAddress("0x00")
	GasLimit     = uint64(6721975)
	GasPrice     = big.NewInt(20000000000)
)

type ContractTestSuite struct {
	suite.Suite
	ContractClient contract.Client
	Brokers        []broker.Broker
}

func (s *ContractTestSuite) SetupSuite() {
	mydir, err := os.Getwd()
	s.Require().NoError(err)

	// List of brokers
	vaultAccount, err := parser.ToVaultAccount(mydir + "/../../contracts/accounts.json")
	s.Require().NoError(err)

	brokers := []broker.Broker{}
	for k, v := range vaultAccount.Accounts {
		privateKey := strings.TrimPrefix(v.PrivateKey, "0x")
		broker := broker.New(common.HexToAddress(v.Address), types.Destination(common.HexToHash(v.Address)), common.Hex2Bytes(privateKey), uint(k))
		brokers = append(brokers, *broker)
	}
	s.Brokers = brokers

	// SC client
	contractObj, err := parser.ToContract(mydir + "/../../contracts/addresses.json")
	s.Require().NoError(err)
	contractAddress := contractObj.ChainIds[0].SC.NitroAdj.Address
	client, err := contract.NewClient(contractAddress, "http://127.0.0.1:8545")
	s.Require().NoError(err)
	s.ContractClient = client
}

func (s *ContractTestSuite) openChannel(amount1, amount2 *big.Int) channel.Channel {
	randNum := rand.Intn(MaxRandNum)
	channelNonce := big.NewInt(int64(randNum))
	preFundState := state.Build(
		s.ContractClient.ChainID, channelNonce, AssetAddress,
		amount1, amount2,
		&s.Brokers[0], &s.Brokers[1], false, 0,
	)

	// Open a channel between brokers
	c, err := channel.New(preFundState, s.Brokers[0].Role)
	s.Require().NoError(err)

	return c
}

func (s *ContractTestSuite) buildFinalState(c *channel.Channel, amount1, amount2 *big.Int, final bool, finalTurnNum uint64) states.State {
	finalState := state.Build(
		s.ContractClient.ChainID, c.ChannelNonce, AssetAddress,
		amount1, amount2,
		&s.Brokers[0], &s.Brokers[1], final, finalTurnNum,
	)
	c.SignedStateForTurnNum[finalTurnNum] = states.NewSignedState(finalState)

	return finalState
}

func (s *ContractTestSuite) concludeChannel(c *channel.Channel, amount1, amount2 *big.Int, finalTurnNum uint64) {
	finalState := s.buildFinalState(c, amount1, amount2, true, finalTurnNum)

	// Sign states by both participants
	signature1, err := s.Brokers[0].SignState(c, finalState)
	s.Require().NoError(err)

	signature2, err := s.Brokers[1].SignState(c, finalState)
	s.Require().NoError(err)

	// Build conclude params
	concludeParams, err := state.BuildConcludeParams(finalState, signature1, signature2)
	s.Require().NoError(err)
	signerFn := s.Brokers[0].SignTransaction
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: s.Brokers[0].Address, Signer: signerFn}

	// Conclude channel
	_, err = s.ContractClient.Contract.ConcludeAndTransferAllAssets(&transactionOpts,
		big.NewInt(int64(finalTurnNum)),
		concludeParams.FixedPart,
		concludeParams.AppPart,
		concludeParams.OutcomeState,
		concludeParams.NumStates,
		concludeParams.WhoSignedWhat,
		concludeParams.Signatures,
	)
	s.Require().NoError(err)
}

func (s *ContractTestSuite) TestDepositSuccess() {
	// Unique channel ID
	amount := big.NewInt(100)
	channel := s.openChannel(amount, amount)
	validChannelID := channel.Id
	signerFn := s.Brokers[0].SignTransaction
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: s.Brokers[0].Address, Signer: signerFn, Value: amount}

	// expectedHeld = 0
	_, err := s.ContractClient.Contract.Deposit(&transactionOpts, AssetAddress, validChannelID, big.NewInt(0), amount)
	s.Require().NoError(err)

	// expectedHeld = initial amount
	_, err = s.ContractClient.Contract.Deposit(&transactionOpts, AssetAddress, validChannelID, amount, amount)
	s.Require().NoError(err)

	// To make this test work every time, system need to conclude channel with right params
	s.concludeChannel(&channel, amount.Add(amount, amount), big.NewInt(0), uint64(1))
}

func (s *ContractTestSuite) TestDepositErrorSufficientHoldings() {
	// Unique channel ID
	amount := big.NewInt(100)
	channel := s.openChannel(amount, amount)
	validChannelID := channel.Id

	signerFn := s.Brokers[0].SignTransaction
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: s.Brokers[0].Address, Signer: signerFn, Value: amount}

	// Deposit initial funds to SC
	_, err := s.ContractClient.Contract.Deposit(&transactionOpts, AssetAddress, validChannelID, big.NewInt(0), amount)
	s.Require().NoError(err)

	// Deposit funds once more time, when channel have already been funded
	_, err = s.ContractClient.Contract.Deposit(&transactionOpts, AssetAddress, validChannelID, big.NewInt(0), amount)
	s.Require().Error(err, "holdings already sufficient")

	// To make this test work every time, system need to conclude channel with right params
	s.concludeChannel(&channel, amount, big.NewInt(0), uint64(1))
}

func (s *ContractTestSuite) TestDepositErrorInsufficientHoldings() {
	// Unique channel ID
	amount := big.NewInt(100)
	channel := s.openChannel(amount, amount)
	validChannelID := channel.Id

	signerFn := s.Brokers[0].SignTransaction
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: s.Brokers[0].Address, Signer: signerFn, Value: amount}

	// Expected held >= initial amount
	_, err := s.ContractClient.Contract.Deposit(&transactionOpts, AssetAddress, validChannelID, amount, amount)
	s.Require().Error(err, "holdings < expectedHeld")
}

func (s *ContractTestSuite) TestDepositErrorInvalidGas() {
	// Unique channel ID
	amount := big.NewInt(100)
	channel := s.openChannel(amount, amount)
	validChannelID := channel.Id

	signerFn := s.Brokers[0].SignTransaction
	GasPrice := big.NewInt(100)

	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: s.Brokers[0].Address, Signer: signerFn, Value: amount}
	_, err := s.ContractClient.Contract.Deposit(&transactionOpts, AssetAddress, validChannelID, big.NewInt(0), amount)
	s.Require().Error(err, "Transaction gasPrice (100) is too low for the next block, which has a baseFeePerGas of 72519483")
}

func (s *ContractTestSuite) TestConcludeChannelSuccess() {
	// Open a channel between brokers
	initialAmount1 := big.NewInt(0)
	initialAmount2 := big.NewInt(0)
	c := s.openChannel(initialAmount1, initialAmount2)

	// Build final state, with finalized = true
	finalTurnNum := uint64(1)
	finalState := s.buildFinalState(&c, initialAmount1, initialAmount2, true, finalTurnNum)

	// Sign states by both participants
	signature1, err := s.Brokers[0].SignState(&c, finalState)
	s.Require().NoError(err)

	signature2, err := s.Brokers[1].SignState(&c, finalState)
	s.Require().NoError(err)

	// Build conclude params
	concludeParams, err := state.BuildConcludeParams(finalState, signature1, signature2)
	s.Require().NoError(err)

	signerFn := s.Brokers[0].SignTransaction
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: s.Brokers[0].Address, Signer: signerFn}

	// Conclude channel
	_, err = s.ContractClient.Contract.ConcludeAndTransferAllAssets(&transactionOpts,
		big.NewInt(int64(finalTurnNum)),
		concludeParams.FixedPart,
		concludeParams.AppPart,
		concludeParams.OutcomeState,
		concludeParams.NumStates,
		concludeParams.WhoSignedWhat,
		concludeParams.Signatures,
	)
	s.Require().NoError(err)
}

func (s *ContractTestSuite) TestConcludeChannelErrorChannelIsNotFinalized() {
	// Open a channel between brokers
	initialAmount1 := big.NewInt(200)
	initialAmount2 := big.NewInt(300)
	c := s.openChannel(initialAmount1, initialAmount2)

	// Build final state, with finalized = false
	finalTurnNum := uint64(1)
	finalState := s.buildFinalState(&c, initialAmount1, initialAmount2, false, finalTurnNum)

	// Sign states by both participants
	signature1, err := s.Brokers[0].SignState(&c, finalState)
	s.Require().NoError(err)

	signature2, err := s.Brokers[1].SignState(&c, finalState)
	s.Require().NoError(err)

	// Build conclude params
	concludeParams, err := state.BuildConcludeParams(finalState, signature1, signature2)
	s.Require().NoError(err)

	signerFn := s.Brokers[0].SignTransaction
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: s.Brokers[0].Address, Signer: signerFn}

	// Conclude channel
	_, err = s.ContractClient.Contract.ConcludeAndTransferAllAssets(&transactionOpts,
		big.NewInt(int64(finalTurnNum)),
		concludeParams.FixedPart,
		concludeParams.AppPart,
		concludeParams.OutcomeState,
		concludeParams.NumStates,
		concludeParams.WhoSignedWhat,
		concludeParams.Signatures,
	)
	s.Require().Error(err, "Invalid signatures / !isFinal")
}

func (s *ContractTestSuite) TestConcludeChannelErrorStateSignedByOne() {
	// Open a channel between brokers
	initialAmount1 := big.NewInt(200)
	initialAmount2 := big.NewInt(300)
	c := s.openChannel(initialAmount1, initialAmount2)

	// Build final state, with finalized = true
	finalTurnNum := uint64(1)
	finalState := s.buildFinalState(&c, initialAmount1, initialAmount2, true, finalTurnNum)

	// Sign states by both participants
	signature1, err := s.Brokers[0].SignState(&c, finalState)
	s.Require().NoError(err)

	// No state signature
	signature2 := states.Signature{}

	// Build conclude params
	concludeParams, err := state.BuildConcludeParams(finalState, signature1, signature2)
	s.Require().NoError(err)

	signerFn := s.Brokers[0].SignTransaction
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: s.Brokers[0].Address, Signer: signerFn}

	// Conclude channel
	_, err = s.ContractClient.Contract.ConcludeAndTransferAllAssets(&transactionOpts,
		big.NewInt(int64(finalTurnNum)),
		concludeParams.FixedPart,
		concludeParams.AppPart,
		concludeParams.OutcomeState,
		concludeParams.NumStates,
		concludeParams.WhoSignedWhat,
		concludeParams.Signatures,
	)
	s.Require().Error(err, "Invalid signature")
}

func (s *ContractTestSuite) TestConcludeChannelErrorInvalidGasPrice() {
	// Open a channel between brokers
	initialAmount1 := big.NewInt(200)
	initialAmount2 := big.NewInt(300)
	c := s.openChannel(initialAmount1, initialAmount2)

	// Build final state, with finalized = true
	finalTurnNum := uint64(1)
	finalState := s.buildFinalState(&c, initialAmount1, initialAmount2, true, finalTurnNum)

	// Sign states by both participants
	signature1, err := s.Brokers[0].SignState(&c, finalState)
	s.Require().NoError(err)

	// No state signature
	signature2, err := s.Brokers[1].SignState(&c, finalState)
	s.Require().NoError(err)

	// Build conclude params
	concludeParams, err := state.BuildConcludeParams(finalState, signature1, signature2)
	s.Require().NoError(err)

	signerFn := s.Brokers[0].SignTransaction
	GasPrice := big.NewInt(100)
	transactionOpts := bind.TransactOpts{GasLimit: GasLimit, GasPrice: GasPrice, From: s.Brokers[0].Address, Signer: signerFn}

	// Conclude channel
	_, err = s.ContractClient.Contract.ConcludeAndTransferAllAssets(&transactionOpts,
		big.NewInt(int64(finalTurnNum)),
		concludeParams.FixedPart,
		concludeParams.AppPart,
		concludeParams.OutcomeState,
		concludeParams.NumStates,
		concludeParams.WhoSignedWhat,
		concludeParams.Signatures,
	)
	s.Require().Error(err, "Transaction gasPrice (100) is too low for the next block, which has a baseFeePerGas of 36641")
}

func TestContractTestSuite(t *testing.T) {
	suite.Run(t, &ContractTestSuite{})
}
