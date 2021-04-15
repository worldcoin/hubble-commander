package commander

import (
	"log"
	"math/big"

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
		txError, appError := ApplyTransfer(stateTree, &transfer, feeReceiverTokenIndex)
		if appError != nil {
			return nil, appError
		}
		if txError == nil {
			validTransfers = append(validTransfers, transfer)
			combinedFee.Add(&combinedFee.Int, &transfer.Fee.Int)
		} else {
			if txError == ErrNonceTooHigh {
				continue
			}
			err := storage.SetTransactionError(transfer.Hash, txError.Error())
			if err != nil {
				log.Printf("Setting transaction error failed: %s", err)
			}
			log.Printf("Transfer failed: %s", txError)
		}

		if uint32(len(validTransfers)) == cfg.TxsPerCommitment {
			break
		}
	}

	if combinedFee.Cmp(big.NewInt(0)) == 1 {
		// TODO cfg.FeeReceiverIndex actually represents PubKeyID and is used as StateID here
		err := ApplyFee(stateTree, cfg.FeeReceiverIndex, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return validTransfers, nil
}
