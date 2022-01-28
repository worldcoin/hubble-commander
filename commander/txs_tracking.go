package commander

import (
	"context"
	"fmt"
	"time"

	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (c *Commander) txsTracking(txsHashChan <-chan common.Hash) error {
	for {
		select {
		case <-c.workersContext.Done():
			return nil
		case txHash := <-txsHashChan:
			err := c.waitUntilTxMinedAndCheckForFail(txHash)
			if err != nil {
				panic(err)
			}
		default:
			time.Sleep(time.Millisecond * 300)
		}
	}
}

func (c *Commander) waitUntilTxMinedAndCheckForFail(txHash common.Hash) error {
	tx, _, err := c.client.Blockchain.GetBackend().TransactionByHash(context.Background(), txHash)
	if err != nil {
		return errors.WithStack(err)
	}

	receipt, err := chain.WaitToBeMined(c.client.Blockchain.GetBackend(), tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if receipt.Status == 1 {
		return nil
	}
	err = c.client.GetRevertMessage(tx, receipt)
	return fmt.Errorf("%w tx_hash=%s", err, txHash.String())
}
