package deployer

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var (
	PollInterval             = 500 * time.Millisecond
	ChainTimeout             = 5 * time.Minute
	ErrWaitToBeMinedTimedOut = errors.New("timeout on waiting for transaction to be mined") // TODO-API here
)

func WaitToBeMined(r ReceiptProvider, tx *types.Transaction) (*types.Receipt, error) {
	timeout := time.After(ChainTimeout)
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	for {
		receipt, err := r.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil && err != ethereum.NotFound {
			return nil, errors.WithStack(err)
		}
		if receipt != nil && receipt.BlockNumber != nil {
			return receipt, nil
		}

		select {
		case <-timeout:
			return nil, errors.WithStack(ErrWaitToBeMinedTimedOut)
		case <-ticker.C:
		}
	}
}
