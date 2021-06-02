package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrNotEnoughCommitments = NewRollupError("not enough commitments")
)

func (t *transactionExecutor) submitBatch(
	batchType txtype.TransactionType,
	commitments []models.Commitment,
) error {
	if len(commitments) < int(t.cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}

	var batch *models.Batch
	var accountRoot *common.Hash
	var err error

	if batchType == txtype.Transfer {
		batch, accountRoot, err = t.client.SubmitTransfersBatch(commitments)
	} else {
		batch, accountRoot, err = t.client.SubmitCreate2TransfersBatch(commitments)
	}
	if err != nil {
		return err
	}

	batchID, err := t.storage.AddBatch(batch)
	if err != nil {
		return err
	}

	err = t.markCommitmentsAsIncluded(commitments, *batchID, accountRoot)
	if err != nil {
		return err
	}

	log.Printf("Submitted %d commitment(s) on chain. Batch ID: %d. Batch Hash: %v", len(commitments), batch.Number.Uint64(), batch.Hash)

	return nil
}

func (t *transactionExecutor) markCommitmentsAsIncluded(commitments []models.Commitment, batchID int32, accountRoot *common.Hash) error {
	for i := range commitments {
		err := t.storage.MarkCommitmentAsIncluded(commitments[i].ID, batchID, accountRoot)
		if err != nil {
			return err
		}
	}
	return nil
}
