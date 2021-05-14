package commander

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ApplyCreate2Transfers(
	storage *st.Storage,
	client *eth.Client,
	transfers []models.Create2Transfer,
	cfg *config.RollupConfig,
) (
	appliedTransfers []models.Create2Transfer,
	invalidTransfers []models.Create2Transfer,
	addedPubKeyIDs []uint32,
	feeReceiverStateID *uint32,
	err error,
) {
	if len(transfers) == 0 {
		return
	}

	stateTree := st.NewStateTree(storage)
	appliedTransfers = make([]models.Create2Transfer, 0, cfg.TxsPerCommitment)
	addedPubKeyIDs = make([]uint32, 0, cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	senderLeaf, err := stateTree.Leaf(transfers[0].FromStateID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	commitmentTokenIndex := senderLeaf.TokenIndex

	for i := range transfers {
		transfer := transfers[i]

		addedPubKeyID, transferError, appError := ApplyCreate2Transfer(storage, client, &transfer, commitmentTokenIndex)
		if appError != nil {
			return nil, nil, nil, nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(storage, &transfer.TransactionBase, transferError)
			invalidTransfers = append(invalidTransfers, transfer)
			continue
		}

		addedPubKeyIDs = append(addedPubKeyIDs, *addedPubKeyID)
		appliedTransfers = append(appliedTransfers, transfer)
		combinedFee = *combinedFee.Add(&transfer.Fee)

		if uint32(len(appliedTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if len(appliedTransfers) > 0 {
		feeReceiverStateID, err = ApplyFee(stateTree, storage, cfg.FeeReceiverPubKeyID, commitmentTokenIndex, combinedFee)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	return appliedTransfers, invalidTransfers, addedPubKeyIDs, feeReceiverStateID, nil
}
