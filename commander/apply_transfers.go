package commander

import (
	"log"

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
		if transferError == nil {
			validTransfers = append(validTransfers, transfer)
			combinedFee = *combinedFee.Add(&transfer.Fee)
		} else {
			if transferError == ErrNonceTooHigh {
				continue
			}
			err := storage.SetTransactionError(transfer.Hash, transferError.Error())
			if err != nil {
				log.Printf("Setting transaction error failed: %s", err)
			}
			log.Printf("Transfer failed: %s", transferError)
		}

		if uint32(len(validTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if combinedFee.CmpN(0) == 1 {
		// TODO cfg.FeeReceiverIndex actually represents PubKeyID and is used as StateID here
		err := ApplyFee(stateTree, cfg.FeeReceiverIndex, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return validTransfers, nil
}
