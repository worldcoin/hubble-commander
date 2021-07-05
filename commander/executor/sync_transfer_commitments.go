package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrInvalidDataLength = NewDisputableTransferErrorWithoutProofs("invalid data length")
	ErrTooManyTx         = NewDisputableTransferErrorWithoutProofs("too many transactions in a commitment")
)

// TODO-AFS rename file (and C2T as well)
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

	if uint32(len(transfers)) > t.cfg.TxsPerCommitment {
		return nil, ErrTooManyTx
	}

	feeReceiver, err := t.getSyncedCommitmentFeeReceiver(commitment)
	if err != nil {
		return nil, err
	}

	appliedTransfers, err := t.ApplyTransfersForSync(transfers, feeReceiver)
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
	feeReceiverState, err := t.storage.GetStateLeaf(commitment.FeeReceiver)
	if err != nil {
		return nil, err
	}
	return &FeeReceiver{
		StateID: commitment.FeeReceiver,
		TokenID: feeReceiverState.TokenID,
	}, nil
}
