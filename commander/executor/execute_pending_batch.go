package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

func (c *TxsContext) ExecutePendingBatch(batch *models.PendingBatch) error {
	err := c.storage.AddBatch(&models.Batch{
		ID:              batch.ID,
		Type:            batch.Type,
		TransactionHash: batch.TransactionHash,
		PrevStateRoot:   &batch.PrevStateRoot,
	})
	if err != nil {
		return err
	}

	for i := range batch.Commitments {
		pendingCommitment := batch.Commitments[i]
		err = c.storage.AddCommitment(pendingCommitment.Commitment)
		if err != nil {
			return err
		}

		err = c.storage.BatchAddTransaction(pendingCommitment.Transactions)
		if err != nil {
			return err
		}

		feeReceiver, err := c.getFeeReceiver(pendingCommitment.Commitment)
		if err != nil {
			return err
		}
		_, err = c.executePendingTxs(pendingCommitment.Transactions, feeReceiver)
		if err != nil {
			return err
		}
	}
	return nil
}

// Idea: deduplicate this function and ExecuteTxs by extracting a struct which is a generalised "source" of transactions
func (c *TxsContext) executePendingTxs(txs models.GenericTransactionArray, feeReceiver *FeeReceiver) (ExecuteTxsResult, error) {
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
		applyResult, txError, appError := c.Executor.ApplyTx(tx, feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if txError != nil {
			return nil, txError
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
		_, err := c.Applier.ApplyFee(feeReceiver.StateID, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}

func (c *TxsContext) getFeeReceiver(commitment models.Commitment) (*FeeReceiver, error) {
	var stateID uint32
	switch commitment.GetCommitmentBase().Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		stateID = commitment.ToTxCommitment().FeeReceiver
	case batchtype.Deposit:
		stateID = commitment.ToMMCommitment().Meta.FeeReceiver
	}

	feeReceiver, err := c.storage.StateTree.Leaf(stateID)
	if err != nil {
		return nil, err
	}

	return &FeeReceiver{
		StateID: feeReceiver.StateID,
		TokenID: feeReceiver.TokenID,
	}, nil
}

func (c *DepositsContext) ExecutePendingBatch(batch *models.PendingBatch) error {
	err := c.storage.AddBatch(&models.Batch{
		ID:              batch.ID,
		Type:            batch.Type,
		TransactionHash: batch.TransactionHash,
		PrevStateRoot:   &batch.PrevStateRoot,
	})
	if err != nil {
		return err
	}

	for i := range batch.Commitments {
		commitment := batch.Commitments[i].Commitment.ToDepositCommitment()

		err = c.storage.AddCommitment(batch.Commitments[i].Commitment)
		if err != nil {
			return err
		}

		subtree, err := c.storage.GetPendingDepositSubtree(commitment.SubtreeID)
		if err != nil {
			return err
		}

		_, err = c.executeDeposits(subtree)
		if err != nil {
			return err
		}
	}
	return nil
}
