package protocol

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GasStation represents information about gas price, gas limit.
type GasStation struct {
	GasPrice *big.Int
	GasLimit uint64
}

// signTransaction adds a participant signature to the transaction.
// An error is thrown if the signature is invalid.
func signTransaction(chainID *big.Int, privateKey []byte) (signerFn bind.SignerFn) {
	signerFn = func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		signer := types.LatestSignerForChainID(chainID)
		hash := signer.Hash(tx)

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

	return
}

// CalculateGasPrice calculates gas price.
func CalculateGasPrice(ethClient *ethclient.Client) (*big.Int, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	gasPrice, err := ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	return gasPrice, nil
}
