package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidDataLength = NewDisputableTransferErrorWithoutProofs("invalid data length")
	ErrTooManyTx         = NewDisputableTransferErrorWithoutProofs("too many transactions in a commitment")
)

func (t *TransactionExecutor) syncTransferCommitment(
	commitment *encoder.DecodedCommitment,
) (models.GenericTransactionArray, error) {
	if len(commitment.Transactions)%encoder.TransferLength != 0 {
		return nil, ErrInvalidDataLength
	}

	transfers, err := encoder.DeserializeTransfers(commitment.Transactions)
	if err != nil {
		return nil, err
	}

	if uint32(len(transfers)) > t.cfg.MaxTxsPerCommitment {
		return nil, ErrTooManyTx
	}

	// TODO check if commitments can have 0 transfers (signature disputes)

	feeReceiver, err := t.getSyncedCommitmentFeeReceiver(commitment)
	if err != nil {
		return nil, err
	}

	appliedTransfers, stateProofs, err := t.ApplyTransfersForSync(transfers, feeReceiver)
	if err != nil {
		return nil, err
	}

	err = t.verifyStateRoot(commitment.StateRoot, stateProofs)
	if err != nil {
		return nil, err
	}

	if !t.cfg.DevMode {
		err = t.verifyTransferSignature(commitment, appliedTransfers)
		if err != nil {
			return nil, err
		}
	}

	return models.TransferArray(appliedTransfers), nil
}

func (t *TransactionExecutor) getSyncedCommitmentFeeReceiver(commitment *encoder.DecodedCommitment) (*FeeReceiver, error) {
	feeReceiverState, err := t.storage.StateTree.Leaf(commitment.FeeReceiver)
	if err != nil {
		return nil, err
	}
	return &FeeReceiver{
		StateID: commitment.FeeReceiver,
		TokenID: feeReceiverState.TokenID,
	}, nil
}

func (t *TransactionExecutor) verifyStateRoot(commitmentPostState common.Hash, proofs []models.StateMerkleProof) error {
	postStateRoot, err := t.storage.StateTree.Root()
	if err != nil {
		return err
	}
	if *postStateRoot != commitmentPostState {
		return NewDisputableTransferError(ErrInvalidCommitmentStateRoot, proofs)
	}
	return nil
}
