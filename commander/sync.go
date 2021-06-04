package commander

import (
	"fmt"
	"log"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

var (
	ErrFraudulentTransfer  = errors.New("fraudulent transfer encountered when syncing")
	ErrTransfersNotApplied = errors.New("could not apply all transfers from synced batch")
)

func (t *transactionExecutor) SyncBatches(startBlock, endBlock uint64) error {
	latestBatchNumber, err := getLatestBatchID(t.storage)
	if err != nil {
		return err
	}

	newBatches, err := t.client.GetBatches(&bind.FilterOpts{
		Start: startBlock,
		End:   &endBlock,
	})
	if err != nil {
		return err
	}

	for i := range newBatches {
		batch := &newBatches[i]
		if batch.Number.Cmp(latestBatchNumber) <= 0 {
			continue
		}
		if err := t.syncBatch(batch); err != nil {
			return err
		}
	}

	return nil
}

func getLatestBatchID(storage *st.Storage) (*models.Uint256, error) {
	latestBatch, err := storage.GetLatestSubmittedBatch()
	if st.IsNotFoundError(err) {
		return models.NewUint256(0), nil
	} else if err != nil {
		return nil, err
	}
	return latestBatch.Number, nil
}

func (t *transactionExecutor) syncBatch(batch *eth.DecodedBatch) error {
	batchID, err := t.storage.AddBatch(&batch.Batch)
	if err != nil {
		return err
	}

	batch.Batch.ID = *batchID

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

	log.Printf("Synced new batch #%s from chain: %d commitments included", batch.Number.String(), len(batch.Commitments))
	return nil
}
