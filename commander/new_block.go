package commander

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Commander) newBlockLoop() error {
	blocks := make(chan *types.Header)
	subscription, err := c.client.ChainConnection.SubscribeNewHead(blocks)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	isProposer, err := c.client.IsActiveProposer()
	if err != nil {
		return errors.WithStack(err)
	}
	c.storage.SetProposer(isProposer)

	for {
		select {
		case <-c.stopChannel:
			return nil
		case err = <-subscription.Err():
			return err
		case newBlock := <-blocks:
			c.storage.SetLatestBlockNumber(uint32(newBlock.Number.Uint64()))

			isProposer, err = c.client.IsActiveProposer()
			if err != nil {
				return errors.WithStack(err)
			}
			c.storage.SetProposer(isProposer)

			err = c.SyncBatches(isProposer)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
}

func (c *Commander) SyncBatches(isProposer bool) (err error) {
	if isProposer {
		return nil
	}

	transactionExecutor, err := newTransactionExecutor(c.storage, c.client, &c.cfg.Rollup)
	if err != nil {
		return err
	}
	defer transactionExecutor.Rollback(&err)
	return transactionExecutor.Commit()
}
