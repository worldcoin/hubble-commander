package commander

import (
	"errors"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

var (
	ErrAccountAlreadyExists = errors.New("account with provided public key already exists")
)

func ApplyCreate2Transfers(
	storage *st.Storage,
	transfers []models.Create2Transfer,
	addedPubKeyIDs map[uint32]struct{},
	cfg *config.RollupConfig,
) (
	appliedTransfers []models.Create2Transfer,
	invalidTransfers []models.Create2Transfer,
	err error,
) {
	stateTree := st.NewStateTree(storage)
	appliedTransfers = make([]models.Create2Transfer, 0, cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	feeReceiverLeaf, err := stateTree.Leaf(cfg.FeeReceiverPubKeyID)
	if err != nil {
		return nil, nil, err
	}

	feeReceiverTokenIndex := feeReceiverLeaf.TokenIndex

	for i := range transfers {
		transfer := transfers[i]

		if _, ok := addedPubKeyIDs[transfer.ToPubKeyID]; ok {
			logAndSaveTransactionError(storage, &transfer.TransactionBase, ErrAccountAlreadyExists)
			invalidTransfers = append(invalidTransfers, transfer)
			continue
		}

		addedPubKeyID, transferError, appError := ApplyCreate2Transfer(storage, &transfer, feeReceiverTokenIndex)
		if appError != nil {
			return nil, nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(storage, &transfer.TransactionBase, transferError)
			invalidTransfers = append(invalidTransfers, transfer)
			continue
		}

		appliedTransfers = append(appliedTransfers, transfer)
		addedPubKeyIDs[*addedPubKeyID] = struct{}{}
		combinedFee = *combinedFee.Add(&transfer.Fee)

		if uint32(len(appliedTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if len(appliedTransfers) > 0 {
		err = ApplyFee(stateTree, cfg.FeeReceiverPubKeyID, combinedFee)
		if err != nil {
			return nil, nil, err
		}
	}

	return appliedTransfers, invalidTransfers, nil
}
