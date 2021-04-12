package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

func submitBatch(commitments []models.Commitment, storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (err error) {
	if len(commitments) < int(cfg.MinCommitmentsPerBatch) {
		return
	}

	batch, accountRoot, err := client.SubmitTransfersBatch(commitments)
	if err != nil {
		return
	}

	err = storage.AddBatch(batch)
	if err != nil {
		return
	}

	err = markCommitmentsAsIncluded(storage, commitments, &batch.Hash, accountRoot)
	if err != nil {
		return
	}

	log.Printf("Submitted %d commitment(s) on chain. Batch ID: %d. Batch Hash: %v", len(commitments), batch.ID.Uint64(), batch.Hash)

	return nil
}

func markCommitmentsAsIncluded(storage *st.Storage, commitments []models.Commitment, batchHash, accountRoot *common.Hash) error {
	for i := range commitments {
		err := storage.MarkCommitmentAsIncluded(commitments[i].ID, batchHash, accountRoot)
		if err != nil {
			return err
		}
	}
	return nil
}
