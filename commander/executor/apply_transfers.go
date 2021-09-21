package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (c *RollupContext) ApplyTxs(
	txs models.GenericTransactionArray,
	maxApplied uint32,
	feeReceiver *FeeReceiver,
) (ApplyTxsResult, error) {
	if txs.Len() == 0 {
		return c.Executor.NewApplyTxsResult(0), nil
	}

	returnStruct := c.Executor.NewApplyTxsResult(c.cfg.MaxTxsPerCommitment)
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
			logAndSaveTransactionError(c.storage, applyResult.AppliedTx(), transferError)
			returnStruct.AddInvalidTx(applyResult.AppliedTx())
			continue
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
