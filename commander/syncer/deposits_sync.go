package syncer

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

func (c *DepositsContext) SyncBatch(remoteBatch eth.DecodedBatch) error {
	batch := remoteBatch.ToDecodedDepositBatch()

	startStateID, depositSubtree, err := c.getDepositSubtree(batch)
	if err != nil {
		return err
	}

	err = c.applier.ApplyDeposits(startStateID, depositSubtree.Deposits)
	if err != nil {
		return err
	}

	err = c.syncCommitment(batch.ID, depositSubtree)
	if err != nil {
		return err
	}

	return c.storage.AddBatch(batch.GetBatch())
}

func (c *DepositsContext) getDepositSubtree(batch *eth.DecodedDepositBatch) (uint32, *models.PendingDepositSubTree, error) {
	subtreeDepth, err := c.client.GetMaxSubTreeDepthParam()
	if err != nil {
		return 0, nil, err
	}
	startStateID := batch.PathAtDepth << *subtreeDepth

	depositSubtree, err := c.storage.GetFirstPendingDepositSubTree()
	if err != nil {
		return 0, nil, err
	}
	return startStateID, depositSubtree, nil
}

func (c *DepositsContext) syncCommitment(batchID models.Uint256, depositSubtree *models.PendingDepositSubTree) error {
	err := c.storage.DeletePendingDepositSubTrees(depositSubtree.ID)
	if err != nil {
		return err
	}

	commitment, err := c.newCommitment(batchID, depositSubtree)
	if err != nil {
		return err
	}

	return c.storage.AddDepositCommitment(commitment)
}

func (c *DepositsContext) newCommitment(
	batchID models.Uint256,
	depositSubtree *models.PendingDepositSubTree,
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
		SubTreeID:   depositSubtree.ID,
		SubTreeRoot: depositSubtree.Root,
		Deposits:    depositSubtree.Deposits,
	}, nil
}
