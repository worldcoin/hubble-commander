package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func (t *TransactionExecutor) prepareTransfers(pendingTransfers []models.Transfer) (
	appliedTransfers, newPendingTransfers []models.Transfer,
	feeReceiverStateID *uint32,
	err error,
) {
	initialStateRoot, err := t.stateTree.Root()
	if err != nil {
		return nil, nil, nil, err
	}

	appliedTransfers = make([]models.Transfer, 0, t.cfg.TxsPerCommitment)
	invalidTransfers := make([]models.Transfer, 0, 1)

	for {
		if len(pendingTransfers) == 0 {
			pendingTransfers, err = t.storage.GetPendingTransfers(pendingTxsCountMultiplier * t.cfg.TxsPerCommitment)
			if err != nil || len(pendingTransfers) == 0 {
				return nil, nil, nil, err
			}
		}

		var transfers *AppliedTransfers

		maxAppliedTransfers := t.cfg.TxsPerCommitment - uint32(len(appliedTransfers))
		transfers, err = t.ApplyTransfers(pendingTransfers, maxAppliedTransfers)
		if err != nil {
			return nil, nil, nil, err
		}
		if transfers == nil {
			return nil, nil, nil, ErrNotEnoughTransfers
		}

		appliedTransfers = append(appliedTransfers, transfers.appliedTransfers...)
		invalidTransfers = append(invalidTransfers, transfers.invalidTransfers...)

		if len(appliedTransfers) == int(t.cfg.TxsPerCommitment) {
			feeReceiverStateID = transfers.feeReceiverStateID
			break
		}

		limit := pendingTxsCountMultiplier*t.cfg.TxsPerCommitment + uint32(len(appliedTransfers)+len(invalidTransfers))
		pendingTransfers, err = t.storage.GetPendingTransfers(limit)
		if err != nil {
			return nil, nil, nil, err
		}

		// TODO - instead of doing this use SQL Offset (needs proper mempool)
		pendingTransfers = removeTransfer(pendingTransfers, append(appliedTransfers, invalidTransfers...))

		if len(pendingTransfers) == 0 {
			err = t.stateTree.RevertTo(*initialStateRoot)
			return nil, nil, nil, err
		}
	}

	newPendingTransfers = removeTransfer(pendingTransfers, append(appliedTransfers, invalidTransfers...))

	return appliedTransfers, newPendingTransfers, feeReceiverStateID, nil
}

func (t *TransactionExecutor) prepareTransferCommitment(
	appliedTransfers []models.Transfer,
	feeReceiverStateID uint32,
	domain *bls.Domain,
) (*models.Commitment, error) {
	serializedTxs, err := encoder.SerializeTransfers(appliedTransfers)
	if err != nil {
		return nil, err
	}

	combinedSignature, err := combineTransferSignatures(appliedTransfers, domain)
	if err != nil {
		return nil, err
	}

	commitment, err := t.createAndStoreCommitment(txtype.Transfer, feeReceiverStateID, serializedTxs, combinedSignature)
	if err != nil {
		return nil, err
	}

	return commitment, nil
}
