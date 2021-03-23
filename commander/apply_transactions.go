package commander

import (
	"log"
	"math/big"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func ApplyTransactions(
	storage *st.Storage,
	transactions []models.Transaction,
	cfg *config.RollupConfig,
) (
	[]models.Transaction,
	error,
) {
	stateTree := st.NewStateTree(storage)
	validTxs := make([]models.Transaction, 0, cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	for i := range transactions {
		tx := transactions[i]
		txError, appError := ApplyTransfer(stateTree, &tx)
		if appError != nil {
			return nil, appError
		}
		if txError == nil {
			validTxs = append(validTxs, tx)
			combinedFee.Add(&combinedFee.Int, &tx.Fee.Int)
		} else {
			if txError == ErrNonceTooHigh {
				continue
			}
			err := storage.SetTransactionError(tx.Hash, txError.Error())
			if err != nil {
				log.Printf("Setting transaction error failed: %s", err)
			}
			log.Printf("Transaction failed: %s", txError)
		}

		if uint32(len(validTxs)) == cfg.TxsPerCommitment {
			break
		}
	}

	if combinedFee.Cmp(big.NewInt(0)) == 1 {
		// TODO cfg.FeeReceiverIndex actually represents PubKeyID and is used and StateID here
		err := ApplyFee(stateTree, cfg.FeeReceiverIndex, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return validTxs, nil
}
