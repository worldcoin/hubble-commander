package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/sirupsen/logrus"
)

func (s *Storage) RevertBatches(startBatch *models.Batch) error {
	err := s.StateTree.RevertTo(*startBatch.PrevStateRoot)
	if err != nil {
		return err
	}
	return s.revertBatchesFrom(&startBatch.ID)
}

func (s *Storage) revertBatchesFrom(startBatchID *models.Uint256) error {
	batches, err := s.GetBatchesInRange(startBatchID, nil)
	if err != nil {
		return err
	}

	batchIDs := make([]models.Uint256, 0, len(batches))
	for i := range batches {
		batchIDs = append(batchIDs, batches[i].ID)
	}
	err = s.revertCommitments(batches)
	if err != nil {
		return err
	}
	err = s.RemoveCommitmentsByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}
	logrus.Debugf("Removing %d local batches", len(batches))
	return s.RemoveBatches(batchIDs...)
}

func (s *Storage) revertCommitments(batches []models.Batch) error {
	txBatchIDs := make([]models.Uint256, 0, len(batches))
	for i := range batches {
		switch batches[i].Type {
		case batchtype.Transfer, batchtype.Create2Transfer, batchtype.MassMigration:
			txBatchIDs = append(txBatchIDs, batches[i].ID)
		case batchtype.Deposit:
			err := s.revertDepositCommitment(batches[i].ID)
			if err != nil {
				return err
			}
		case batchtype.Genesis:
			panic("batch types not supported")
		}
	}
	return s.excludeTransactionsFromCommitment(txBatchIDs...)
}

func (s *Storage) revertDepositCommitment(batchID models.Uint256) error {
	commitment, err := s.GetCommitment(&models.CommitmentID{
		BatchID:      batchID,
		IndexInBatch: 0,
	})
	if err != nil {
		return err
	}

	depositCommitment := commitment.ToDepositCommitment()
	return s.AddPendingDepositSubtree(&models.PendingDepositSubtree{
		ID:       depositCommitment.SubtreeID,
		Root:     depositCommitment.SubtreeRoot,
		Deposits: depositCommitment.Deposits,
	})
}

func (s *Storage) excludeTransactionsFromCommitment(batchIDs ...models.Uint256) error {
	if len(batchIDs) == 0 {
		return nil
	}

	logIDs := make([]uint64, len(batchIDs))
	for _, batchID := range batchIDs {
		logIDs = append(logIDs, batchID.Uint64())
	}

	// TODO: confirm this returns transactions from oldest to newest,
	//       even when we are reverting multiple batches
	slots, err := s.GetTransactionIDsByBatchIDs(batchIDs...)
	if err != nil {
		return err
	}

	// TODO: confirm these are recieved in the correct order

	batchedTxs, err := s.DeleteBatchedTxs(slots)
	if err != nil {
		return err
	}
	return s.UnbatchTransactions(batchedTxs)
}
