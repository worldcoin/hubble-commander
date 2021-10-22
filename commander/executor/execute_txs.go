package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (c *RollupContext) ExecuteTxs(
	txs models.GenericTransactionArray,
	maxApplied uint32,
	feeReceiver *FeeReceiver,
) (ExecuteTxsResult, error) {
	if txs.Len() == 0 {
		return c.Executor.NewExecuteTxsResult(0), nil
	}

	returnStruct := c.Executor.NewExecuteTxsResult(c.cfg.MaxTxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	for i := 0; i < txs.Len(); i++ {
		if returnStruct.AppliedTxs().Len() == int(maxApplied) {
			break
		}

		applyResult, transferError, appError := c.Executor.ApplyTx(txs.At(i), feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(c.storage, txs.At(i), transferError)
			returnStruct.AddInvalidTx(txs.At(i))
			c.TxErrorsToStore = append(c.TxErrorsToStore, models.TransactionError{
				Hash:         txs.At(i).GetBase().Hash,
				ErrorMessage: transferError.Error(),
			})
			continue
		}

		err := c.Executor.AddPendingAccount(applyResult)
		if err != nil {
			return nil, err
		}

		returnStruct.AddApplied(applyResult)
		fee := applyResult.AppliedTx().GetFee()
		combinedFee = *combinedFee.Add(&fee)
	}

	if returnStruct.AppliedTxs().Len() > 0 {
		_, err := c.ApplyFee(feeReceiver.StateID, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}
