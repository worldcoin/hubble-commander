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

func createCommitments(
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

		initialStateRoot, err := stateTree.Root()
		if err != nil {
			return nil, err
		}

		startTime := time.Now()

		includedTransfers, err := ApplyTransfers(storage, pendingTransfers, cfg)
		if err != nil {
			return nil, err
		}

		if len(includedTransfers) < int(cfg.TxsPerCommitment) {
			err = stateTree.RevertTo(*initialStateRoot)
			if err != nil {
				return nil, err
			}
			break
		}

		pendingTransfers = removeTransfer(pendingTransfers, includedTransfers)

		log.Printf("Creating a commitment from %d transactions in %d ms", len(includedTransfers), time.Since(startTime).Milliseconds())
		commitment, err := createAndStoreCommitment(storage, includedTransfers, cfg.FeeReceiverIndex)
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)

		err = markTransactionsAsIncluded(storage, includedTransfers, commitment.ID)
		if err != nil {
			return nil, err
		}
	}

	return commitments, nil
}

func removeTransfer(transferList, toRemove []models.Transfer) []models.Transfer {
	outputIndex := 0
	for i := range transferList {
		transfer := &transferList[i]
		if !transactionExists(toRemove, transfer) {
			transferList[outputIndex] = *transfer
			outputIndex++
		}
	}

	return transferList[:outputIndex]
}

func transactionExists(transferList []models.Transfer, tx *models.Transfer) bool {
	for i := range transferList {
		if transferList[i].Hash == tx.Hash {
			return true
		}
	}
	return false
}

func createAndStoreCommitment(storage *st.Storage, transfers []models.Transfer, feeReceiverIndex uint32) (*models.Commitment, error) {
	serializedTxs, err := encoder.SerializeTransfers(transfers)
	if err != nil {
		return nil, err
	}

	combinedSignature, err := combineSignatures(transfers)
	if err != nil {
		return nil, err
	}

	stateRoot, err := st.NewStateTree(storage).Root()
	if err != nil {
		return nil, err
	}

	commitment := models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      serializedTxs,
		FeeReceiver:       feeReceiverIndex,
		CombinedSignature: *combinedSignature,
		PostStateRoot:     *stateRoot,
	}

	id, err := storage.AddCommitment(&commitment)
	if err != nil {
		return nil, err
	}

	commitment.ID = *id

	return &commitment, nil
}

func combineSignatures(transfers []models.Transfer) (*models.Signature, error) {
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

func markTransactionsAsIncluded(storage *st.Storage, transfers []models.Transfer, commitmentID int32) error {
	for i := range transfers {
		err := storage.MarkTransactionAsIncluded(transfers[i].Hash, commitmentID)
		if err != nil {
			return err
		}
	}
	return nil
}
