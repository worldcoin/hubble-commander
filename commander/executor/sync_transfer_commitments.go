package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrTooManyTx         = NewDisputableTransferError(TransitionError, "too many transactions in a commitment")
	ErrInvalidDataLength = NewDisputableTransferError(TransitionError, "invalid data length")
)

func (t *TransactionExecutor) syncTransferCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		err := t.syncTransferCommitment(batch, &batch.Commitments[i])
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

func (t *TransactionExecutor) syncTransferCommitment(
	batch *eth.DecodedBatch,
	commitment *encoder.DecodedCommitment,
) error {
	if len(commitment.Transactions)%encoder.TransferLength != 0 {
		return ErrInvalidDataLength
	}

	transactions, err := t.syncTransferCommitmentsInternal(commitment)
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

func (t *TransactionExecutor) syncTransferCommitmentsInternal(commitment *encoder.DecodedCommitment) (models.GenericTransactionArray, error) {
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
