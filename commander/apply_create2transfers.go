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
	successfullyAddedPubKeyIDs []uint32,
	cfg *config.RollupConfig,
) (
	[]models.Create2Transfer,
	[]uint32,
	error,
) {
	stateTree := st.NewStateTree(storage)
	validTransfers := make([]models.Create2Transfer, 0, cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	feeReceiverLeaf, err := stateTree.Leaf(cfg.FeeReceiverIndex)
	if err != nil {
		return nil, nil, err
	}

	feeReceiverTokenIndex := feeReceiverLeaf.TokenIndex

	for i := range transfers {
		transfer := transfers[i]
		currentTransferError := ErrAccountAlreadyExists

		if !uint32InSlice(transfer.ToPubKeyID, successfullyAddedPubKeyIDs) {
			addedPubKeyID, transferError, appError := ApplyCreate2Transfer(storage, &transfer, feeReceiverTokenIndex)
			if appError != nil {
				return nil, nil, appError
			}
			if transferError == nil {
				validTransfers = append(validTransfers, transfer)
				successfullyAddedPubKeyIDs = append(successfullyAddedPubKeyIDs, *addedPubKeyID)
				combinedFee = *combinedFee.Add(&transfer.Fee)
			}
			currentTransferError = transferError
		}

		logAndSaveTransactionError(storage, &transfer.TransactionBase, currentTransferError)

		if uint32(len(validTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	err = ApplyFee(stateTree, cfg.FeeReceiverIndex, combinedFee)
	if err != nil {
		return nil, nil, err
	}

	return validTransfers, successfullyAddedPubKeyIDs, nil
}

func uint32InSlice(a uint32, list []uint32) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
