package executor

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

var (
	ErrNotEnoughTransfers = NewRollupError("not enough transfers")
)

type FeeReceiver struct {
	StateID uint32
	TokenID models.Uint256
}

func (t *TransactionExecutor) CreateTransferCommitments(
	domain *bls.Domain,
) (commitments []models.Commitment, err error) {
	pendingTransfers, err := t.queryPendingTransfers()
	if err != nil {
		return nil, err
	}

	commitments = make([]models.Commitment, 0, t.cfg.MaxCommitmentsPerBatch)

	for len(commitments) != int(t.cfg.MaxCommitmentsPerBatch) {
		var commitment *models.Commitment

		pendingTransfers, commitment, err = t.createTransferCommitment(pendingTransfers, domain)
		if err == ErrNotEnoughTransfers {
			break
		}
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)
	}

	if len(commitments) == 0 {
		return nil, ErrNotEnoughTransfers
	}

	return commitments, nil
}

func (t *TransactionExecutor) createTransferCommitment(
	pendingTransfers []models.Transfer,
	domain *bls.Domain,
) (
	newPendingTransfers []models.Transfer,
	commitment *models.Commitment,
	err error,
) {
	startTime := time.Now()

	pendingTransfers, err = t.refillPendingTransfers(pendingTransfers)
	if err != nil {
		return nil, nil, err
	}

	feeReceiver, err := t.getCommitmentFeeReceiver()
	if err != nil {
		return nil, nil, err
	}

	initialStateRoot, err := t.stateTree.Root()
	if err != nil {
		return nil, nil, err
	}

	appliedTransfers, newPendingTransfers, err := t.applyTransfersForCommitment(pendingTransfers, feeReceiver)
	if err == ErrNotEnoughTransfers {
		if revertErr := t.stateTree.RevertTo(*initialStateRoot); revertErr != nil {
			return nil, nil, revertErr
		}
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}

	commitment, err = t.buildTransferCommitment(appliedTransfers, feeReceiver.StateID, domain)
	if err != nil {
		return nil, nil, err
	}

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		txtype.Transfer,
		len(appliedTransfers),
		time.Since(startTime).Round(time.Millisecond).String(),
	)

	return newPendingTransfers, commitment, nil
}

func (t *TransactionExecutor) applyTransfersForCommitment(pendingTransfers []models.Transfer, feeReceiver *FeeReceiver) (
	appliedTransfers, newPendingTransfers []models.Transfer,
	err error,
) {
	appliedTransfers = make([]models.Transfer, 0, t.cfg.TxsPerCommitment)
	invalidTransfers := make([]models.Transfer, 0, 1)

	for {
		var transfers *AppliedTransfers

		numNeededTransfers := t.cfg.TxsPerCommitment - uint32(len(appliedTransfers))
		transfers, err = t.ApplyTransfers(pendingTransfers, numNeededTransfers, feeReceiver, false)
		if err != nil {
			return nil, nil, err
		}

		appliedTransfers = append(appliedTransfers, transfers.appliedTransfers...)
		invalidTransfers = append(invalidTransfers, transfers.invalidTransfers...)

		if len(appliedTransfers) == int(t.cfg.TxsPerCommitment) {
			newPendingTransfers = removeTransfers(pendingTransfers, append(appliedTransfers, invalidTransfers...))
			return appliedTransfers, newPendingTransfers, nil
		}

		pendingTransfers, err = t.queryMorePendingTransfers(appliedTransfers)
		if err != nil {
			return nil, nil, err
		}
	}
}

func (t *TransactionExecutor) refillPendingTransfers(pendingTransfers []models.Transfer) ([]models.Transfer, error) {
	if len(pendingTransfers) < int(t.cfg.TxsPerCommitment) {
		return t.queryPendingTransfers()
	}
	return pendingTransfers, nil
}

func (t *TransactionExecutor) queryPendingTransfers() ([]models.Transfer, error) {
	pendingTransfers, err := t.storage.GetPendingTransfers(t.cfg.MaxCommitmentsPerBatch * t.cfg.TxsPerCommitment)
	if err != nil {
		return nil, err
	}
	if len(pendingTransfers) < int(t.cfg.TxsPerCommitment) {
		return nil, ErrNotEnoughTransfers
	}
	return pendingTransfers, nil
}

func (t *TransactionExecutor) queryMorePendingTransfers(appliedTransfers []models.Transfer) ([]models.Transfer, error) {
	numAppliedTransfers := uint32(len(appliedTransfers))
	// TODO use SQL Offset instead
	pendingTransfers, err := t.storage.GetPendingTransfers(
		t.cfg.MaxCommitmentsPerBatch*t.cfg.TxsPerCommitment + numAppliedTransfers,
	)
	if err != nil {
		return nil, err
	}
	pendingTransfers = removeTransfers(pendingTransfers, appliedTransfers)

	numNeededTransfers := t.cfg.TxsPerCommitment - numAppliedTransfers
	if len(pendingTransfers) < int(numNeededTransfers) {
		return nil, ErrNotEnoughTransfers
	}
	return pendingTransfers, nil
}

func (t *TransactionExecutor) getCommitmentFeeReceiver() (*FeeReceiver, error) {
	commitmentTokenID := models.MakeUint256(0) // TODO support multiple tokens
	feeReceiverState, err := t.storage.GetFeeReceiverStateLeaf(t.cfg.FeeReceiverPubKeyID, commitmentTokenID)
	if err != nil {
		return nil, err
	}
	return &FeeReceiver{
		StateID: feeReceiverState.StateID,
		TokenID: feeReceiverState.TokenID,
	}, nil
}

func removeTransfers(transferList, toRemove []models.Transfer) []models.Transfer {
	outputIndex := 0
	for i := range transferList {
		transfer := &transferList[i]
		if !transferExists(toRemove, transfer) {
			transferList[outputIndex] = *transfer
			outputIndex++
		}
	}

	return transferList[:outputIndex]
}

func transferExists(transferList []models.Transfer, tx *models.Transfer) bool {
	for i := range transferList {
		if transferList[i].Hash == tx.Hash {
			return true
		}
	}
	return false
}
