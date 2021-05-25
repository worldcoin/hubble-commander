package commander

import (
	"fmt"
	"log"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var (
	ErrFraudulentTransfer  = errors.New("fraudulent transfer encountered when syncing")
	ErrTransfersNotApplied = errors.New("could not apply all transfers from synced batch")
)

func (t *transactionExecutor) SyncBatches() error {
	newBatches, err := t.client.GetBatches() // TODO query batches starting from the submission block of the latest known batch.
	if err != nil {
		return err
	}

	latestBatchID, err := getLatestBatchID(t.storage)
	if err != nil {
		return err
	}

	for i := range newBatches {
		batch := &newBatches[i]
		if batch.ID.Cmp(latestBatchID) <= 0 {
			continue
		}
		if err := t.syncBatch(batch); err != nil {
			return err
		}
	}

	return nil
}

func getLatestBatchID(storage *st.Storage) (*models.Uint256, error) {
	var latestBatchID models.Uint256
	latestBatch, err := storage.GetLatestBatch()
	if st.IsNotFoundError(err) {
		latestBatchID = models.MakeUint256(0)
	} else if err != nil {
		return nil, err
	} else {
		latestBatchID = latestBatch.ID
	}
	return &latestBatchID, nil
}

func (t *transactionExecutor) syncBatch(batch *eth.DecodedBatch) error {
	err := t.storage.AddBatch(&batch.Batch)
	if err != nil {
		return err
	}

	switch batch.Type {
	case txtype.Transfer:
		err = t.syncTransferCommitments(batch)
		if err != nil {
			return err
		}
	case txtype.Create2Transfer:
		err = t.syncCreate2TransferCommitments(batch)
		if err != nil {
			return err
		}
	case txtype.MassMigration:
		return fmt.Errorf("unsupported batch type for sync: %s", batch.Type)
	}

	log.Printf("Synced new batch #%s from chain: %d commitments included", batch.ID.String(), len(batch.Commitments))
	return nil
}
