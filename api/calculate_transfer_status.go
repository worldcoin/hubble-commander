package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func CalculateTransferStatus(storage *st.Storage, transfer *models.Transfer, latestBlockNumber uint32) (*models.TransactionStatus, error) {
	if transfer.ErrorMessage != nil {
		return models.Error.Ref(), nil
	}

	if transfer.IncludedInCommitment == nil {
		return models.Pending.Ref(), nil
	}

	batch, err := storage.GetBatchByCommitmentID(*transfer.IncludedInCommitment)
	if err != nil {
		return nil, err
	}

	if latestBlockNumber < batch.FinalisationBlock {
		return models.InBatch.Ref(), nil
	}
	return models.Finalised.Ref(), nil
}
