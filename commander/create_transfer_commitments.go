package commander

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

var mockDomain = bls.Domain{1, 2, 3, 4} // TODO use real domain

func createTransferCommitments(
	pendingTransfers []models.Transfer,
	storage *st.Storage,
	cfg *config.RollupConfig,
) ([]models.Commitment, error) {
	stateTree := st.NewStateTree(storage)
	commitments := make([]models.Commitment, 0, 32)

	for {
		if len(commitments) >= int(cfg.MaxCommitmentsPerBatch) {
			break
		}
		startTime := time.Now()

		initialStateRoot, err := stateTree.Root()
		if err != nil {
			return nil, err
		}

		appliedTxs, invalidTxs, feeReceiverStateID, err := ApplyTransfers(storage, pendingTransfers, cfg)
		if err != nil {
			return nil, err
		}

		if len(appliedTxs) < int(cfg.TxsPerCommitment) {
			err = stateTree.RevertTo(*initialStateRoot)
			if err != nil {
				return nil, err
			}
			break
		}

		pendingTransfers = removeTransfer(pendingTransfers, append(appliedTxs, invalidTxs...))

		serializedTxs, err := encoder.SerializeTransfers(appliedTxs)
		if err != nil {
			return nil, err
		}

		combinedSignature, err := combineTransferSignatures(appliedTxs)
		if err != nil {
			return nil, err
		}

		commitment, err := createAndStoreCommitment(storage, txtype.Transfer, *feeReceiverStateID, serializedTxs, combinedSignature)
		if err != nil {
			return nil, err
		}

		err = markTransfersAsIncluded(storage, appliedTxs, commitment.ID)
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)
		log.Printf(
			"Created a %s commitment from %d transactions in %d ms",
			txtype.Transfer,
			len(appliedTxs),
			time.Since(startTime).Milliseconds(),
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

func combineTransferSignatures(transfers []models.Transfer) (*models.Signature, error) {
	signatures := make([]*bls.Signature, 0, len(transfers))
	for i := range transfers {
		sig, err := bls.NewSignatureFromBytes(transfers[i].Signature[:], mockDomain)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sig)
	}
	return bls.NewAggregatedSignature(signatures).ModelsSignature(), nil
}

func markTransfersAsIncluded(storage *st.Storage, transfers []models.Transfer, commitmentID int32) error {
	for i := range transfers {
		err := storage.MarkTransactionAsIncluded(transfers[i].Hash, commitmentID)
		if err != nil {
			return err
		}
	}
	return nil
}
