package executor

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (c *TxsContext) ExecuteTxs(mempoolHeap *storage.MempoolHeap, feeReceiver *FeeReceiver) (ExecuteTxsResult, error) {
	returnStruct := c.Executor.NewExecuteTxsResult(c.cfg.MaxTxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	peekFn := mempoolHeap.PeekHighestFeeExecutableTx
	for tx := peekFn(); tx != nil; tx = peekFn() {
		if returnStruct.AppliedTxs().Len() == int(c.cfg.MaxTxsPerCommitment) {
			break
		}

		applyResult, txError, appError := c.Executor.ApplyTx(tx, feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if txError != nil {
			// TODO: should we return an appError here? This is very bad, it
			//       might be better to _not_ continue and mess up the state
			//       any further.
			c.handleTxError(returnStruct, tx, txError)
			continue
		}

		err := c.Executor.AddPendingAccount(applyResult)
		if err != nil {
			return nil, err
		}

		err = mempoolHeap.DropHighestFeeExecutableTx()
		if err != nil {
			return nil, err
		}

		returnStruct.AddApplied(applyResult)
		fee := applyResult.AppliedTx().GetFee()
		combinedFee = *combinedFee.Add(&fee)
	}

	if returnStruct.AppliedTxs().Len() > 0 {
		_, err := c.Applier.ApplyFee(feeReceiver.StateID, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}

// TODO: thread the context down here so we can attach this to the rollup span
func (c *TxsContext) handleTxError(
	result ExecuteTxsResult,
	tx models.GenericTransaction,
	err error,
) {
	if errors.Is(err, applier.ErrNonceTooHigh) {
		panic("got ErrNonceTooHigh in ExecuteTxs; this should never happen")
	}

	// TODO: If this happens we need to scan through the mempool and cascade the
	//       failure to other txns which were relying on our successful execution.

	log.WithFields(log.Fields{
		"tx.Hash":        tx.GetBase().Hash.String(),
		"tx.FromStateID": tx.GetBase().FromStateID,
		"tx.Nonce":       tx.GetBase().Nonce.Uint64(),
		"tx.Type":        tx.Type().String(),
		"errMessage":     err.Error(),
		"err":            err,
	}).Errorf("Unimplemented: failed to batch transaction. State might be inconsistent")
	result.AddInvalidTx(tx)
}
