package commander

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (t *transactionExecutor) syncCreate2TransferCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		if err := syncCreate2TransferCommitment(t.storage, t.cfg, batch, &batch.Commitments[i]); err != nil {
			return err
		}
	}
	return nil
}

func syncCreate2TransferCommitment(
	storage *st.Storage,
	cfg *config.RollupConfig,
	batch *eth.DecodedBatch,
	commitment *encoder.DecodedCommitment,
) error {
	transfers, pubKeyIDs, err := encoder.DeserializeCreate2Transfers(commitment.Transactions)
	if err != nil {
		return err
	}

	appliedTxs, invalidTxs, err := ApplyCreate2TransfersForSync(storage, transfers, pubKeyIDs, cfg)
	if err != nil {
		return err
	}

	if len(invalidTxs) > 0 {
		return ErrFraudulentTransfer
	}

	if len(appliedTxs) != len(transfers) {
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
