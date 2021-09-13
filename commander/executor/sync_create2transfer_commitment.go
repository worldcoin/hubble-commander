package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
)

func (c *ExecutionContext) syncCreate2TransferCommitment(
	commitment *encoder.DecodedCommitment,
) (models.GenericTransactionArray, error) {
	deserializedTransfers, pubKeyIDs, err := encoder.DeserializeCreate2Transfers(commitment.Transactions)
	if err != nil {
		return nil, err
	}

	if uint32(len(deserializedTransfers)) > c.cfg.MaxTxsPerCommitment {
		return nil, ErrTooManyTxs
	}

	appliedTransfers, stateProofs, err := c.ApplyCreate2TransfersForSync(deserializedTransfers, pubKeyIDs, commitment.FeeReceiver)
	if err != nil {
		return nil, err
	}

	err = c.verifyStateRoot(commitment.StateRoot, stateProofs)
	if err != nil {
		return nil, err
	}

	err = c.setPublicKeys(appliedTransfers, pubKeyIDs)
	if err != nil {
		return nil, err
	}
	if !c.cfg.DisableSignatures {
		err = c.verifyCreate2TransferSignature(commitment, appliedTransfers)
		if err != nil {
			return nil, err
		}
	}

	return models.Create2TransferArray(appliedTransfers), nil
}

func (c *ExecutionContext) setPublicKeys(transfers []models.Create2Transfer, pubKeyIDs []uint32) error {
	for i := range transfers {
		leaf, err := c.storage.AccountTree.Leaf(pubKeyIDs[i])
		if err != nil {
			return err
		}
		transfers[i].ToPublicKey = leaf.PublicKey
	}
	return nil
}
