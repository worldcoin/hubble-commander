package commander

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	var syncedBlock *uint64

	for {
		select {
		case <-c.stopChannel:
			return nil
		case err = <-subscription.Err():
			return err
		case newBlock := <-blocks:
			latestBlockNumber := newBlock.Number.Uint64()
			c.storage.SetLatestBlockNumber(uint32(latestBlockNumber))
			syncedBlock, err = c.storage.GetSyncedBlock(c.client.ChainState.ChainID)
			if err != nil {
				return err
			}
			endBlock := min(latestBlockNumber, *syncedBlock+uint64(c.cfg.Rollup.SyncSize))

			isProposer, err = c.client.IsActiveProposer()
			if err != nil {
				return errors.WithStack(err)
			}
			c.storage.SetProposer(isProposer)
			err = c.syncAccounts(*syncedBlock, endBlock)
			if err != nil {
				return errors.WithStack(err)
			}
			err = c.syncBatches(*syncedBlock, endBlock)
			if err != nil {
				return errors.WithStack(err)
			}
			err = c.storage.SetSyncedBlock(c.client.ChainState.ChainID, newBlock.Number.Uint64())
			if err != nil {
				return err
			}
		}
	}
}

func (c *Commander) syncBatches(startBlock, endBlock uint64) (err error) {
	transactionExecutor, err := newTransactionExecutor(c.storage, c.client, c.cfg.Rollup)
	if err != nil {
		return err
	}
	defer transactionExecutor.Rollback(&err)

	err = transactionExecutor.SyncBatches(startBlock, endBlock)
	if err != nil {
		return err
	}
	return transactionExecutor.Commit()
}

func (c *Commander) syncAccounts(start, end uint64) error {
	it, err := c.client.AccountRegistry.FilterPubkeyRegistered(&bind.FilterOpts{
		Start: start,
		End:   &end,
	})
	if err != nil {
		return err
	}
	defer it.Close()
	for it.Next() {
		err = ProcessPubkeyRegistered(c.storage, it.Event)
		if err != nil {
			return err
		}
	}
	return nil
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
