package commander

import (
	"context"
	"log"

	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func saveBatch(storage *st.Storage, client *eth.Client) error {
	pendingBatch, err := storage.GetOldestPendingBatch()
	if st.IsNotFoundError(err) {
		return nil
	} else if err != nil {
		return err
	}

	batchTxReceipt, err := client.ChainConnection.GetBackend().TransactionReceipt(context.Background(), *&pendingBatch.TransactionHash)
	if err != nil {
		return err
	}

	if batchTxReceipt != nil {
		if batchTxReceipt.Status == 1 {
			submissionBlock, latestBatchNumber, err := getLatestSubmissionBlockAndBatchNumber(storage, client)
			if err != nil {
				return err
			}
			newBatches, err := client.GetBatches(submissionBlock)
			if err != nil {
				return err
			}

			for i := range newBatches {
				batch := &newBatches[i]
				if batch.Number.Cmp(latestBatchNumber) <= 0 {
					continue
				}
				
				batch.TransactionHash = pendingBatch.TransactionHash
				err := storage.MarkBatchAsSubmitted(&batch.Batch)
				if err != nil {
					return err
				}

				err = storage.UpdateCommitmentsAccountTreeRoot(batch.Batch.TransactionHash, batch.AccountRoot)
				if err != nil {
					return err
				}

				log.Printf("Submitted %d commitment(s) on chain. Batch ID: %d. Batch Hash: %v", len(batch.Commitments), batch.Number.Uint64(), batch.Hash)
			}

		} else {
			// TODO - Have a discussion with the team on how to handle the situation if transaction was mined unsuccessfully
		}
	} else {
		// TODO - Have a discussion with the team on how to handle the situation when the sent transaction is stuck in the mempool
	}

	return nil
}
