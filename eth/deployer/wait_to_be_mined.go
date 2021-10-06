package deployer

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	PollInterval             = 500 * time.Millisecond
	ChainTimeout             = 5 * time.Minute
	ErrWaitToBeMinedTimedOut = fmt.Errorf("timeout on waiting for transaction to be mined")
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
			err = errors.WithStack(ErrWaitToBeMinedTimedOut)
			log.Warnf("%+v", err)
			return nil, err
		case <-ticker.C:
		}
	}
}
