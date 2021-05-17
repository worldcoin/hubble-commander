package commander

import (
	"fmt"
	"log"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func SyncBatches(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	// TODO start a database transaction

	// TODO query batches starting from the submission block of the latest known batch.
	newBatches, err := client.GetBatches()
	if err != nil {
		return err
	}

	var latestBatchID models.Uint256
	latestBatch, err := storage.GetLatestBatch()
	if st.IsNotFoundError(err) {
		latestBatchID = models.MakeUint256(0)
	} else if err != nil {
		return err
	} else {
		latestBatchID = latestBatch.ID
	}

	for i := range newBatches {
		batch := &newBatches[i]
		if batch.ID.Cmp(&latestBatchID) <= 0 {
			continue
		}

		if batch.Type != txtype.Transfer {
			// TODO support create2Transfers
			return fmt.Errorf("unsupported batch type for sync: %s", batch.Type)
		}

		// Apply batch

		err = storage.AddBatch(&batch.Batch)
		if err != nil {
			return err
		}

		for i := range batch.Commitments {
			commitment := &batch.Commitments[i]
			transfers, err := encoder.DeserializeTransfers(commitment.Transactions)
			if err != nil {
				return err
			}

			appliedTransfers, invalidTransfers, _, err := ApplyTransfers(storage, transfers, cfg)
			if err != nil {
				return err
			}

			if len(invalidTransfers) > 0 {
				return fmt.Errorf("fraduelent transfer encountered when syncing")
			}

			if len(appliedTransfers) != len(transfers) {
				return fmt.Errorf("could not apply all transfers from synced batch")
			}

			_, err = storage.AddCommitment(&models.Commitment{
				Type:              batch.Type,
				Transactions:      commitment.Transactions,
				FeeReceiver:       commitment.FeeReceiver,
				CombinedSignature: commitment.CombinedSignature,
				PostStateRoot:     commitment.StateRoot,
				AccountTreeRoot:   &batch.AccountRoot,
				IncludedInBatch:   &batch.Hash,
			})
			if err != nil {
				return err
			}
		}

		log.Printf("Synced new batch #%s from chain: %d commitments included", batch.ID.String(), len(batch.Commitments))
	}

	return nil
}
