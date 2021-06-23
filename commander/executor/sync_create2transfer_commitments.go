package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *TransactionExecutor) syncCreate2TransferCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		err := t.syncCreate2TransferCommitment(batch, &batch.Commitments[i])
		if err == ErrInvalidSignature {
			// TODO: dispute fraudulent commitment
			return err
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TransactionExecutor) syncCreate2TransferCommitment(
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

	err = t.setPublicKeys(transfers.appliedTransfers)
	if err != nil {
		return err
	}
	if !t.cfg.DevMode {
		err = t.verifyCreate2TransferSignature(commitment, transfers.appliedTransfers)
		if err != nil {
			return err
		}
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

func (t *TransactionExecutor) setPublicKeys(transfers []models.Create2Transfer) error {
	for i := range transfers {
		publicKey, err := t.storage.GetPublicKeyByStateID(*transfers[i].ToStateID)
		if err != nil {
			return err
		}
		transfers[i].ToPublicKey = *publicKey
	}
	return nil
}
