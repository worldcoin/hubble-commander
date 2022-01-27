package commander

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Commander) txsTracking() error {
	for {
		select {
		case <-c.workersContext.Done():
			return nil
		case txHash := <-c.client.TxsHashesChan:
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
	var receipt *types.Receipt

	for {
		receipt, err = c.client.Blockchain.GetBackend().TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return errors.WithStack(err)
		}
		if receipt != nil {
			break
		}
		time.Sleep(time.Millisecond * 300)
	}

	if receipt.Status == 1 {
		return nil
	}
	return c.client.GetRevertMessage(tx, receipt)
}
