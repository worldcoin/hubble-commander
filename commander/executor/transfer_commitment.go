package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

type PreparedTransfers struct {
	appliedTransfers    []models.Transfer
	newPendingTransfers []models.Transfer
	feeReceiverStateID  *uint32
}

func (t *TransactionExecutor) prepareTransfers(pendingTransfers []models.Transfer) (*PreparedTransfers, error) {
	initialStateRoot, err := t.stateTree.Root()
	if err != nil {
		return nil, err
	}

	preparedTransfers := &PreparedTransfers{
		appliedTransfers: make([]models.Transfer, 0, t.cfg.TxsPerCommitment),
	}

	invalidTransfers := make([]models.Transfer, 0, 1)

	for {
		if len(pendingTransfers) == 0 {
			pendingTransfers, err = t.storage.GetPendingTransfers(pendingTxsCountMultiplier * t.cfg.TxsPerCommitment)
			if err != nil || len(pendingTransfers) == 0 {
				return nil, err
			}
		}

		var transfers *AppliedTransfers

		maxAppliedTransfers := t.cfg.TxsPerCommitment - uint32(len(preparedTransfers.appliedTransfers))
		transfers, err = t.ApplyTransfers(pendingTransfers, maxAppliedTransfers, false)
		if err != nil {
			return nil, err
		}
		if transfers == nil {
			return nil, ErrNotEnoughTransfers
		}

		preparedTransfers.appliedTransfers = append(preparedTransfers.appliedTransfers, transfers.appliedTransfers...)
		invalidTransfers = append(invalidTransfers, transfers.invalidTransfers...)

		if len(preparedTransfers.appliedTransfers) == int(t.cfg.TxsPerCommitment) {
			preparedTransfers.feeReceiverStateID = transfers.feeReceiverStateID
			break
		}

		limit := pendingTxsCountMultiplier*t.cfg.TxsPerCommitment + uint32(len(preparedTransfers.appliedTransfers)+len(invalidTransfers))
		pendingTransfers, err = t.storage.GetPendingTransfers(limit)
		if err != nil {
			return nil, err
		}

		// TODO - instead of doing this use SQL Offset (needs proper mempool)
		pendingTransfers = removeTransfer(pendingTransfers, append(preparedTransfers.appliedTransfers, invalidTransfers...))

		if len(pendingTransfers) == 0 {
			err = t.stateTree.RevertTo(*initialStateRoot)
			return nil, err
		}
	}

	preparedTransfers.newPendingTransfers = removeTransfer(pendingTransfers, append(preparedTransfers.appliedTransfers, invalidTransfers...))

	return preparedTransfers, nil
}

func (t *TransactionExecutor) prepareTransferCommitment(
	preparedTransfers *PreparedTransfers,
	domain *bls.Domain,
) (*models.Commitment, error) {
	serializedTxs, err := encoder.SerializeTransfers(preparedTransfers.appliedTransfers)
	if err != nil {
		return nil, err
	}

	combinedSignature, err := combineTransferSignatures(preparedTransfers.appliedTransfers, domain)
	if err != nil {
		return nil, err
	}

	commitment, err := t.createAndStoreCommitment(
		txtype.Transfer,
		*preparedTransfers.feeReceiverStateID,
		serializedTxs,
		combinedSignature,
	)
	if err != nil {
		return nil, err
	}

	return commitment, nil
}
