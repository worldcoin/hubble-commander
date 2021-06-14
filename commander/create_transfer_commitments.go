package commander

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

func (t *transactionExecutor) createTransferCommitments(
	pendingTransfers []models.Transfer,
	domain *bls.Domain,
) ([]models.Commitment, error) {
	commitments := make([]models.Commitment, 0, 32)

	for {
		if len(commitments) >= int(t.cfg.MaxCommitmentsPerBatch) {
			break
		}
		if len(pendingTransfers) < int(t.cfg.TxsPerCommitment) {
			break
		}

		commitment, err := t.createTransferCommitment(pendingTransfers, domain)
		if err != nil {
			return nil, err
		}
		if commitment == nil {
			return []models.Commitment{}, nil
		}

		commitments = append(commitments, *commitment)
	}

	return commitments, nil
}

func (t *transactionExecutor) createTransferCommitment(pendingTransfers []models.Transfer, domain *bls.Domain) (*models.Commitment, error) {
	stateTree := st.NewStateTree(t.storage)

	startTime := time.Now()

	initialStateRoot, err := stateTree.Root()
	if err != nil {
		return nil, err
	}

	var feeReceiverStateID uint32
	appliedTransfers := make([]models.Transfer, 0, t.cfg.TxsPerCommitment)
	invalidTransfers := make([]models.Transfer, 0, 1)

	for {
		// nolint:govet
		transfers, err := t.ApplyTransfers(pendingTransfers)
		if err != nil {
			return nil, err
		}

		appliedTransfers = append(appliedTransfers, transfers.appliedTransfers...)
		invalidTransfers = append(invalidTransfers, transfers.invalidTransfers...)

		if len(appliedTransfers) >= int(t.cfg.TxsPerCommitment) {
			feeReceiverStateID = *transfers.feeReceiverStateID
			break
		}

		pendingTransfers, err = t.storage.GetPendingTransfers(t.cfg.TxsPerCommitment, transfers.lastTransactionNonce.AddN(1))
		if err != nil {
			return nil, err
		}

		if len(pendingTransfers) == 0 {
			err = stateTree.RevertTo(*initialStateRoot)
			if err != nil {
				return nil, err
			}
			return nil, nil
		}
	}

	pendingTransfers = removeTransfer(pendingTransfers, append(appliedTransfers, invalidTransfers...))

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

	err = t.markTransfersAsIncluded(appliedTransfers, commitment.ID)
	if err != nil {
		return nil, err
	}

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		txtype.Transfer,
		len(appliedTransfers),
		time.Since(startTime).Round(time.Millisecond).String(),
	)

	return commitment, nil
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
