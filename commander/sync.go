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
	latestBatchNumber, err := getLatestBatchNumber(t.storage)
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

		localBatch, err := t.storage.GetBatchByNumber(*batch.Number)
		if st.IsNotFoundError(err) {
			err = t.syncBatch(batch)
			if err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}

		err = t.syncExistingBatch(batch, localBatch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *transactionExecutor) syncExistingBatch(batch *eth.DecodedBatch, localBatch *models.Batch) error {
	if batch.TransactionHash == localBatch.TransactionHash {
		batch.ID = localBatch.ID
		err := t.storage.MarkBatchAsSubmitted(&batch.Batch)
		if err != nil {
			return err
		}

		err = t.storage.UpdateCommitmentsAccountTreeRoot(localBatch.ID, batch.AccountRoot)
		if err != nil {
			return err
		}

		log.Printf(
			"Submitted %d commitment(s) on chain. Batch number: %d. Batch Hash: %v",
			len(batch.Commitments),
			batch.Number.Uint64(),
			batch.Hash,
		)
	} else { // nolint:staticcheck
		// TODO: handle race condition
	}
	return nil
}

func getLatestBatchNumber(storage *st.Storage) (*models.Uint256, error) {
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
