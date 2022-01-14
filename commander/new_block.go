package commander

import (
	"context"
	stdErrors "errors"
	"math/big"

	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	ErrIncompleteBlockRangeSync = stdErrors.New("syncing of a block range was stopped prematurely")
	ErrRollbackInProgress       = stdErrors.New("rollback is in progress")
)

func (c *Commander) newBlockLoop() error {
	latestBlockNumber, err := c.client.Blockchain.GetLatestBlockNumber()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"latestBlockNumber": *latestBlockNumber}).Debug("Starting newBlockLoop")

	blocks := make(chan *types.Header, 5)
	blocks <- &types.Header{Number: new(big.Int).SetUint64(*latestBlockNumber)}
	subscription, err := c.client.Blockchain.SubscribeNewHead(blocks)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	var rollupCancel context.CancelFunc
	for {
		select {
		case <-c.workersContext.Done():
			return nil
		case err = <-subscription.Err():
			return err
		case currentBlock := <-blocks:
			if currentBlock.Number.Uint64() <= uint64(c.storage.GetLatestBlockNumber()) {
				continue
			}

			err = c.syncAndManageRollbacks()
			if errors.Is(err, ErrIncompleteBlockRangeSync) {
				return nil
			}
			if errors.Is(err, ErrRollbackInProgress) {
				continue
			}
			if errors.Is(err, chain.ErrWaitToBeMinedTimedOut) {
				// Can happen for dispute or keepRollingBack transactions, continue the loop to retry if necessary
				continue
			}
			if err != nil {
				return err
			}

			isProposer, err := c.client.IsActiveProposer()
			if err != nil {
				return errors.WithStack(err)
			}

			err = c.withdrawRemainingStakes(currentBlock.Number.Uint64())
			if err != nil {
				return errors.WithStack(err)
			}

			rollupCancel = c.manageRollupLoop(rollupCancel, isProposer)
		}
	}
}

func (c *Commander) syncAndManageRollbacks() error {
	err := c.syncToLatestBlock()
	if err != nil && !errors.Is(err, ErrRollbackInProgress) {
		return err
	}

	return c.keepRollingBackIfNecessary()
}

func (c *Commander) keepRollingBackIfNecessary() (err error) {
	c.invalidBatchID, err = c.client.GetInvalidBatchID()
	if err != nil {
		return err
	}
	if c.invalidBatchID != nil {
		err = c.client.KeepRollingBack()
		if err != nil {
			return err
		}
		return ErrRollbackInProgress
	}
	return nil
}

func (c *Commander) syncToLatestBlock() (err error) {
	latestBlockNumber, err := c.updateLatestBlockNumber()
	if err != nil {
		return err
	}
	syncedBlock, err := c.storage.GetSyncedBlock()
	if err != nil {
		return err
	}

	for *syncedBlock != *latestBlockNumber {
		c.invalidBatchID, err = c.client.GetInvalidBatchID()
		if err != nil {
			return err
		}

		syncedBlock, err = c.syncForward(*latestBlockNumber)
		if err != nil {
			return err
		}

		latestBlockNumber, err = c.updateLatestBlockNumber()
		if err != nil {
			return err
		}

		select {
		case <-c.workersContext.Done():
			return nil
		default:
			continue
		}
	}
	return nil
}

func (c *Commander) updateLatestBlockNumber() (*uint64, error) {
	latestBlockNumber, err := c.client.Blockchain.GetLatestBlockNumber()
	if err != nil {
		return nil, err
	}
	c.storage.SetLatestBlockNumber(uint32(*latestBlockNumber))
	return latestBlockNumber, nil
}

func (c *Commander) syncForward(latestBlockNumber uint64) (*uint64, error) {
	syncedBlock, err := c.storage.GetSyncedBlock()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	startBlock := *syncedBlock + 1
	endBlock := min(latestBlockNumber, startBlock+uint64(c.cfg.Rollup.SyncSize))

	duration, err := metrics.MeasureDuration(func() error {
		return c.syncRange(startBlock, endBlock)
	})
	if err != nil {
		return nil, err
	}

	err = c.storage.SetSyncedBlock(endBlock)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	metrics.SaveHistogramMeasurement(duration, c.metrics.SyncingMethodDuration, prometheus.Labels{
		"method": metrics.SyncRangeMethod,
	})

	return &endBlock, nil
}

func (c *Commander) syncRange(startBlock, endBlock uint64) error {
	logSyncedBlocks(startBlock, endBlock)

	err := c.syncAccounts(startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncTokens(startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncSpokes(startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncDeposits(startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncBatches(startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncStakeWithdrawals(startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Commander) withdrawRemainingStakes(finalizationBlock uint64) error {
	stakes, err := c.storage.GetReadyStateWithdrawals(uint32(finalizationBlock))
	if err != nil {
		return err
	}
	for _, stake := range stakes {
		err = c.client.WithdrawStake(&stake.BatchID)
		if err != nil {
			return err
		}
		err = c.storage.RemovePendingStakeWithdrawal(stake.BatchID)
		if err != nil {
			return err
		}
	}
	return nil
}

func logSyncedBlocks(startBlock, endBlock uint64) {
	if startBlock == endBlock {
		log.Printf("Syncing block %d", startBlock)
	} else {
		log.Printf("Syncing blocks from %d to %d", startBlock, endBlock)
	}
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
