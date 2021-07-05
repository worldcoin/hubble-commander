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

	appliedTransfers, err := t.ApplyCreate2TransfersForSync(deserializedTransfers, pubKeyIDs, feeReceiver)
	if err != nil {
		return nil, err
	}

	err = t.setPublicKeys(appliedTransfers)
	if err != nil {
		return nil, err
	}
	if !t.cfg.DevMode {
		err = t.verifyCreate2TransferSignature(commitment, appliedTransfers)
		if err != nil {
			return nil, err
		}
	}

	return models.Create2TransferArray(appliedTransfers), nil
}

// TODO-AFS consider using existing pubKeyIDs slice
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
