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
	feeReceiverStateID *uint32,
	err error,
) {
	if len(transfers) == 0 {
		return
	}

	appliedTransfers = make([]models.Transfer, 0, cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	senderLeaf, err := storage.GetStateLeafByStateID(transfers[0].FromStateID)
	if err != nil {
		return nil, nil, nil, err
	}

	commitmentTokenIndex := senderLeaf.TokenIndex

	for i := range transfers {
		transfer := transfers[i]
		transferError, appError := ApplyTransfer(storage, &transfer, commitmentTokenIndex)
		if appError != nil {
			return nil, nil, nil, appError
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
		feeReceiverStateID, err = ApplyFee(storage, cfg.FeeReceiverPubKeyID, commitmentTokenIndex, combinedFee)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return appliedTransfers, invalidTransfers, feeReceiverStateID, nil
}
