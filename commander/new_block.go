package commander

import (
	"context"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Commander) newBlockLoop() error {
	err := c.initialSync()
	if err != nil {
		return err
	}

	blocks := make(chan *types.Header)
	subscription, err := c.client.ChainConnection.SubscribeNewHead(blocks)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	var rollupCancel context.CancelFunc
	for {
		select {
		case <-c.stopChannel:
			return nil
		case err = <-subscription.Err():
			return err
		case <-blocks:
			latestBlockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
			if err != nil {
				return err
			}
			c.storage.SetLatestBlockNumber(uint32(*latestBlockNumber))

			isProposer, err := c.client.IsActiveProposer()
			if err != nil {
				return errors.WithStack(err)
			}
			syncedBlock, err := c.syncForward(*latestBlockNumber, isProposer)
			if err != nil {
				return err
			}

			if *syncedBlock == *latestBlockNumber {
				rollupCancel = c.manageRollupLoop(rollupCancel, isProposer)
			}
		}
	}
}

func (c *Commander) initialSync() error {
	latestBlockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
	if err != nil {
		return err
	}
	c.storage.SetLatestBlockNumber(uint32(*latestBlockNumber))

	syncedBlock := ref.Uint64(uint64(0))
	for *syncedBlock != *latestBlockNumber {
		syncedBlock, err = c.syncForward(*latestBlockNumber, false)
		if err != nil {
			return err
		}

		latestBlockNumber, err = c.client.ChainConnection.GetLatestBlockNumber()
		if err != nil {
			return err
		}
		c.storage.SetLatestBlockNumber(uint32(*latestBlockNumber))

		select {
		case <-c.stopChannel:
			return nil
		default:
			continue
		}
	}
	return nil
}

func (c *Commander) syncForward(latestBlockNumber uint64, isProposer bool) (*uint64, error) {
	syncedBlock, err := c.storage.GetSyncedBlock(c.client.ChainState.ChainID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	startBlock := *syncedBlock + 1
	endBlock := min(latestBlockNumber, startBlock+uint64(c.cfg.Rollup.SyncSize))

	err = c.syncRange(startBlock, endBlock, isProposer)
	if err != nil {
		return nil, err
	}

	err = c.storage.SetSyncedBlock(c.client.ChainState.ChainID, endBlock)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &endBlock, nil
}

func (c *Commander) syncRange(startBlock, endBlock uint64, isProposer bool) error {
	err := c.syncAccounts(startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}
	if !isProposer {
		err = c.syncBatches(startBlock, endBlock)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Commander) syncBatches(startBlock, endBlock uint64) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
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

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
