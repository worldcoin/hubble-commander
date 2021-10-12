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
	PollInterval = 500 * time.Millisecond
	ChainTimeout = 5 * time.Minute

	ErrWaitToBeMinedTimeout = fmt.Errorf("timeout on waiting for transcation to be mined")
)

func WaitToBeMined(c ChainBackend, tx *types.Transaction) (*types.Receipt, error) {
	immediately := make(chan struct{}, 1)
	immediately <- struct{}{}
	defer close(immediately)

	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	timeout := time.After(ChainTimeout)

	for {
		select {
		case <-immediately:
			receipt, err := c.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil && err != ethereum.NotFound {
				return nil, errors.WithStack(err)
			}
			if receipt != nil && receipt.BlockNumber != nil {
				return receipt, nil
			}
		case <-ticker.C:
			receipt, err := c.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil && err != ethereum.NotFound {
				return nil, errors.WithStack(err)
			}
			if receipt != nil && receipt.BlockNumber != nil {
				return receipt, nil
			}
		case <-timeout:
			err := errors.WithStack(ErrWaitToBeMinedTimeout)
			log.Warnf("%+v", err)
			return nil, err
		}
	}
}
