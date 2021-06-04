package commander

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *transactionExecutor) syncCreate2TransferCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		if err := t.syncCreate2TransferCommitment(batch, &batch.Commitments[i]); err != nil {
			return err
		}
	}
	return nil
}

func (t *transactionExecutor) syncCreate2TransferCommitment(
	batch *eth.DecodedBatch,
	commitment *encoder.DecodedCommitment,
) error {
	deserializedTransfers, pubKeyIDs, err := encoder.DeserializeCreate2Transfers(commitment.Transactions)
	if err != nil {
		return err
	}

	transfers, err := t.ApplyCreate2TransfersForSync(deserializedTransfers, pubKeyIDs)
	if err != nil {
		return err
	}

	if len(transfers.invalidTransfers) > 0 {
		return ErrFraudulentTransfer
	}

	if len(transfers.appliedTransfers) != len(deserializedTransfers) {
		return ErrTransfersNotApplied
	}

	_, err = t.storage.AddCommitment(&models.Commitment{
		Type:              batch.Type,
		Transactions:      commitment.Transactions,
		FeeReceiver:       commitment.FeeReceiver,
		CombinedSignature: commitment.CombinedSignature,
		PostStateRoot:     commitment.StateRoot,
		AccountTreeRoot:   &batch.AccountRoot,
		IncludedInBatch:   &batch.ID,
	})
	return err
}
