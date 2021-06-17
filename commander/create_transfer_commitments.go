package commander

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrNotEnoughTransfers = NewRollupError("not enough transfers")
)

func (t *transactionExecutor) createTransferCommitments(
	pendingTransfers []models.Transfer,
	domain *bls.Domain,
) ([]models.Commitment, error) {
	if len(pendingTransfers) < int(t.cfg.TxsPerCommitment) {
		return []models.Commitment{}, nil
	}

	commitments := make([]models.Commitment, 0, t.cfg.MaxCommitmentsPerBatch)

	for len(commitments) != int(t.cfg.MaxCommitmentsPerBatch) {
		var commitment *models.Commitment
		var err error

		pendingTransfers, commitment, err = t.createTransferCommitment(pendingTransfers, domain)
		if err != nil {
			return nil, err
		}
		if commitment == nil {
			break
		}

		commitments = append(commitments, *commitment)
	}

	return commitments, nil
}

func (t *transactionExecutor) createTransferCommitment(
	pendingTransfers []models.Transfer,
	domain *bls.Domain,
) (
	newPendingTransfers []models.Transfer,
	commitment *models.Commitment,
	err error,
) {
	startTime := time.Now()

	initialStateRoot, err := t.stateTree.Root()
	if err != nil {
		return nil, nil, err
	}

	var feeReceiverStateID uint32
	appliedTransfers := make([]models.Transfer, 0, t.cfg.TxsPerCommitment)
	invalidTransfers := make([]models.Transfer, 0, 1)

	for {
		if len(pendingTransfers) == 0 {
			pendingTransfers, err = t.storage.GetPendingTransfers(t.cfg.PendingTxsCountMultiplier * t.cfg.TxsPerCommitment)
			if err != nil || len(pendingTransfers) == 0 {
				return nil, nil, err
			}
		}

		var transfers *AppliedTransfers

		transfers, err = t.ApplyTransfers(pendingTransfers, t.cfg.TxsPerCommitment-uint64(len(appliedTransfers)))
		if err != nil {
			return nil, nil, err
		}
		if transfers == nil {
			return nil, nil, ErrNotEnoughTransfers
		}

		appliedTransfers = append(appliedTransfers, transfers.appliedTransfers...)
		invalidTransfers = append(invalidTransfers, transfers.invalidTransfers...)

		if len(appliedTransfers) == int(t.cfg.TxsPerCommitment) {
			feeReceiverStateID = *transfers.feeReceiverStateID
			break
		}

		limit := t.cfg.PendingTxsCountMultiplier*t.cfg.TxsPerCommitment + uint64(len(appliedTransfers)+len(invalidTransfers))
		pendingTransfers, err = t.storage.GetPendingTransfers(limit)
		if err != nil {
			return nil, nil, err
		}

		pendingTransfers = removeTransfer(pendingTransfers, append(appliedTransfers, invalidTransfers...))

		if len(pendingTransfers) == 0 {
			err = t.stateTree.RevertTo(*initialStateRoot)
			return nil, nil, err
		}
	}

	newPendingTransfers = removeTransfer(pendingTransfers, append(appliedTransfers, invalidTransfers...))

	serializedTxs, err := encoder.SerializeTransfers(appliedTransfers)
	if err != nil {
		return nil, nil, err
	}

	combinedSignature, err := combineTransferSignatures(appliedTransfers, domain)
	if err != nil {
		return nil, nil, err
	}

	commitment, err = t.createAndStoreCommitment(txtype.Transfer, feeReceiverStateID, serializedTxs, combinedSignature)
	if err != nil {
		return nil, nil, err
	}

	err = t.markTransfersAsIncluded(appliedTransfers, commitment.ID)
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

func removeTransfer(transferList, toRemove []models.Transfer) []models.Transfer {
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
		sig, err := bls.NewSignatureFromBytes(transfers[i].Signature[:], *domain)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sig)
	}
	return bls.NewAggregatedSignature(signatures).ModelsSignature(), nil
}

func (t *transactionExecutor) markTransfersAsIncluded(transfers []models.Transfer, commitmentID int32) error {
	hashes := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		hashes = append(hashes, transfers[i].Hash)
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, &commitmentID)
}
