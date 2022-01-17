package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

func (c *DepositsContext) newCommitment(
	batchID models.Uint256,
	depositSubtree *models.PendingDepositSubtree,
) (*models.DepositCommitment, error) {
	stateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return nil, err
	}

	return &models.DepositCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      batchID,
				IndexInBatch: 0,
			},
			Type:          batchtype.Deposit,
			PostStateRoot: *stateRoot,
		},
		SubtreeID:   depositSubtree.ID,
		SubtreeRoot: depositSubtree.Root,
		Deposits:    depositSubtree.Deposits,
	}, nil
}

func (c *TxsContext) NextCommitmentID() (*models.CommitmentID, error) {
	nextBatchID, err := c.storage.GetNextBatchID()
	if err != nil {
		return nil, err
	}
	return &models.CommitmentID{BatchID: *nextBatchID}, nil
}
