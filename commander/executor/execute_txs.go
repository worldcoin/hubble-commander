package executor

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
)

func (c *TransactionsContext) ExecuteTxs(txs models.GenericTransactionArray, feeReceiver *FeeReceiver) (ExecuteTxsResult, error) {
	if txs.Len() == 0 {
		return c.Executor.NewExecuteTxsResult(0), nil
	}

	returnStruct := c.Executor.NewExecuteTxsResult(c.cfg.MaxTxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	for i := 0; i < txs.Len(); i++ {
		if returnStruct.AppliedTxs().Len() == int(c.cfg.MaxTxsPerCommitment) {
			break
		}

		tx := txs.At(i)
		applyResult, transferError, appError := c.Executor.ApplyTx(tx, feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			c.handleTransactionError(returnStruct, tx, transferError)
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

func (c *TransactionsContext) handleTransactionError(result ExecuteTxsResult, tx models.GenericTransaction, err error) {
	if errors.Is(err, applier.ErrNonceTooHigh) {
		result.AddSkippedTx(tx)
		return
	}

	logAndSaveTransactionError(c.storage, tx, err)
	result.AddInvalidTx(tx)
	c.txErrorsToStore = append(c.txErrorsToStore, TransactionError{
		Hash:         tx.GetBase().Hash,
		ErrorMessage: err.Error(),
	})
}
