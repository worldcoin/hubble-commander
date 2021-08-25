package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func CalculateTransactionStatus(
	storage *st.Storage,
	transfer *models.TransactionBase,
	latestBlockNumber uint32,
) (*txstatus.TransactionStatus, error) {
	if transfer.ErrorMessage != nil {
		return txstatus.Error.Ref(), nil
	}

	if transfer.BatchID == nil {
		return txstatus.Pending.Ref(), nil
	}

	batch, err := storage.GetBatch(*transfer.BatchID)
	if err != nil {
		return nil, err
	}

	if batch.FinalisationBlock == nil {
		return txstatus.Pending.Ref(), nil
	}

	return calculateFinalisedStatus(latestBlockNumber, *batch.FinalisationBlock), nil
}

func calculateFinalisedStatus(latestBlockNumber, batchFinalisationBlock uint32) *txstatus.TransactionStatus {
	if latestBlockNumber < batchFinalisationBlock {
		return txstatus.InBatch.Ref()
	}
	return txstatus.Finalised.Ref()
}
