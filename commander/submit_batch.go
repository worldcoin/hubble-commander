package commander

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrNotEnoughCommitments = NewRollupError("not enough commitments")
)

func submitBatch(
	ctx context.Context,
	batchType txtype.TransactionType,
	commitments []models.Commitment,
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
) error {
	if len(commitments) < int(cfg.MinCommitmentsPerBatch) {
		return ErrNotEnoughCommitments
	}

	var tx *types.Transaction
	var err error

	select {
	case <-ctx.Done():
		return NewRollupError("commander is no longer an active proposer")
	default:
	}

	if batchType == txtype.Transfer {
		tx, err = client.SubmitTransfersBatch(commitments)
	} else {
		tx, err = client.SubmitCreate2TransfersBatch(commitments)
	}
	if err != nil {
		return err
	}

	newBatch := models.Batch{
		Type:            batchType,
		TransactionHash: tx.Hash(),
	}
	batchID, err := storage.AddBatch(&newBatch)
	if err != nil {
		return err
	}

	err = markCommitmentsAsIncluded(storage, commitments, *batchID)
	if err != nil {
		return err
	}

	return nil
}

func markCommitmentsAsIncluded(storage *st.Storage, commitments []models.Commitment, batchID int32) error {
	for i := range commitments {
		err := storage.MarkCommitmentAsIncluded(commitments[i].ID, batchID)
		if err != nil {
			return err
		}
	}
	return nil
}
