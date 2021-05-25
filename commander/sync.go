package commander

import (
	"fmt"
	"log"

	"github.com/Worldcoin/hubble-commander/config"
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

func SyncBatches(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (err error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return
	}
	defer tx.Rollback(&err)

	err = unsafeSyncBatches(txStorage, client, cfg)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func unsafeSyncBatches(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	submissionBlock, latestBatchID, err := getLatestSubmissionBlockAndBatchID(storage, client)
	if err != nil {
		return err
	}

	newBatches, err := client.GetBatches(submissionBlock)
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

func getLatestSubmissionBlockAndBatchID(storage *st.Storage, client *eth.Client) (*uint32, *models.Uint256, error) {
	var submissionBlock uint32
	var latestBatchID *models.Uint256

	latestBatch, err := storage.GetLatestBatch()
	if st.IsNotFoundError(err) {
		submissionBlock = 0
		latestBatchID = models.NewUint256(0)
	} else if err != nil {
		return nil, nil, err
	} else {
		blocks, err := client.GetBlocksToFinalise()
		if err != nil {
			return nil, nil, err
		}
		submissionBlock = latestBatch.FinalisationBlock - uint32(*blocks)
		latestBatchID = &latestBatch.ID
	}

	return &submissionBlock, latestBatchID, nil
}

func syncBatch(storage *st.Storage, cfg *config.RollupConfig, batch *eth.DecodedBatch) error {
	err := storage.AddBatch(&batch.Batch)
	if err != nil {
		return err
	}

	switch batch.Type {
	case txtype.Transfer:
		err = syncTransferCommitments(storage, cfg, batch)
		if err != nil {
			return err
		}
	case txtype.Create2Transfer:
		err = syncCreate2TransferCommitments(storage, cfg, batch)
		if err != nil {
			return err
		}
	case txtype.MassMigration:
		return fmt.Errorf("unsupported batch type for sync: %s", batch.Type)
	}

	log.Printf("Synced new batch #%s from chain: %d commitments included", batch.ID.String(), len(batch.Commitments))
	return nil
}
