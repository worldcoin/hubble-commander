package commander

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *transactionExecutor) syncCreate2TransferCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		err := t.syncCreate2TransferCommitment(batch, &batch.Commitments[i])
		if err != nil {
			if err == ErrInvalidSignature { //nolint: staticcheck
				//TODO: dispute fraudulent commitment
			}
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

	isValid, err := t.verifyCreate2TransferSignature(commitment, transfers.appliedTransfers)
	if err != nil {
		return err
	}
	if !isValid {
		return ErrInvalidSignature
	}

	commitmentID, err := t.storage.AddCommitment(&models.Commitment{
		Type:              batch.Type,
		Transactions:      commitment.Transactions,
		FeeReceiver:       commitment.FeeReceiver,
		CombinedSignature: commitment.CombinedSignature,
		PostStateRoot:     commitment.StateRoot,
		IncludedInBatch:   &batch.ID,
	})
	if err != nil {
		return err
	}
	for i := range transfers.appliedTransfers {
		transfers.appliedTransfers[i].IncludedInCommitment = commitmentID
	}

	for i := range transfers.appliedTransfers {
		hashTransfer, err := encoder.HashCreate2Transfer(&transfers.appliedTransfers[i])
		if err != nil {
			return err
		}
		transfers.appliedTransfers[i].Hash = *hashTransfer
	}

	return t.storage.BatchAddCreate2Transfer(transfers.appliedTransfers)
}
