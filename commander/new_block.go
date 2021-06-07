package commander

import (
	"context"
	"math/big"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Commander) newBlockLoop() error {
	latestBlockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
	if err != nil {
		return err
	}

	blocks := make(chan *types.Header, 5)
	// TODO: remove mining every 1s from github actions as it will be no longer needed with below line
	blocks <- &types.Header{Number: new(big.Int).SetUint64(*latestBlockNumber)}
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
		case currentBlock := <-blocks:
			if currentBlock.Number.Uint64() <= uint64(c.storage.GetLatestBlockNumber()) {
				continue
			}

			err = c.syncToLastBlock()
			if err != nil {
				return err
			}

			isProposer, err := c.client.IsActiveProposer()
			if err != nil {
				return errors.WithStack(err)
			}
			rollupCancel = c.manageRollupLoop(rollupCancel, isProposer)
		}
	}
}

func (c *Commander) syncToLastBlock() error {
	latestBlockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
	if err != nil {
		return err
	}
	c.storage.SetLatestBlockNumber(uint32(*latestBlockNumber))

	syncedBlock := ref.Uint64(uint64(0))
	for *syncedBlock != *latestBlockNumber {
		syncedBlock, err = c.syncForward(*latestBlockNumber)
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

func (c *Commander) syncForward(latestBlockNumber uint64) (*uint64, error) {
	syncedBlock, err := c.storage.GetSyncedBlock(c.client.ChainState.ChainID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	startBlock := *syncedBlock + 1
	endBlock := min(latestBlockNumber, startBlock+uint64(c.cfg.Rollup.SyncSize))

	err = c.syncRange(startBlock, endBlock)
	if err != nil {
		return nil, err
	}

	err = c.storage.SetSyncedBlock(c.client.ChainState.ChainID, endBlock)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &endBlock, nil
}

func (c *Commander) syncRange(startBlock, endBlock uint64) error {
	err := c.syncAccounts(startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncBatches(startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Commander) syncBatches(startBlock, endBlock uint64) (err error) {
	transactionExecutor, err := newTransactionExecutor(c.storage, c.client, c.cfg.Rollup)
	if err != nil {
		return err
	}
	defer transactionExecutor.Rollback(&err)

	err = transactionExecutor.SyncBatches(&c.stateMutex, startBlock, endBlock)
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
