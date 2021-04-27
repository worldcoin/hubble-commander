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
	[]models.Transfer,
	error,
) {
	stateTree := st.NewStateTree(storage)
	validTransfers := make([]models.Transfer, 0, cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	feeReceiverLeaf, err := stateTree.Leaf(cfg.FeeReceiverIndex)
	if err != nil {
		return nil, err
	}

	feeReceiverTokenIndex := feeReceiverLeaf.TokenIndex

	for i := range transfers {
		transfer := transfers[i]
		transferError, appError := ApplyTransfer(stateTree, &transfer, feeReceiverTokenIndex)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			// TODO optimise to not process failed transactions for subsequent commitments
			logAndSaveTransactionError(storage, &transfer.TransactionBase, transferError)
			continue
		}

		validTransfers = append(validTransfers, transfer)
		combinedFee = *combinedFee.Add(&transfer.Fee)

		if uint32(len(validTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if len(validTransfers) > 0 {
		err = ApplyFee(stateTree, cfg.FeeReceiverIndex, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return validTransfers, nil
}
