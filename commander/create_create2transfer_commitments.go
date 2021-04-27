package commander

import (
	"log"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func createCreate2TransferCommitments(
	pendingTransfers []models.Create2Transfer,
	storage *st.Storage,
	cfg *config.RollupConfig,
) ([]models.Commitment, error) {
	stateTree := st.NewStateTree(storage)
	commitments := make([]models.Commitment, 0, 32)
	successfullyAddedPubKeyIDs := make([]uint32, 0, 1)

	for {
		if len(commitments) >= int(cfg.MaxCommitmentsPerBatch) {
			break
		}

		initialStateRoot, err := stateTree.Root()
		if err != nil {
			return nil, err
		}

		includedTransfers, addedPubKeyIDs, err := ApplyCreate2Transfers(storage, pendingTransfers, successfullyAddedPubKeyIDs, cfg)
		if err != nil {
			return nil, err
		}

		for i := range addedPubKeyIDs {
			successfullyAddedPubKeyIDs = append(successfullyAddedPubKeyIDs, addedPubKeyIDs[i])
		}

		if len(includedTransfers) < int(cfg.TxsPerCommitment) {
			err = stateTree.RevertTo(*initialStateRoot)
			if err != nil {
				return nil, err
			}
			break
		}

		pendingTransfers = removeCreate2Transfer(pendingTransfers, includedTransfers)

		serializedTxs, err := encoder.SerializeCreate2Transfers(includedTransfers)
		if err != nil {
			return nil, err
		}

		combinedSignature, err := combineCreate2TransferSignatures(includedTransfers)
		if err != nil {
			return nil, err
		}

		log.Printf("Creating a create2Transfer commitment from %d transactions", len(includedTransfers))
		commitment, err := createAndStoreCommitment(storage, txtype.Create2Transfer, cfg.FeeReceiverIndex, serializedTxs, combinedSignature)
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)

		err = markCreate2TransfersAsIncluded(storage, includedTransfers, commitment.ID)
		if err != nil {
			return nil, err
		}
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

func combineCreate2TransferSignatures(transfers []models.Create2Transfer) (*models.Signature, error) {
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

func markCreate2TransfersAsIncluded(storage *st.Storage, transfers []models.Create2Transfer, commitmentID int32) error {
	includedTxBases := make([]models.TransactionBase, 0, len(transfers))
	for i := range transfers {
		includedTxBases = append(includedTxBases, transfers[i].TransactionBase)
	}

	err := markTransactionsAsIncluded(storage, includedTxBases, commitmentID)
	if err != nil {
		return err
	}

	return nil
}
