package commander

import (
	log "github.com/sirupsen/logrus"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

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

func logBatchesCount(newRemoteBatches []eth.DecodedBatch) {
	newBatchesCount := len(newRemoteBatches)
	if newBatchesCount > 0 {
		log.Printf("Found %d batch(es)", newBatchesCount)
	}
}
