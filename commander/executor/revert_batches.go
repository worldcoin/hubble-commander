package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/sirupsen/logrus"
)

func (t *ExecutionContext) RevertBatches(startBatch *models.Batch) error {
	err := t.storage.StateTree.RevertTo(*startBatch.PrevStateRoot)
	if err != nil {
		return err
	}
	return t.revertBatchesFrom(&startBatch.ID)
}

func (t *ExecutionContext) revertBatchesFrom(startBatchID *models.Uint256) error {
	batches, err := t.storage.GetBatchesInRange(startBatchID, nil)
	if err != nil {
		return err
	}
	numBatches := len(batches)
	batchIDs := make([]models.Uint256, 0, numBatches)
	for i := range batches {
		batchIDs = append(batchIDs, batches[i].ID)
	}
	err = t.excludeTransactionsFromCommitment(batchIDs...)
	if err != nil {
		return err
	}
	err = t.storage.DeleteCommitmentsByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	logrus.Debugf("Removing %d local batches", numBatches)
	return t.storage.DeleteBatches(batchIDs...)
}

func (t *ExecutionContext) excludeTransactionsFromCommitment(batchIDs ...models.Uint256) error {
	hashes, err := t.storage.GetTransactionHashesByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	return t.storage.MarkTransactionsAsPending(hashes)
}
