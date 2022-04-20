package gasprice

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

// Station represents information about gas price, gas limit.
type Station struct {
	GasPrice *big.Int
	GasLimit uint64
}

// Calculate calculates gas price.
func Calculate(ethClient *ethclient.Client) (*big.Int, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	gasPrice, err := ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	return gasPrice, nil
}
