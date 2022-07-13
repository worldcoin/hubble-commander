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
			// TODO: throw a big error here, if this happens something has
			//       gone terribly wrong, we might even want to take downtime
			// c.handleTxError(txMempool, returnStruct, tx, txError)
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

func (c *TxsContext) handleTxError(result ExecuteTxsResult, tx models.GenericTransaction, err error) {
	if errors.Is(err, applier.ErrNonceTooHigh) {
		panic("got ErrNonceTooHigh in ExecuteTxs; this should never happen")
	}

	// TODO: Why does this happen? What could cause a transaction which we previously
	//       accepted to fail? If this happens we need to scan through the mempool and
	//       cascade the failure to other txns which were relying on our successful
	//       execution.

	log.WithField("txHash", tx.GetBase().Hash.String()).
		Errorf("%s failed: %s", tx.Type().String(), err)
	result.AddInvalidTx(tx)
	c.txErrorsToStore = append(c.txErrorsToStore, models.TxError{
		TxHash:        tx.GetBase().Hash,
		SenderStateID: tx.GetFromStateID(),
		ErrorMessage:  err.Error(),
	})
}
