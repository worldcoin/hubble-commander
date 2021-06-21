package executor

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

func (t *TransactionExecutor) createTransferCommitments(
	pendingTransfers []models.Transfer,
	domain *bls.Domain,
) ([]models.Commitment, error) {
	stateTree := st.NewStateTree(t.storage)
	commitments := make([]models.Commitment, 0, 32)

	for {
		if len(commitments) >= int(t.cfg.MaxCommitmentsPerBatch) {
			break
		}
		if len(pendingTransfers) < int(t.cfg.TxsPerCommitment) {
			break
		}

		startTime := time.Now()

		initialStateRoot, err := stateTree.Root()
		if err != nil {
			return nil, err
		}

		transfers, err := t.ApplyTransfers(pendingTransfers)
		if err != nil {
			return nil, err
		}

		if len(transfers.appliedTransfers) < int(t.cfg.TxsPerCommitment) {
			err = stateTree.RevertTo(*initialStateRoot)
			if err != nil {
				return nil, err
			}
			break
		}

		pendingTransfers = removeTransfer(pendingTransfers, append(transfers.appliedTransfers, transfers.invalidTransfers...))

		serializedTxs, err := encoder.SerializeTransfers(transfers.appliedTransfers)
		if err != nil {
			return nil, err
		}

		combinedSignature, err := combineTransferSignatures(transfers.appliedTransfers, domain)
		if err != nil {
			return nil, err
		}

		commitment, err := t.createAndStoreCommitment(txtype.Transfer, *transfers.feeReceiverStateID, serializedTxs, combinedSignature)
		if err != nil {
			return nil, err
		}

		err = t.markTransfersAsIncluded(transfers.appliedTransfers, commitment.ID)
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)
		log.Printf(
			"Created a %s commitment from %d transactions in %s",
			txtype.Transfer,
			len(transfers.appliedTransfers),
			time.Since(startTime).Round(time.Millisecond).String(),
		)
	}

	return commitments, nil
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
