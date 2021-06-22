package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

type PreparedCreate2Transfers struct {
	appliedTransfers    []models.Create2Transfer
	addedPubKeyIDs      []uint32
	newPendingTransfers []models.Create2Transfer
	feeReceiverStateID  *uint32
}

func (t *TransactionExecutor) prepareCreate2Transfers(pendingTransfers []models.Create2Transfer) (
	*PreparedCreate2Transfers,
	error,
) {
	initialStateRoot, err := t.stateTree.Root()
	if err != nil {
		return nil, err
	}

	preparedTransfers := &PreparedCreate2Transfers{
		appliedTransfers: make([]models.Create2Transfer, 0, t.cfg.TxsPerCommitment),
		addedPubKeyIDs:   make([]uint32, 0, t.cfg.TxsPerCommitment),
	}

	invalidTransfers := make([]models.Create2Transfer, 0, 1)

	for {
		if len(pendingTransfers) == 0 {
			pendingTransfers, err = t.storage.GetPendingCreate2Transfers(pendingTxsCountMultiplier * t.cfg.TxsPerCommitment)
			if err != nil || len(pendingTransfers) == 0 {
				return nil, err
			}
		}

		var transfers *AppliedC2Transfers

		maxAppliedTransfers := t.cfg.TxsPerCommitment - uint32(len(preparedTransfers.appliedTransfers))
		transfers, err = t.ApplyCreate2Transfers(pendingTransfers, maxAppliedTransfers)
		if err != nil {
			return nil, err
		}
		if transfers == nil {
			return nil, ErrNotEnoughC2Transfers
		}

		preparedTransfers.appliedTransfers = append(preparedTransfers.appliedTransfers, transfers.appliedTransfers...)
		invalidTransfers = append(invalidTransfers, transfers.invalidTransfers...)
		preparedTransfers.addedPubKeyIDs = append(preparedTransfers.addedPubKeyIDs, transfers.addedPubKeyIDs...)

		if len(preparedTransfers.appliedTransfers) == int(t.cfg.TxsPerCommitment) {
			preparedTransfers.feeReceiverStateID = transfers.feeReceiverStateID
			break
		}

		limit := pendingTxsCountMultiplier*t.cfg.TxsPerCommitment + uint32(len(preparedTransfers.appliedTransfers)+len(invalidTransfers))
		pendingTransfers, err = t.storage.GetPendingCreate2Transfers(limit)
		if err != nil {
			return nil, err
		}

		// TODO - instead of doing this use SQL Offset (needs proper mempool)
		pendingTransfers = removeCreate2Transfer(pendingTransfers, append(preparedTransfers.appliedTransfers, invalidTransfers...))

		if len(pendingTransfers) == 0 {
			err = t.stateTree.RevertTo(*initialStateRoot)
			return nil, err
		}
	}

	preparedTransfers.newPendingTransfers = removeCreate2Transfer(
		pendingTransfers,
		append(preparedTransfers.appliedTransfers, invalidTransfers...),
	)

	return preparedTransfers, nil
}

func (t *TransactionExecutor) prepareC2TCommitment(
	preparedTransfers *PreparedCreate2Transfers,
	domain *bls.Domain,
) (
	*models.Commitment,
	error,
) {
	serializedTxs, err := encoder.SerializeCreate2Transfers(preparedTransfers.appliedTransfers, preparedTransfers.addedPubKeyIDs)
	if err != nil {
		return nil, err
	}

	combinedSignature, err := combineCreate2TransferSignatures(preparedTransfers.appliedTransfers, domain)
	if err != nil {
		return nil, err
	}

	return t.createAndStoreCommitment(
		txtype.Create2Transfer,
		*preparedTransfers.feeReceiverStateID,
		serializedTxs,
		combinedSignature,
	)
}
