package commander

import (
	"log"

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

	for {
		if len(commitments) >= int(cfg.MaxCommitmentsPerBatch) {
			break
		}

		initialStateRoot, err := stateTree.Root()
		if err != nil {
			return nil, err
		}

		includedTransfers, err := ApplyCreate2Transfers(storage, pendingTransfers, cfg)
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

		pendingTransfers = removeCreate2Transfer(pendingTransfers, includedTransfers)

		log.Printf("Creating a create2Transfer commitment from %d transactions", len(includedTransfers))
		commitment, err := createAndStoreCreate2TransferCommitment(storage, includedTransfers, cfg.FeeReceiverIndex)
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)

		includedTxBases := make([]models.TransactionBase, 0, len(includedTransfers))
		for i := range includedTransfers {
			includedTxBases = append(includedTxBases, includedTransfers[i].TransactionBase)
		}

		err = markTransactionsAsIncluded(storage, includedTxBases, commitment.ID)
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

func createAndStoreCreate2TransferCommitment(storage *st.Storage, transfers []models.Create2Transfer, feeReceiverIndex uint32) (*models.Commitment, error) {
	combinedSignature := models.MakeSignature(1, 2) // TODO: Actually combine signatures

	serializedTxs, err := encoder.SerializeCreate2Transfers(transfers)
	if err != nil {
		return nil, err
	}

	stateRoot, err := st.NewStateTree(storage).Root()
	if err != nil {
		return nil, err
	}

	commitment := models.Commitment{
		Type:              txtype.Create2Transfer,
		Transactions:      serializedTxs,
		FeeReceiver:       feeReceiverIndex,
		CombinedSignature: combinedSignature,
		PostStateRoot:     *stateRoot,
	}

	id, err := storage.AddCommitment(&commitment)
	if err != nil {
		return nil, err
	}

	commitment.ID = *id

	return &commitment, nil
}
