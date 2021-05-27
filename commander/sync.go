package commander

import (
	"fmt"
	"log"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	ErrFraudulentTransfer  = errors.New("fraudulent transfer encountered when syncing")
	ErrTransfersNotApplied = errors.New("could not apply all transfers from synced batch")
)

func (t *transactionExecutor) SyncBatches() error {
	submissionBlock, latestBatchID, err := getLatestSubmissionBlockAndBatchNumber(t.storage, t.client)
	if err != nil {
		return err
	}

	newBatches, txHashes, err := t.client.GetBatches(submissionBlock)
	if err != nil {
		return err
	}

	for i := range newBatches {
		batch := &newBatches[i]
		if batch.Number.Cmp(latestBatchID) <= 0 {
			continue
		}
		if err := t.syncBatch(batch, txHashes[i]); err != nil {
			return err
		}
	}

	return nil
}

func getLatestSubmissionBlockAndBatchNumber(storage *st.Storage, client *eth.Client) (*uint32, *models.Uint256, error) {
	var submissionBlock uint32
	var latestBatchNumber *models.Uint256

	latestBatch, err := storage.GetLatestBatch()
	if st.IsNotFoundError(err) {
		submissionBlock = 0
		latestBatchNumber = models.NewUint256(0)
	} else if err != nil {
		return nil, nil, err
	} else {
		blocks, err := client.GetBlocksToFinalise()
		if err != nil {
			return nil, nil, err
		}
		submissionBlock = *latestBatch.FinalisationBlock - uint32(*blocks)
		latestBatchNumber = latestBatch.Number
	}

	return &submissionBlock, latestBatchNumber, nil
}

func (t *transactionExecutor) syncBatch(batch *eth.DecodedBatch, txHash common.Hash) error {
	batch.Batch.TransactionHash = txHash
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
