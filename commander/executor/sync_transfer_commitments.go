package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrTooManyTx         = NewDisputableTransferError("too many transactions in a commitment", []models.Witness{})
	ErrInvalidDataLength = NewDisputableTransferError("invalid data length", []models.Witness{})
)

func (t *TransactionExecutor) syncTransferCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		err := t.syncTransferCommitment(batch, &batch.Commitments[i])
		if err == ErrInvalidSignature {
			// TODO: dispute fraudulent commitment
			return err
		}
		if IsDisputableTransferError(err) {
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

	transfers, err := encoder.DeserializeTransfers(commitment.Transactions)
	if err != nil {
		return err
	}

	if uint32(len(transfers)) > t.cfg.TxsPerCommitment {
		return ErrTooManyTx // TODO-AFS shouldn't we check using != here ??
	}

	feeReceiver, err := t.getSyncedCommitmentFeeReceiver(commitment)
	if err != nil {
		return err
	}

	appliedTransfers, err := t.ApplyTransfersForSync(transfers, feeReceiver)
	if err != nil {
		return err
	}

	if !t.cfg.DevMode {
		err = t.verifyTransferSignature(commitment, appliedTransfers)
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

	for i := range appliedTransfers {
		transferHash, err := encoder.HashTransfer(&appliedTransfers[i])
		if err != nil {
			return err
		}
		appliedTransfers[i].Hash = *transferHash
		appliedTransfers[i].IncludedInCommitment = commitmentID
	}

	return t.storage.BatchAddTransfer(appliedTransfers)
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
