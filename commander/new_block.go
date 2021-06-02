package commander

import (
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Commander) newBlockLoop() error {
	endBlock := ref.Uint64(0)
	cancelRollup := make(chan struct{}, 1)
	latestBlockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
	if err != nil {
		return err
	}

	for *endBlock != *latestBlockNumber {
		endBlock, err = c.newBlockIteration(cancelRollup, *latestBlockNumber)
		if err != nil {
			return err
		}

		latestBlockNumber, err = c.client.ChainConnection.GetLatestBlockNumber()
		if err != nil {
			return err
		}
		select {
		case <-c.stopChannel:
			return nil
		default:
			continue
		}
	}

	blocks := make(chan *types.Header)
	subscription, err := c.client.ChainConnection.SubscribeNewHead(blocks)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case <-c.stopChannel:
			return nil
		case err = <-subscription.Err():
			return err
		case newBlock := <-blocks:
			currentBlockNumber := newBlock.Number.Uint64()
			c.storage.SetLatestBlockNumber(uint32(currentBlockNumber))
			_, err = c.newBlockIteration(cancelRollup, currentBlockNumber)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Commander) newBlockIteration(cancelRollup chan struct{}, latestBlockNumber uint64) (*uint64, error) {
	syncedBlock, err := c.storage.GetSyncedBlock(c.client.ChainState.ChainID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	endBlock := min(latestBlockNumber, *syncedBlock+uint64(c.cfg.Rollup.SyncSize))

	isProposer, err := c.client.IsActiveProposer()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = c.syncAccounts(*syncedBlock, endBlock)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = c.syncBatches(*syncedBlock, endBlock)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = c.storage.SetSyncedBlock(c.client.ChainState.ChainID, endBlock)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if endBlock == latestBlockNumber {
		c.manageRollupLoop(isProposer, cancelRollup)
	}
	return &endBlock, nil
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

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
