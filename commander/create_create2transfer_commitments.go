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
	ErrNotEnoughC2Transfers = NewRollupError("not enough create2transfers")
)

func (t *transactionExecutor) createCreate2TransferCommitments(
	pendingTransfers []models.Create2Transfer,
	domain *bls.Domain,
) ([]models.Commitment, error) {
	if len(pendingTransfers) < int(t.cfg.TxsPerCommitment) {
		return []models.Commitment{}, nil
	}

	commitments := make([]models.Commitment, 0, t.cfg.MaxCommitmentsPerBatch)

	for len(commitments) != int(t.cfg.MaxCommitmentsPerBatch) {
		var commitment *models.Commitment
		var err error

		pendingTransfers, commitment, err = t.createC2TCommitment(pendingTransfers, domain)
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

func (t *transactionExecutor) createC2TCommitment(
	pendingTransfers []models.Create2Transfer,
	domain *bls.Domain,
) (
	newPendingTransfers []models.Create2Transfer,
	commitment *models.Commitment,
	err error,
) {
	startTime := time.Now()

	initialStateRoot, err := t.stateTree.Root()
	if err != nil {
		return nil, nil, err
	}

	var feeReceiverStateID uint32
	appliedTransfers := make([]models.Create2Transfer, 0, t.cfg.TxsPerCommitment)
	invalidTransfers := make([]models.Create2Transfer, 0, 1)
	addedPubKeyIDs := make([]uint32, 0, t.cfg.TxsPerCommitment)

	for {
		if len(pendingTransfers) == 0 {
			pendingTransfers, err = t.storage.GetPendingCreate2Transfers(pendingTxsCountMultiplier * t.cfg.TxsPerCommitment)
			if err != nil || len(pendingTransfers) == 0 {
				return nil, nil, err
			}
		}

		var transfers *AppliedC2Transfers

		maxAppliedTransfers := t.cfg.TxsPerCommitment - uint64(len(appliedTransfers))
		transfers, err = t.ApplyCreate2Transfers(pendingTransfers, maxAppliedTransfers)
		if err != nil {
			return nil, nil, err
		}
		if transfers == nil {
			return nil, nil, ErrNotEnoughC2Transfers
		}

		appliedTransfers = append(appliedTransfers, transfers.appliedTransfers...)
		invalidTransfers = append(invalidTransfers, transfers.invalidTransfers...)
		addedPubKeyIDs = append(addedPubKeyIDs, transfers.addedPubKeyIDs...)

		if len(appliedTransfers) == int(t.cfg.TxsPerCommitment) {
			feeReceiverStateID = *transfers.feeReceiverStateID
			break
		}

		limit := pendingTxsCountMultiplier*t.cfg.TxsPerCommitment + uint64(len(appliedTransfers)+len(invalidTransfers))
		pendingTransfers, err = t.storage.GetPendingCreate2Transfers(limit)
		if err != nil {
			return nil, nil, err
		}

		// TODO - instead of doing this use SQL Offset (needs proper mempool)
		pendingTransfers = removeCreate2Transfer(pendingTransfers, append(appliedTransfers, invalidTransfers...))

		if len(pendingTransfers) == 0 {
			err = t.stateTree.RevertTo(*initialStateRoot)
			return nil, nil, err
		}
	}

	newPendingTransfers = removeCreate2Transfer(pendingTransfers, append(appliedTransfers, invalidTransfers...))

	serializedTxs, err := encoder.SerializeCreate2Transfers(appliedTransfers, addedPubKeyIDs)
	if err != nil {
		return nil, nil, err
	}

	combinedSignature, err := combineCreate2TransferSignatures(appliedTransfers, domain)
	if err != nil {
		return nil, nil, err
	}

	commitment, err = t.createAndStoreCommitment(txtype.Create2Transfer, feeReceiverStateID, serializedTxs, combinedSignature)
	if err != nil {
		return nil, nil, err
	}

	err = t.markCreate2TransfersAsIncluded(appliedTransfers, commitment.ID)
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

func removeCreate2Transfer(transferList, toRemove []models.Create2Transfer) []models.Create2Transfer {
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

func combineCreate2TransferSignatures(transfers []models.Create2Transfer, domain *bls.Domain) (*models.Signature, error) {
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

func (t *transactionExecutor) markCreate2TransfersAsIncluded(transfers []models.Create2Transfer, commitmentID int32) error {
	hashes := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		hashes = append(hashes, transfers[i].Hash)

		err := t.storage.SetCreate2TransferToStateID(transfers[i].Hash, *transfers[i].ToStateID)
		if err != nil {
			return err
		}
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, &commitmentID)
}
