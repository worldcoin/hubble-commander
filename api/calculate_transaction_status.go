package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func CalculateTransactionStatus(storage *st.Storage, tx *models.Transaction, latestBlockNumber uint32) (*models.TransactionStatus, error) {
	if tx.ErrorMessage != nil {
		return models.Error.Ref(), nil
	}

	if tx.IncludedInCommitment == nil {
		return models.Pending.Ref(), nil
	}

	commitment, err := storage.GetCommitment(*tx.IncludedInCommitment)
	if err != nil {
		return nil, err
	}

	if commitment.IncludedInBatch == nil {
		return models.Committed.Ref(), nil
	}

	batch, err := storage.GetBatch(*commitment.IncludedInBatch)
	if err != nil {
		return nil, err
	}

	if latestBlockNumber < batch.FinalisationBlock {
		return models.InBatch.Ref(), nil
	}
	return models.Finalised.Ref(), nil
}
