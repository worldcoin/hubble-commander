package commander

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
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

	newBatch := models.PendingBatch{
		Type:            batchType,
		TransactionHash: tx.Hash(),
	}
	_, err = storage.AddPendingBatch(&newBatch)
	if err != nil {
		return err
	}

	return nil
}

// TODO - consinder changing the types here
func markCommitmentsAsIncluded(storage *st.Storage, firstCommitmentID int32, numberOfCommitments int, batchID int32, accountRoot *common.Hash) error {
	for i := 0; i < numberOfCommitments; i++ {
		commitmentID := firstCommitmentID + int32(i)
		err := storage.MarkCommitmentAsIncluded(commitmentID, batchID, accountRoot)
		if err != nil {
			return err
		}
	}
	return nil
}
