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
	alreadyAddedPubKeyIDs []uint32,
	cfg *config.RollupConfig,
) (
	[]models.Create2Transfer,
	[]models.Create2Transfer,
	[]uint32,
	error,
) {
	stateTree := st.NewStateTree(storage)
	appliedTransfers := make([]models.Create2Transfer, 0, cfg.TxsPerCommitment)
	invalidTransfers := make([]models.Create2Transfer, 0)
	combinedFee := models.MakeUint256(0)

	feeReceiverLeaf, err := stateTree.Leaf(cfg.FeeReceiverIndex)
	if err != nil {
		return nil, nil, nil, err
	}

	feeReceiverTokenIndex := feeReceiverLeaf.TokenIndex

	for i := range transfers {
		transfer := transfers[i]

		if uint32InSlice(transfer.ToPubKeyID, alreadyAddedPubKeyIDs) {
			logAndSaveTransactionError(storage, &transfer.TransactionBase, ErrAccountAlreadyExists)
			invalidTransfers = append(invalidTransfers, transfer)
			continue
		}

		addedPubKeyID, transferError, appError := ApplyCreate2Transfer(storage, &transfer, feeReceiverTokenIndex)
		if appError != nil {
			return nil, nil, nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(storage, &transfer.TransactionBase, transferError)
			invalidTransfers = append(invalidTransfers, transfer)
			continue
		}

		appliedTransfers = append(appliedTransfers, transfer)
		alreadyAddedPubKeyIDs = append(alreadyAddedPubKeyIDs, *addedPubKeyID)
		combinedFee = *combinedFee.Add(&transfer.Fee)

		if uint32(len(appliedTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if len(appliedTransfers) > 0 {
		err = ApplyFee(stateTree, cfg.FeeReceiverIndex, combinedFee)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return appliedTransfers, invalidTransfers, alreadyAddedPubKeyIDs, nil
}

func uint32InSlice(a uint32, list []uint32) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
