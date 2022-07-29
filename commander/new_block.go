package commander

import (
	"context"
	stdErrors "errors"
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrIncompleteBlockRangeSync = stdErrors.New("syncing of a block range was stopped prematurely")
	ErrRollbackInProgress       = stdErrors.New("rollback is in progress")

	newBlockTracer = otel.Tracer("newBlockLoop")
)

func (c *Commander) newBlockLoop() error {
	latestBlockNumber, err := c.client.Blockchain.GetLatestBlockNumber()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"latestBlockNumber": *latestBlockNumber}).Debug("Starting newBlockLoop")
	c.metrics.LatestBlockNumber.Set(float64(*latestBlockNumber))

	blocks := make(chan *types.Header, 5)
	blocks <- &types.Header{Number: new(big.Int).SetUint64(*latestBlockNumber)}
	subscription, err := c.client.Blockchain.SubscribeNewHead(blocks)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-c.workersContext.Done():
				return
			case <-ticker.C:
				log.Debug("Running GC in background")
			again:
				innerErr := c.storage.TriggerGC()
				if innerErr == nil {
					goto again
				}
				log.Debug("Finished Running GC: ", innerErr)
			}
		}
	}()

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

			err = c.newBlockLoopIteration(currentBlock)
			if errors.Is(err, ErrIncompleteBlockRangeSync) {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}

func (c *Commander) newBlockLoopIteration(currentBlock *types.Header) error {
	err := c.syncAndManageRollbacks()
	if errors.Is(err, ErrRollbackInProgress) {
		return nil
	}
	if errors.Is(err, chain.ErrWaitToBeMinedTimedOut) {
		// Can happen for dispute or keepRollingBack transactions, continue the loop to retry if necessary
		return nil
	}
	if err != nil {
		return err
	}

	err = c.withdrawRemainingStakes(currentBlock.Number.Uint64())
	if err != nil {
		return errors.WithStack(err)
	}

	if c.isMigrating() {
		return c.migrate()
	}

	isProposer, err := c.client.IsActiveProposer()
	if err != nil {
		return errors.WithStack(err)
	}

	c.manageRollupLoop(isProposer)
	return nil
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
	c.metrics.LatestBlockNumber.Set(float64(*latestBlockNumber))
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

	c.metrics.SyncedBlockNumber.Set(float64(endBlock))
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

	ctx, span := newBlockTracer.Start(context.Background(), "syncRange")
	defer span.End()
	span.SetAttributes(
		attribute.Int("hubble.startBlock", int(startBlock)),
		attribute.Int("hubble.endBlock", int(endBlock)),
	)

	err := c.syncAccounts(ctx, startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncTokens(ctx, startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncSpokes(ctx, startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncDeposits(ctx, startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncBatches(ctx, startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.syncStakeWithdrawals(ctx, startBlock, endBlock)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Commander) withdrawRemainingStakes(currentBlock uint64) error {
	stakes, err := c.storage.GetReadyStateWithdrawals(uint32(currentBlock))
	if err != nil {
		return err
	}
	for i := range stakes {
		_, err = c.client.WithdrawStake(&stakes[i].BatchID)
		if err != nil {
			return err
		}
		err = c.storage.RemovePendingStakeWithdrawal(stakes[i].BatchID)
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
