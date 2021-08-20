package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotEnoughC2Transfers = NewRollupError("not enough create2transfers")
)

func (t *TransactionExecutor) CreateCreate2TransferCommitments(domain *bls.Domain) (commitments []models.Commitment, err error) {
	pendingTransfers, err := t.queryPendingC2Ts()
	if err != nil {
		return nil, err
	}

	commitmentKey, err := t.createCommitmentKey()
	if err != nil {
		return nil, err
	}
	commitments = make([]models.Commitment, 0, t.cfg.MaxCommitmentsPerBatch)

	for i := uint32(0); len(commitments) != int(t.cfg.MaxCommitmentsPerBatch); i++ {
		var commitment *models.Commitment
		commitmentKey.IndexInBatch = i

		pendingTransfers, commitment, err = t.createC2TCommitment(pendingTransfers, commitmentKey, domain)
		if err == ErrNotEnoughC2Transfers {
			break
		}
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)
	}

	if len(commitments) == 0 {
		return nil, ErrNotEnoughC2Transfers
	}

	return commitments, nil
}

func (t *TransactionExecutor) createC2TCommitment(
	pendingTransfers []models.Create2Transfer,
	commitmentKey *models.CommitmentKey,
	domain *bls.Domain,
) (newPendingTransfers []models.Create2Transfer, commitment *models.Commitment, err error) {
	startTime := time.Now()

	pendingTransfers, err = t.refillPendingC2Ts(pendingTransfers)
	if err != nil {
		return nil, nil, err
	}

	feeReceiver, err := t.getCommitmentFeeReceiver()
	if err != nil {
		return nil, nil, err
	}

	initialStateRoot, err := t.storage.StateTree.Root()
	if err != nil {
		return nil, nil, err
	}

	appliedTransfers, newPendingTransfers, addedPubKeyIDs, err := t.applyC2TsForCommitment(pendingTransfers, feeReceiver)
	if err == ErrNotEnoughC2Transfers {
		if revertErr := t.storage.StateTree.RevertTo(*initialStateRoot); revertErr != nil {
			return nil, nil, revertErr
		}
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}

	commitment, err = t.buildC2TCommitment(appliedTransfers, addedPubKeyIDs, commitmentKey, feeReceiver.StateID, domain)
	if err != nil {
		return nil, nil, err
	}

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		txtype.Create2Transfer,
		len(appliedTransfers),
		time.Since(startTime).Round(time.Millisecond).String(),
	)

	return newPendingTransfers, commitment, nil
}

func (t *TransactionExecutor) applyC2TsForCommitment(pendingTransfers []models.Create2Transfer, feeReceiver *FeeReceiver) (
	appliedTransfers, newPendingTransfers []models.Create2Transfer,
	addedPubKeyIDs []uint32,
	err error,
) {
	appliedTransfers = make([]models.Create2Transfer, 0, t.cfg.MaxTxsPerCommitment)
	invalidTransfers := make([]models.Create2Transfer, 0, 1)
	addedPubKeyIDs = make([]uint32, 0, t.cfg.MaxTxsPerCommitment)

	for {
		var transfers *AppliedC2Transfers

		numNeededTransfers := t.cfg.MaxTxsPerCommitment - uint32(len(appliedTransfers))
		transfers, err = t.ApplyCreate2Transfers(pendingTransfers, numNeededTransfers, feeReceiver)
		if err != nil {
			return nil, nil, nil, err
		}

		appliedTransfers = append(appliedTransfers, transfers.appliedTransfers...)
		invalidTransfers = append(invalidTransfers, transfers.invalidTransfers...)
		addedPubKeyIDs = append(addedPubKeyIDs, transfers.addedPubKeyIDs...)

		if len(appliedTransfers) == int(t.cfg.MaxTxsPerCommitment) {
			newPendingTransfers = removeC2Ts(pendingTransfers, append(appliedTransfers, invalidTransfers...))
			return appliedTransfers, newPendingTransfers, addedPubKeyIDs, nil
		}

		morePendingTransfers, err := t.queryMorePendingC2Ts(appliedTransfers)
		if err == ErrNotEnoughC2Transfers {
			newPendingTransfers = removeC2Ts(pendingTransfers, append(appliedTransfers, invalidTransfers...))
			return appliedTransfers, newPendingTransfers, addedPubKeyIDs, nil
		}
		if err != nil {
			return nil, nil, nil, err
		}
		pendingTransfers = morePendingTransfers
	}
}

func (t *TransactionExecutor) refillPendingC2Ts(pendingTransfers []models.Create2Transfer) ([]models.Create2Transfer, error) {
	if len(pendingTransfers) < int(t.cfg.MaxTxsPerCommitment) {
		return t.queryPendingC2Ts()
	}
	return pendingTransfers, nil
}

func (t *TransactionExecutor) queryPendingC2Ts() ([]models.Create2Transfer, error) {
	pendingTransfers, err := t.storage.GetPendingCreate2Transfers(t.cfg.MaxCommitmentsPerBatch * t.cfg.MaxTxsPerCommitment)
	if err != nil {
		return nil, err
	}
	if len(pendingTransfers) < int(t.cfg.MinTxsPerCommitment) {
		return nil, ErrNotEnoughC2Transfers
	}
	return pendingTransfers, nil
}

func (t *TransactionExecutor) queryMorePendingC2Ts(appliedTransfers []models.Create2Transfer) ([]models.Create2Transfer, error) {
	numAppliedTransfers := uint32(len(appliedTransfers))
	// TODO use SQL Offset instead
	pendingTransfers, err := t.storage.GetPendingCreate2Transfers(
		t.cfg.MaxCommitmentsPerBatch*t.cfg.MaxTxsPerCommitment + numAppliedTransfers,
	)
	if err != nil {
		return nil, err
	}
	pendingTransfers = removeC2Ts(pendingTransfers, appliedTransfers)

	numNeededTransfers := t.cfg.MaxTxsPerCommitment - numAppliedTransfers
	if len(pendingTransfers) < int(numNeededTransfers) {
		return nil, ErrNotEnoughC2Transfers
	}
	return pendingTransfers, nil
}

func removeC2Ts(transferList, toRemove []models.Create2Transfer) []models.Create2Transfer {
	outputIndex := 0
	for i := range transferList {
		transfer := &transferList[i]
		if !create2TransferExists(toRemove, transfer) {
			transferList[outputIndex] = *transfer
			outputIndex++
		}
	}

	return transferList[:outputIndex]
}

func create2TransferExists(transferList []models.Create2Transfer, tx *models.Create2Transfer) bool {
	for i := range transferList {
		if transferList[i].Hash == tx.Hash {
			return true
		}
	}
	return false
}
