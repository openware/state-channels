package nitro

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// StateChannelContract represents available functions from Nitro protocol
type StateChannelContract interface {
	Transfer(opts *bind.TransactOpts, assetIndex *big.Int, fromChannelId [32]byte, outcomeBytes []byte, stateHash [32]byte, indices []*big.Int) (*types.Transaction, error)
	TransferAllAssets(opts *bind.TransactOpts, channelId [32]byte, outcomeBytes []byte, stateHash [32]byte) (*types.Transaction, error)
	Deposit(opts *bind.TransactOpts, asset common.Address, channelId [32]byte, expectedHeld *big.Int, amount *big.Int) (*types.Transaction, error)
	ValidTransition(opts *bind.CallOpts, nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error)
	UnpackStatus(opts *bind.CallOpts, channelId [32]byte) (struct {
		TurnNumRecord *big.Int
		FinalizesAt   *big.Int
		Fingerprint   *big.Int
	}, error)
	Checkpoint(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error)
	ConcludeAndTransferAllAssets(opts *bind.TransactOpts, largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcomeBytes []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error)
	GetChainID(opts *bind.CallOpts) (*big.Int, error)
	Holdings(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte) (*big.Int, error)
}

// Client stores information about adjudicator and chainID
type Client struct {
	Adjudicator StateChannelContract
	ChainID     *big.Int
	Eth         *ethclient.Client
}

// NewClient returns a new Client from supplied params.
func NewClient(contractAddr, rpcUrl string) (Client, error) {
	contractAddress := common.HexToAddress(contractAddr)
	ethClient, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return Client{}, err
	}

	adjudicator, err := NewNitroAdjudicator(contractAddress, ethClient)
	if err != nil {
		return Client{}, err
	}

	chainID, err := adjudicator.GetChainID(nil)
	if err != nil {
		return Client{}, err
	}

	return Client{
		Adjudicator: adjudicator,
		Eth:         ethClient,
		ChainID:     chainID,
	}, nil
}
