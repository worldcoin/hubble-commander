package commander

import (
	"context"
	stdErrors "errors"
	"log"
	"math/big"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

var (
	ErrIncompleteBlockRangeSync = stdErrors.New("syncing of a block range was stopped prematurely")
)

func (c *Commander) newBlockLoop() error {
	latestBlockNumber, err := c.client.ChainConnection.GetLatestBlockNumber()
	if err != nil {
		return err
	}

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
		case <-c.stopChannel:
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

func (c *Commander) syncToLatestBlock() error {
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

func (c *Commander) syncBatches(startBlock, endBlock uint64) error {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	return c.unsafeSyncBatches(startBlock, endBlock)
}

func (c *Commander) unsafeSyncBatches(startBlock, endBlock uint64) error {
	latestBatchID, err := c.getLatestBatchID()
	if err != nil {
		return err
	}

	newRemoteBatches, err := c.client.GetBatches(&bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	})
	if err != nil {
		return err
	}
	logBatchesCount(newRemoteBatches)

	for i := range newRemoteBatches {
		remoteBatch := &newRemoteBatches[i]
		if remoteBatch.ID.Cmp(latestBatchID) <= 0 {
			log.Printf("Batch #%d already synced. Skipping...", remoteBatch.ID.Uint64())
			continue
		}

		err = c.syncRemoteBatch(remoteBatch)
		if err != nil {
			return err
		}

		select {
		case <-c.stopChannel:
			return ErrIncompleteBlockRangeSync
		default:
			continue
		}
	}

	return nil
}

func (c *Commander) syncRemoteBatch(remoteBatch *eth.DecodedBatch) (err error) {
	txExecutor, err := executor.NewTransactionExecutor(c.storage, c.client, c.cfg.Rollup, executor.TransactionExecutorOpts{AssumeNonces: true})
	if err != nil {
		return err
	}
	defer txExecutor.Rollback(&err)

	err = txExecutor.SyncBatch(remoteBatch)
	if err != nil {
		return err
	}
	return txExecutor.Commit()
}

func (c *Commander) getLatestBatchID() (*models.Uint256, error) {
	latestBatch, err := c.storage.GetLatestSubmittedBatch()
	if st.IsNotFoundError(err) {
		return models.NewUint256(0), nil
	} else if err != nil {
		return nil, err
	}
	return &latestBatch.ID, nil
}

func logSyncedBlocks(startBlock, endBlock uint64) {
	if startBlock == endBlock {
		log.Printf("Syncing block %d", startBlock)
	} else {
		log.Printf("Syncing blocks from %d to %d", startBlock, endBlock)
	}
}

func logBatchesCount(newRemoteBatches []eth.DecodedBatch) {
	newBatchesCount := len(newRemoteBatches)
	if newBatchesCount > 0 {
		log.Printf("Found %d batch(es)", newBatchesCount)
	}
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
