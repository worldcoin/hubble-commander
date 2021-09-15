package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/sirupsen/logrus"
)

func (c *ExecutionContext) RevertBatches(startBatch *models.Batch) error {
	err := c.storage.StateTree.RevertTo(*startBatch.PrevStateRoot)
	if err != nil {
		return err
	}
	return c.revertBatchesFrom(&startBatch.ID)
}

func (c *ExecutionContext) revertBatchesFrom(startBatchID *models.Uint256) error {
	batches, err := c.storage.GetBatchesInRange(startBatchID, nil)
	if err != nil {
		return err
	}
	numBatches := len(batches)
	batchIDs := make([]models.Uint256, 0, numBatches)
	for i := range batches {
		batchIDs = append(batchIDs, batches[i].ID)
	}
	err = c.excludeTransactionsFromCommitment(batchIDs...)
	if err != nil {
		return err
	}
	err = c.storage.DeleteCommitmentsByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	logrus.Debugf("Removing %d local batches", numBatches)
	return c.storage.DeleteBatches(batchIDs...)
}

func (c *ExecutionContext) excludeTransactionsFromCommitment(batchIDs ...models.Uint256) error {
	hashes, err := c.storage.GetTransactionHashesByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	return c.storage.MarkTransactionsAsPending(hashes)
}
