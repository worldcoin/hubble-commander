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

func (t *TransactionExecutor) createCreate2TransferCommitments(
	pendingTransfers []models.Create2Transfer,
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

		transfers, err := t.ApplyCreate2Transfers(pendingTransfers)
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

		pendingTransfers = removeCreate2Transfer(pendingTransfers, append(transfers.appliedTransfers, transfers.invalidTransfers...))

		serializedTxs, err := encoder.SerializeCreate2Transfers(transfers.appliedTransfers, transfers.addedPubKeyIDs)
		if err != nil {
			return nil, err
		}

		combinedSignature, err := combineCreate2TransferSignatures(transfers.appliedTransfers, domain)
		if err != nil {
			return nil, err
		}

		commitment, err := t.createAndStoreCommitment(txtype.Create2Transfer, *transfers.feeReceiverStateID, serializedTxs, combinedSignature)
		if err != nil {
			return nil, err
		}

		err = t.markCreate2TransfersAsIncluded(transfers.appliedTransfers, commitment.ID)
		if err != nil {
			return nil, err
		}
		err = t.setCreate2TransferToStateID(transfers.appliedTransfers)
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)
		log.Printf(
			"Created a %s commitment from %d transactions in %s",
			txtype.Create2Transfer,
			len(transfers.appliedTransfers),
			time.Since(startTime).Round(time.Millisecond).String(),
		)
	}

	return commitments, nil
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
		sig, err := bls.NewSignatureFromBytes(transfers[i].Signature.Bytes(), *domain)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sig)
	}
	return bls.NewAggregatedSignature(signatures).ModelsSignature(), nil
}

func (t *TransactionExecutor) markCreate2TransfersAsIncluded(transfers []models.Create2Transfer, commitmentID int32) error {
	hashes := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		hashes = append(hashes, transfers[i].Hash)
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, &commitmentID)
}

func (t *TransactionExecutor) setCreate2TransferToStateID(transfers []models.Create2Transfer) error {
	for i := range transfers {
		err := t.storage.SetCreate2TransferToStateID(transfers[i].Hash, *transfers[i].ToStateID)
		if err != nil {
			return err
		}
	}
	return nil
}
