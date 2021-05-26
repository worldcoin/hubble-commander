package commander

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (t *transactionExecutor) syncTransferCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		if err := syncTransferCommitment(t.storage, t.cfg, batch, &batch.Commitments[i]); err != nil {
			return err
		}
	}
	return nil
}

func syncTransferCommitment(
	storage *st.Storage,
	cfg *config.RollupConfig,
	batch *eth.DecodedBatch,
	commitment *encoder.DecodedCommitment,
) error {
	transfers, err := encoder.DeserializeTransfers(commitment.Transactions)
	if err != nil {
		return err
	}

	appliedTransfers, invalidTransfers, _, err := ApplyTransfers(storage, transfers, cfg)
	if err != nil {
		return err
	}

	if len(invalidTransfers) > 0 {
		return ErrFraudulentTransfer
	}

	if len(appliedTransfers) != len(transfers) {
		return ErrTransfersNotApplied
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
