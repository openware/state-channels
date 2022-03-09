package main

import (
	"app/examples"
	"app/models/broker"
	"app/pkg/contract"
	"app/pkg/parser"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

var Broker1, Broker2 *broker.Broker
var ChainID *big.Int
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

// State channel examples
func main() {
	// Initialize participants (brokers), deploy smart-contract
	mydir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vaultAccount, err := parser.ToVaultAccount(mydir + "/../contracts/accounts.json")
	if err != nil {
		panic(err)
	}

	contractObj, err := parser.ToContract(mydir + "/../contracts/addresses.json")
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

	// Initialize SC client
	client, err := contract.NewClient(ContractAddress, RpcUrl)
	if err != nil {
		panic(err)
	}

	// Get Chain ID from the client
	ChainID = client.ChainID
	fmt.Printf("Chain Id: %v\n", ChainID)

	// Simple Example
	fmt.Println("Simple Example\n\n\n")
	err = examples.Simple(Broker1, Broker2, client, AssetAddress, ChainID)
	if err != nil {
		panic(err)
	}

	// Simple Trade Example
	fmt.Println("Simple Trade Example\n\n\n")
	err = examples.SimpleTrade(Broker1, Broker2, client, AssetAddress, ChainID)
	if err != nil {
		panic(err)
	}
}
