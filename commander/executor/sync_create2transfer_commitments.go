package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *TransactionExecutor) syncCreate2TransferCommitment(
	commitment *encoder.DecodedCommitment,
) (models.GenericTransactionArray, error) {
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
