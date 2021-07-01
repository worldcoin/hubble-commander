package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *TransactionExecutor) syncTransferCommitments(commitment *encoder.DecodedCommitment) (models.GenericTransactionArray, error) {
	deserializedTransfers, err := encoder.DeserializeTransfers(commitment.Transactions)
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

	transfers, err := t.ApplyTransfers(deserializedTransfers, uint32(len(deserializedTransfers)), feeReceiver, true)
	if err != nil {
		return nil, err
	}

	if len(transfers.invalidTransfers) > 0 {
		return nil, ErrFraudulentTransfer
	}
	if len(transfers.appliedTransfers) != len(deserializedTransfers) {
		return nil, ErrTransfersNotApplied
	}

	if !t.cfg.DevMode {
		err = t.verifyTransferSignature(commitment, transfers.appliedTransfers)
		if err != nil {
			return nil, err
		}
	}
	return models.TransferArray(transfers.appliedTransfers), nil
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
