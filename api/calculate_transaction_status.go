package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func CalculateTransactionStatus(storage *st.Storage, tx *models.Transaction, latestBlockNumber uint32) (*models.TransactionStatus, error) {
	var status models.TransactionStatus

	if tx.IncludedInCommitment == nil {
		status = models.Pending
	} else {
		status = models.Committed
	}

	if tx.ErrorMessage != nil {
		status = models.Error
	}

	if status == models.Committed {
		commitment, err := storage.GetCommitment(*tx.IncludedInCommitment)
		if err != nil {
			return nil, err
		}

		if commitment.IncludedInBatch != nil {
			status = models.InBatch

			batch, err := storage.GetBatch(*commitment.IncludedInBatch)
			if err != nil {
				return nil, err
			}

			if latestBlockNumber >= batch.FinalisationBlock {
				status = models.Finalized
			}
		}
	}

	return &status, nil
}
