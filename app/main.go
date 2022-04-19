package main

import (
	"app/examples"
	"app/internal/parser"
	"app/pkg/nitro"
	"app/pkg/protocol"
	"math/big"
	"os"
	"strings"

	"github.com/caitlinelfring/go-env-default"
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/types"
)

var (
	NodeUrl          = env.GetDefault("NODE_URL", "http://127.0.0.1:8545")
	AccountsFileName = env.GetDefault("ACCOUNTS_FILENAME", "accounts.json")
	Network          = env.GetDefault("NETWORK", "localhost")
	AssetAddress     = common.HexToAddress("0x0")
	ParticipantCount = 3
)

// State channel examples
func main() {
	// Initialize participants, deploy smart-contract
	myDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vaultAccount, err := parser.ToVaultAccount(myDir + "/../contracts/" + AccountsFileName)
	if err != nil {
		panic(err)
	}

	contractObj, err := parser.ToContract(myDir + "/../contracts/addresses.json")
	if err != nil {
		panic(err)
	}

	var contractAddress string
	for _, obj := range contractObj {
		if obj[0].Name == Network {
			contractAddress = obj[0].SC.NitroAdj.Address
		}
	}

	// Initialize Participants
	var participants []protocol.Participant
	participantPrivateKeys := make(map[protocol.Participant][]byte)

	for i := 0; i < ParticipantCount; i++ {
		vault := vaultAccount.Accounts[i]
		privateKey := strings.TrimPrefix(vault.PrivateKey, "0x")
		amount := big.NewInt(0).Mul(big.NewInt(1+int64(i)), big.NewInt(100))

		participantObj := protocol.NewParticipant(common.HexToAddress(vault.Address), types.Destination(common.HexToHash(vault.Address)), uint(i), amount)
		participants = append(participants, *participantObj)
		participantPrivateKeys[*participantObj] = common.Hex2Bytes(privateKey)
	}

	// Initialize SC client
	client, err := nitro.NewClient(contractAddress, NodeUrl)
	if err != nil {
		panic(err)
	}

	// Initialize contract
	c := protocol.NewContract(client, AssetAddress)

	// Demo example
	err = examples.Demo(participants, participantPrivateKeys, c)
	if err != nil {
		panic(err)
	}
}
