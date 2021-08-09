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

	if uint32(len(deserializedTransfers)) > t.cfg.MaxTxsPerCommitment {
		return nil, ErrTooManyTxs
	}

	appliedTransfers, stateProofs, err := t.ApplyCreate2TransfersForSync(deserializedTransfers, pubKeyIDs, commitment.FeeReceiver)
	if err != nil {
		return nil, err
	}

	err = t.verifyStateRoot(commitment.StateRoot, stateProofs)
	if err != nil {
		return nil, err
	}

	err = t.setPublicKeys(appliedTransfers, pubKeyIDs)
	if err != nil {
		return nil, err
	}
	if !t.cfg.DisableSignatures {
		err = t.verifyCreate2TransferSignature(commitment, appliedTransfers)
		if err != nil {
			return nil, err
		}
	}

	return models.Create2TransferArray(appliedTransfers), nil
}

func (t *TransactionExecutor) setPublicKeys(transfers []models.Create2Transfer, pubKeyIDs []uint32) error {
	for i := range transfers {
		leaf, err := t.storage.AccountTree.Leaf(pubKeyIDs[i])
		if err != nil {
			return err
		}
		transfers[i].ToPublicKey = leaf.PublicKey
	}
	return nil
}
