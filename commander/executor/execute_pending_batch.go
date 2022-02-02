package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

func (c *TxsContext) ExecutePendingBatch(batch *dto.PendingBatch) error {
	err := c.storage.AddBatch(&models.Batch{
		ID:              batch.ID,
		Type:            batch.Type,
		TransactionHash: batch.TransactionHash,
	})
	if err != nil {
		return err
	}

	for i := range batch.Commitments {
		err = c.storage.AddCommitment(batch.Commitments[i].Commitment)
		if err != nil {
			return err
		}

		err = c.storage.BatchAddTransaction(batch.Commitments[i].Transactions)
		if err != nil {
			return err
		}

		var feeReceiver *FeeReceiver
		feeReceiver, err = c.getFeeReceiver(batch.Commitments[i].Commitment)
		if err != nil {
			return err
		}
		//TODO-mig: at least check invalid and skipped txs
		_, err = c.ExecuteTxs(batch.Commitments[i].Transactions, feeReceiver)
		if err != nil {
			return err
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

func (c *DepositsContext) ExecutePendingBatch(batch *dto.PendingBatch) error {
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
