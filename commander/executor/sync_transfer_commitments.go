package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrTooManyTx         = NewDisputableTransferError("too many transactions in a commitment")
	ErrInvalidDataLength = NewDisputableTransferError("invalid data length")
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

	deserializedTransfers, err := encoder.DeserializeTransfers(commitment.Transactions)
	if err != nil {
		return err
	}

	if uint32(len(deserializedTransfers)) > t.cfg.TxsPerCommitment {
		return ErrTooManyTx
	}

	transfers, err := t.ApplyTransfers(deserializedTransfers, uint32(len(deserializedTransfers)), true)
	if err != nil {
		return err
	}

	if len(transfers.invalidTransfers) > 0 {
		return ErrFraudulentTransfer
	}
	if len(transfers.appliedTransfers) != len(deserializedTransfers) {
		return ErrTransfersNotApplied
	}

	if !t.cfg.DevMode {
		err = t.verifyTransferSignature(commitment, transfers.appliedTransfers)
		if err != nil {
			return err
		}
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
	for i := range transfers.appliedTransfers {
		transfers.appliedTransfers[i].IncludedInCommitment = commitmentID
	}

	for i := range transfers.appliedTransfers {
		hashTransfer, err := encoder.HashTransfer(&transfers.appliedTransfers[i])
		if err != nil {
			return err
		}
		transfers.appliedTransfers[i].Hash = *hashTransfer
	}

	return t.storage.BatchAddTransfer(transfers.appliedTransfers)
}
