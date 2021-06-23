package executor

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrNotEnoughTransfers = NewRollupError("not enough transfers")
)

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

	initialStateRoot, err := t.stateTree.Root()
	if err != nil {
		return nil, nil, err
	}

	appliedTransfers, newPendingTransfers, feeReceiverStateID, err := t.applyTransfersForCommitment(pendingTransfers)
	if err == ErrNotEnoughTransfers {
		if revertErr := t.stateTree.RevertTo(*initialStateRoot); revertErr != nil {
			return nil, nil, revertErr
		}
		return nil, nil, err
	}
	if err != nil {
		return nil, nil, err
	}

	commitment, err = t.buildTransferCommitment(appliedTransfers, *feeReceiverStateID, domain)
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

func (t *TransactionExecutor) applyTransfersForCommitment(pendingTransfers []models.Transfer) (
	appliedTransfers, newPendingTransfers []models.Transfer,
	feeReceiverStateID *uint32,
	err error,
) {
	appliedTransfers = make([]models.Transfer, 0, t.cfg.TxsPerCommitment)
	invalidTransfers := make([]models.Transfer, 0, 1)

	for {
		var transfers *AppliedTransfers

		numNeededTransfers := t.cfg.TxsPerCommitment - uint32(len(appliedTransfers))
		transfers, err = t.ApplyTransfers(pendingTransfers, numNeededTransfers)
		if err != nil {
			return nil, nil, nil, err
		}

		appliedTransfers = append(appliedTransfers, transfers.appliedTransfers...)
		invalidTransfers = append(invalidTransfers, transfers.invalidTransfers...)

		if len(appliedTransfers) == int(t.cfg.TxsPerCommitment) {
			feeReceiverStateID = transfers.feeReceiverStateID
			newPendingTransfers = removeTransfers(pendingTransfers, append(appliedTransfers, invalidTransfers...))
			return appliedTransfers, newPendingTransfers, feeReceiverStateID, nil
		}

		pendingTransfers, err = t.queryMorePendingTransfers(appliedTransfers, invalidTransfers)
		if err != nil {
			return nil, nil, nil, err
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

func (t *TransactionExecutor) queryMorePendingTransfers(appliedTransfers, invalidTransfers []models.Transfer) ([]models.Transfer, error) {
	// TODO use SQL Offset instead
	alreadySeenTransfers := append(appliedTransfers, invalidTransfers...) // nolint:gocritic
	pendingTransfers, err := t.storage.GetPendingTransfers(
		t.cfg.MaxCommitmentsPerBatch*t.cfg.TxsPerCommitment + uint32(len(alreadySeenTransfers)),
	)
	if err != nil {
		return nil, err
	}
	pendingTransfers = removeTransfers(pendingTransfers, alreadySeenTransfers)

	numNeededTransfers := t.cfg.TxsPerCommitment - uint32(len(appliedTransfers))
	if len(pendingTransfers) < int(numNeededTransfers) {
		return nil, ErrNotEnoughTransfers
	}
	return pendingTransfers, nil
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

func combineTransferSignatures(transfers []models.Transfer, domain *bls.Domain) (*models.Signature, error) {
	signatures := make([]*bls.Signature, 0, len(transfers))
	for i := range transfers {
		sig, err := bls.NewSignatureFromBytes(transfers[i].Signature.Bytes(), *domain)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sig)
	}
	return bls.NewAggregatedSignature(signatures).ModelsSignature(), nil
}

func (t *TransactionExecutor) markTransfersAsIncluded(transfers []models.Transfer, commitmentID int32) error {
	hashes := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		hashes = append(hashes, transfers[i].Hash)
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, &commitmentID)
}
