package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidDataLength = NewDisputableError(Transition, "invalid data length")
	ErrTooManyTxs        = NewDisputableError(Transition, "too many transactions in a commitment")
)

func (c *ExecutionContext) syncTransferCommitment(
	commitment *encoder.DecodedCommitment,
) (models.GenericTransactionArray, error) {
	if len(commitment.Transactions)%encoder.TransferLength != 0 {
		return nil, ErrInvalidDataLength
	}

	transfers, err := encoder.DeserializeTransfers(commitment.Transactions)
	if err != nil {
		return nil, err
	}

	if uint32(len(transfers)) > c.cfg.MaxTxsPerCommitment {
		return nil, ErrTooManyTxs
	}

	appliedTransfers, stateProofs, err := c.ApplyTransfersForSync(transfers, commitment.FeeReceiver)
	if err != nil {
		return nil, err
	}

	err = c.verifyStateRoot(commitment.StateRoot, stateProofs)
	if err != nil {
		return nil, err
	}

	if !c.cfg.DisableSignatures {
		err = c.verifyTransferSignature(commitment, appliedTransfers)
		if err != nil {
			return nil, err
		}
	}

	return models.TransferArray(appliedTransfers), nil
}

func (c *ExecutionContext) verifyStateRoot(commitmentPostState common.Hash, proofs []models.StateMerkleProof) error {
	postStateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return err
	}
	if *postStateRoot != commitmentPostState {
		return NewDisputableErrorWithProofs(Transition, ErrInvalidCommitmentStateRoot.Error(), proofs)
	}
	return nil
}
