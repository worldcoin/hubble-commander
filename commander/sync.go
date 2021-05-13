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
	// TODO: Start a database transaction

	// TODO: Query batches starting from the submission block of the latest known batch.
	newBatches, err := client.GetBatches()
	if err != nil {
		return err
	}

	latestBatch, err := storage.GetLatestBatch()
	if err != nil {
		return err
	}

	for _, batch := range newBatches {
		if batch.ID.Cmp(&latestBatch.ID) <= 0 {
			continue
		}

		if batch.Type != txtype.Transfer {
			return fmt.Errorf("unsupported batch type for sync: " + txtype.TransactionTypes[batch.Type])
		}

		// Apply batch

		for _, commitment := range batch.Commitments {
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

		err = storage.AddBatch(&models.Batch{
			Hash:              batch.Hash,
			Type:              batch.Type,
			ID:                batch.ID,
			FinalisationBlock: batch.FinalisationBlock,
		})

		log.Printf("synced new batch from chain #%s: %d commitments included", batch.ID.String(), len(batch.Commitments))
	}

	return nil
}
