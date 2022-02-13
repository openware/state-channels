package contract

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

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
	ConcludeAndTransferAllAssets(opts *bind.TransactOpts, largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appPartHash [32]byte, outcomeBytes []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error)
	GetChainID(opts *bind.CallOpts) (*big.Int, error)
}

type Client struct {
	Contract StateChannelContract
	ChainID  *big.Int
}

func NewClient(contractAddr, rpcUrl string) (Client, error) {
	contractAddress := common.HexToAddress(contractAddr)
	ethClient, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return Client{}, err
	}

	adjucator, err := NewNitroAdjucator(contractAddress, ethClient)
	if err != nil {
		return Client{}, err
	}

	chainID, err := adjucator.GetChainID(nil)
	if err != nil {
		return Client{}, err
	}

	return Client{
		Contract: adjucator,
		ChainID:  chainID,
	}, nil
}
