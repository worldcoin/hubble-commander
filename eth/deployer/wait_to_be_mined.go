package deployer

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var (
	PollInterval = 500 * time.Millisecond
	ChainTimeout = 5 * time.Minute
)

func WaitToBeMined(c ChainBackend, tx *types.Transaction) (*types.Receipt, error) {
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			receipt, err := c.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil && err != ethereum.NotFound {
				return nil, errors.WithStack(err)
			}
			if receipt != nil && receipt.BlockNumber != nil {
				return receipt, nil
			}
		case <-time.After(ChainTimeout):
			return nil, errors.Errorf("timeout on waiting for transcation to be mined")
		}
	}
}
