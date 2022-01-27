package commander

import (
	"context"
	"fmt"
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
	err = c.client.GetRevertMessage(tx, receipt)
	return fmt.Errorf("%w tx_hash=%s", err, txHash.String())
}
