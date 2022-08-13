package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/sirupsen/logrus"
)

func (c *ExecutionContext) RevertBatches(startBatch *models.Batch) error {
	err := c.storage.StateTree.RevertTo(startBatch.PrevStateRoot)
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

	batchIDs := make([]models.Uint256, 0, len(batches))
	for i := range batches {
		batchIDs = append(batchIDs, batches[i].ID)
	}
	err = c.revertCommitments(batches)
	if err != nil {
		return err
	}
	err = c.storage.RemoveCommitmentsByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	logrus.Debugf("Removing %d local batches", len(batches))
	return c.storage.RemoveBatches(batchIDs...)
}

func (c *ExecutionContext) revertCommitments(batches []models.Batch) error {
	txBatchIDs := make([]models.Uint256, 0, len(batches))
	for i := range batches {
		switch batches[i].Type {
		case batchtype.Transfer, batchtype.Create2Transfer, batchtype.MassMigration:
			txBatchIDs = append(txBatchIDs, batches[i].ID)
		case batchtype.Deposit:
			err := c.revertDepositCommitment(batches[i].ID)
			if err != nil {
				return err
			}
		case batchtype.Genesis:
			panic("batch types not supported")
		}
	}
	return c.excludeTransactionsFromCommitment(txBatchIDs...)
}

func (c *ExecutionContext) revertDepositCommitment(batchID models.Uint256) error {
	commitment, err := c.storage.GetCommitment(&models.CommitmentID{
		BatchID:      batchID,
		IndexInBatch: 0,
	})
	if err != nil {
		return err
	}

	depositCommitment := commitment.ToDepositCommitment()
	return c.storage.AddPendingDepositSubtree(&models.PendingDepositSubtree{
		ID:       depositCommitment.SubtreeID,
		Root:     depositCommitment.SubtreeRoot,
		Deposits: depositCommitment.Deposits,
	})
}

func (c *ExecutionContext) excludeTransactionsFromCommitment(batchIDs ...models.Uint256) error {
	if len(batchIDs) == 0 {
		return nil
	}

	logIDs := make([]uint64, len(batchIDs))
	for _, batchID := range batchIDs {
		logIDs = append(logIDs, batchID.Uint64())
	}

	logrus.WithFields(logrus.Fields{
		"hubble.batches": logIDs,
	}).Error("Rolling back batches, transactions not returned to the mempool")

	slots, err := c.storage.GetTransactionIDsByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}

	return c.storage.MarkTransactionsAsPending(slots)
}
