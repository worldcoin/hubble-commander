package commander

import (
	"context"
	"errors"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	log "github.com/sirupsen/logrus"
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

func (c *Commander) syncRemoteBatch(remoteBatch *eth.DecodedBatch) error {
	var icError *executor.InconsistentBatchError

	err := c.syncOrDisputeRemoteBatch(remoteBatch)
	if errors.As(err, &icError) {
		return c.replaceBatch(icError.LocalBatch, remoteBatch)
	}
	return err
}

func (c *Commander) syncOrDisputeRemoteBatch(remoteBatch *eth.DecodedBatch) error {
	var dcError *executor.DisputableCommitmentError

	err := c.syncBatch(remoteBatch)
	if errors.As(err, &dcError) {
		logFraudulentBatch(remoteBatch, dcError)
		return c.disputeFraudulentBatch(remoteBatch, dcError.CommitmentIndex, dcError.Proofs)
	}
	return err
}

func (c *Commander) syncBatch(remoteBatch *eth.DecodedBatch) error {
	txExecutor, err := executor.NewTransactionExecutor(c.storage.InternalStorage, c.client, c.cfg.Rollup, context.Background())
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

func (c *Commander) replaceBatch(localBatch *models.Batch, remoteBatch *eth.DecodedBatch) error {
	log.WithFields(log.Fields{"batchID": localBatch.ID.String()}).
		Debug("Local batch inconsistent with remote batch, reverting local batch(es)")

	err := c.revertBatches(localBatch)
	if err != nil {
		return err
	}
	return c.syncOrDisputeRemoteBatch(remoteBatch)
}

func (c *Commander) disputeFraudulentBatch(
	remoteBatch *eth.DecodedBatch,
	commitmentIndex int,
	proofs []models.StateMerkleProof,
) error {
	// TODO transaction executor may not be needed here. Revisit this when extracting disputer package.
	txExecutor, err := executor.NewTransactionExecutor(c.storage.InternalStorage, c.client, c.cfg.Rollup, context.Background())
	if err != nil {
		return err
	}
	defer txExecutor.Rollback(&err)

	err = txExecutor.DisputeTransition(remoteBatch, commitmentIndex, proofs)
	if err != nil {
		return err
	}
	return txExecutor.Commit()
}

func (c *Commander) revertBatches(startBatch *models.Batch) error {
	txExecutor, err := executor.NewTransactionExecutor(c.storage.InternalStorage, c.client, c.cfg.Rollup, context.Background())
	if err != nil {
		return err
	}
	defer txExecutor.Rollback(&err)

	err = txExecutor.RevertBatches(startBatch)
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

func logFraudulentBatch(batch *eth.DecodedBatch, err *executor.DisputableCommitmentError) {
	log.WithFields(log.Fields{"batchID": batch.ID.String()}).
		Infof("Found fraudulent batch. Reason: %s", err.Reason)
}
