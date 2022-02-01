package commander

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Commander) txsTracking(txsChan <-chan *types.Transaction) error {
	for {
		select {
		case <-c.workersContext.Done():
			return nil
		case tx := <-txsChan:
			err := c.waitUntilTxMinedAndCheckForFail(tx)
			if err != nil {
				panic(err)
			}
		default:
			time.Sleep(time.Millisecond * 300)
		}
	}
}

func (c *Commander) waitUntilTxMinedAndCheckForFail(tx *types.Transaction) error {
	receipt, err := c.client.WaitToBeMined(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if receipt.Status == 1 {
		return nil
	}
	err = c.client.GetRevertMessage(tx, receipt)
	return fmt.Errorf("%w tx_hash=%s", err, tx.Hash().String())
}
