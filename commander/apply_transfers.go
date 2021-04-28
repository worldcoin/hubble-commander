package commander

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ApplyTransfers(
	storage *st.Storage,
	transfers []models.Transfer,
	cfg *config.RollupConfig,
) (
	appliedTransfers []models.Transfer,
	invalidTransfers []models.Transfer,
	err error,
) {
	stateTree := st.NewStateTree(storage)
	appliedTransfers = make([]models.Transfer, 0, cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	feeReceiverLeaf, err := stateTree.Leaf(cfg.FeeReceiverIndex)
	if err != nil {
		return nil, nil, err
	}

	feeReceiverTokenIndex := feeReceiverLeaf.TokenIndex

	for i := range transfers {
		transfer := transfers[i]
		transferError, appError := ApplyTransfer(stateTree, &transfer, feeReceiverTokenIndex)
		if appError != nil {
			return nil, nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(storage, &transfer.TransactionBase, transferError)
			invalidTransfers = append(invalidTransfers, transfer)
			continue
		}

		appliedTransfers = append(appliedTransfers, transfer)
		combinedFee = *combinedFee.Add(&transfer.Fee)

		if uint32(len(appliedTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if len(appliedTransfers) > 0 {
		err = ApplyFee(stateTree, cfg.FeeReceiverIndex, combinedFee)
		if err != nil {
			return nil, nil, err
		}
	}

	return appliedTransfers, invalidTransfers, nil
}
