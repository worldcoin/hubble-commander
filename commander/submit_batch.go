package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrNotEnoughCommitments = NewRollupError("not enough commitments")
)

func submitBatch(
	batchType txtype.TransactionType,
	commitments []models.Commitment,
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
) error {
	if len(commitments) < int(cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}

	var batch *models.Batch
	var accountRoot *common.Hash
	var err error

	// TODO check if rollup hasn't been cancelled

	if batchType == txtype.Transfer {
		batch, accountRoot, err = client.SubmitTransfersBatch(commitments)
	} else {
		batch, accountRoot, err = client.SubmitCreate2TransfersBatch(commitments)
	}
	if err != nil {
		return err
	}

	batchID, err := storage.AddBatch(batch)
	if err != nil {
		return err
	}

	err = markCommitmentsAsIncluded(storage, commitments, *batchID, accountRoot)
	if err != nil {
		return err
	}

	log.Printf("Submitted %d commitment(s) on chain. Batch ID: %d. Batch Hash: %v", len(commitments), batch.Number.Uint64(), batch.Hash)

	return nil
}

func markCommitmentsAsIncluded(storage *st.Storage, commitments []models.Commitment, batchID int32, accountRoot *common.Hash) error {
	for i := range commitments {
		err := storage.MarkCommitmentAsIncluded(commitments[i].ID, batchID, accountRoot)
		if err != nil {
			return err
		}
	}
	return nil
}
