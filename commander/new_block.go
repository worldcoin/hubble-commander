package commander

import (
	"context"
	stdErrors "errors"
	"math/big"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrIncompleteBlockRangeSync = stdErrors.New("syncing of a block range was stopped prematurely")
	ErrRollbackInProgress       = stdErrors.New("rollback is in progress")
)

func (c *Commander) newBlockLoop() error {
	latestBlockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"latestBlockNumber": *latestBlockNumber}).Debug("Starting newBlockLoop")

	blocks := make(chan *types.Header, 5)
	blocks <- &types.Header{Number: new(big.Int).SetUint64(*latestBlockNumber)}
	subscription, err := c.client.ChainConnection.SubscribeNewHead(blocks)
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

			err = c.syncToLatestBlock()
			if errors.Is(err, ErrIncompleteBlockRangeSync) {
				return nil
			}
			if errors.Is(err, ErrRollbackInProgress) {
				continue
			}
			if err != nil {
				return err
			}
			// TODO handle ErrRollbackInProgressHere and `continue` the loop in this case

			isProposer, err := c.client.IsActiveProposer()
			if err != nil {
				return errors.WithStack(err)
			}
			rollupCancel = c.manageRollupLoop(rollupCancel, isProposer)
		}
	}
}

func (c *Commander) syncToLatestBlock() error {
	latestBlockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
	if err != nil {
		return err
	}
	c.storage.SetLatestBlockNumber(uint32(*latestBlockNumber))
	// TODO move invalidBatchID setting here

	syncedBlock := ref.Uint64(uint64(0))
	for *syncedBlock != *latestBlockNumber {
		syncedBlock, err = c.syncForward(*latestBlockNumber)
		if err != nil {
			return err
		}
		// TODO if err == ErrRollbackInProgress break the loop

		latestBlockNumber, err = c.client.ChainConnection.GetLatestBlockNumber()
		if err != nil {
			return err
		}
		c.storage.SetLatestBlockNumber(uint32(*latestBlockNumber))
		// TODO update invalidBatchID here

		select {
		case <-c.workersContext.Done():
			return nil
		default:
			continue
		}
	}

	// TODO update invalidBatchID here
	// TODO if invalidBatchID != nil call keepRollingBack and return ErrRollbackInProgress
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
	logSyncedBlocks(startBlock, endBlock)

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
