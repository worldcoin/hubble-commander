package deployer

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var ChainTimeout = 5 * time.Minute

func WaitToBeMined(c ChainBackend, tx *types.Transaction) (*types.Receipt, error) {
	begin := time.Now()
	for {
		receipt, err := c.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil && err != ethereum.NotFound {
			return nil, errors.WithStack(err)
		}

		if receipt != nil && receipt.BlockNumber != nil {
			return receipt, nil
		}

		if time.Since(begin) > ChainTimeout {
			return nil, errors.Errorf("timeout on waiting for transcation to be mined")
		}
	}
}
