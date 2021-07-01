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
	if len(commitment.Transactions)%encoder.Create2TransferLength != 0 {
		return ErrInvalidDataLength
	}

	transactions, err := t.syncCreate2TransferCommitmentsInternal(commitment)
	if err != nil {
		return err
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
	for i := 0; i < transactions.Len(); i++ {
		transactions.At(i).GetBase().IncludedInCommitment = commitmentID
	}

	for i := 0; i < transactions.Len(); i++ {
		hashTransfer, err := encoder.HashGenericTransaction(transactions.At(i))
		if err != nil {
			return err
		}
		transactions.At(i).GetBase().Hash = *hashTransfer
	}

	return t.storage.BatchAddGenericTransaction(transactions)
}

func (t *TransactionExecutor) syncCreate2TransferCommitmentsInternal(commitment *encoder.DecodedCommitment) (models.GenericTransactionArray, error) {
	deserializedTransfers, pubKeyIDs, err := encoder.DeserializeCreate2Transfers(commitment.Transactions)
	if err != nil {
		return nil, err
	}

	if uint32(len(deserializedTransfers)) > t.cfg.TxsPerCommitment {
		return nil, ErrTooManyTx
	}

	feeReceiver, err := t.getSyncedCommitmentFeeReceiver(commitment)
	if err != nil {
		return nil, err
	}

	transfers, err := t.ApplyCreate2TransfersForSync(deserializedTransfers, pubKeyIDs, feeReceiver)
	if err != nil {
		return nil, err
	}

	if len(transfers.invalidTransfers) > 0 {
		return nil, ErrFraudulentTransfer
	}
	if len(transfers.appliedTransfers) != len(deserializedTransfers) {
		return nil, ErrTransfersNotApplied
	}

	err = t.setPublicKeys(transfers.appliedTransfers)
	if err != nil {
		return nil, err
	}
	if !t.cfg.DevMode {
		err = t.verifyCreate2TransferSignature(commitment, transfers.appliedTransfers)
		if err != nil {
			return nil, err
		}
	}

	return models.Create2TransferArray(transfers.appliedTransfers), nil
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
