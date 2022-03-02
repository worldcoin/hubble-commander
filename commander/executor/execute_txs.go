package executor

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/mempool"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (c *TxsContext) ExecuteTxs(txMempool *mempool.TxMempool, feeReceiver *FeeReceiver) (ExecuteTxsResult, error) {
	returnStruct := c.Executor.NewExecuteTxsResult(c.cfg.MaxTxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	for tx := c.heap.Peek(); tx != nil; tx = c.heap.Peek() {
		if returnStruct.AppliedTxs().Len() == int(c.cfg.MaxTxsPerCommitment) {
			break
		}

		applyResult, txError, appError := c.Executor.ApplyTx(tx, feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if txError != nil {
			c.handleTxError(txMempool, returnStruct, tx, txError)
			c.heap.Pop()
			continue
		}

		err := c.Executor.AddPendingAccount(applyResult)
		if err != nil {
			return nil, err
		}

		err = c.updateHeap(txMempool, tx)
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

func (c *TxsContext) updateHeap(txMempool *mempool.TxMempool, tx models.GenericTransaction) error {
	nextTx, err := txMempool.GetNextExecutableTx(txtype.TransactionType(c.BatchType), tx.GetFromStateID())
	if err != nil {
		return err
	}
	if nextTx != nil {
		c.heap.Replace(nextTx)
		return nil
	}

	c.heap.Pop()
	return nil
}

func (c *TxsContext) handleTxError(txMempool *mempool.TxMempool, result ExecuteTxsResult, tx models.GenericTransaction, err error) {
	if errors.Is(err, applier.ErrNonceTooHigh) {
		panic("got ErrNonceTooHigh in ExecuteTxs; this should never happen")
	}
	removeErr := txMempool.RemoveFailedTx(tx.GetFromStateID())
	if removeErr != nil {
		panic(removeErr) // should never happen
	}

	log.WithField("txHash", tx.GetBase().Hash.String()).
		Errorf("%s failed: %s", tx.Type().String(), err)
	result.AddInvalidTx(tx)
	c.txErrorsToStore = append(c.txErrorsToStore, models.TxError{
		TxHash:        tx.GetBase().Hash,
		SenderStateID: tx.GetFromStateID(),
		ErrorMessage:  err.Error(),
	})
}
