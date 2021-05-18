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

	newBatches, err := client.GetBatches() // TODO query batches starting from the submission block of the latest known batch.
	if err != nil {
		return err
	}

	latestBatchID, err := getLatestBatchID(storage)
	if err != nil {
		return err
	}

	for i := range newBatches {
		batch := &newBatches[i]
		if batch.ID.Cmp(latestBatchID) <= 0 {
			continue
		}
		if err := syncBatch(storage, cfg, batch); err != nil {
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

func syncBatch(storage *st.Storage, cfg *config.RollupConfig, batch *eth.DecodedBatch) error {
	if batch.Type != txtype.Transfer {
		return fmt.Errorf("unsupported batch type for sync: %s", batch.Type) // TODO support create2Transfers
	}

	err := storage.AddBatch(&batch.Batch)
	if err != nil {
		return err
	}

	for i := range batch.Commitments {
		commitment := &batch.Commitments[i]
		if err := syncCommitment(storage, cfg, batch, commitment); err != nil {
			return err
		}
	}

	log.Printf("Synced new batch #%s from chain: %d commitments included", batch.ID.String(), len(batch.Commitments))
	return nil
}

func syncCommitment(storage *st.Storage, cfg *config.RollupConfig, batch *eth.DecodedBatch, commitment *encoder.DecodedCommitment) error {
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
	return err
}
