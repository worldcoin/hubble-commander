package executor

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/pkg/errors"
)

var errInvalidPendingBatchTxs = fmt.Errorf("failed to apply all transactions from pending batch")

func (c *TxsContext) ExecutePendingBatch(batch *models.PendingBatch) error {
	err := c.storage.AddBatch(&models.Batch{
		ID:              batch.ID,
		Type:            batch.Type,
		TransactionHash: batch.TransactionHash,
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
		executeTxsResult, err := c.ExecuteTxs(pendingCommitment.Transactions, feeReceiver)
		if err != nil {
			return err
		}
		if executeTxsResult.AppliedTxs().Len() != pendingCommitment.Transactions.Len() {
			return errors.WithStack(errInvalidPendingBatchTxs)
		}
	}
	return nil
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
