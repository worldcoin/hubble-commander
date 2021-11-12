package syncer

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

func (c *DepositsContext) SyncCommitments(remoteBatch eth.DecodedBatch) error {
	batch := remoteBatch.ToDecodedDepositBatch()
	startStateID, depositSubtree, err := c.getDepositSubtree(batch)
	if err != nil {
		return err
	}

	err = c.applier.ApplyDeposits(startStateID, depositSubtree.Deposits)
	if err != nil {
		return err
	}

	err = c.storage.DeletePendingDepositSubTrees(depositSubtree.ID)
	if err != nil {
		return err
	}

	return c.addCommitment(batch.ID, depositSubtree)
}

func (c *DepositsContext) UpdateExistingBatch(batch eth.DecodedBatch) error {
	return c.storage.UpdateBatch(batch.GetBatch())
}

func (c *DepositsContext) getDepositSubtree(batch *eth.DecodedDepositBatch) (uint32, *models.PendingDepositSubTree, error) {
	subtreeDepth, err := c.client.GetMaxSubTreeDepthParam()
	if err != nil {
		return 0, nil, err
	}
	startStateID := batch.PathAtDepth << *subtreeDepth

	depositSubtree, err := c.storage.GetFirstPendingDepositSubTree() // TODO use SubtreeID from DecodedDepositBatch to query directly with GetPendingDepositSubTree
	if err != nil {
		return 0, nil, err
	}
	return startStateID, depositSubtree, nil
}

func (c *DepositsContext) addCommitment(batchID models.Uint256, depositSubtree *models.PendingDepositSubTree) error {
	stateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return err
	}

	return c.storage.AddDepositCommitment(
		&models.DepositCommitment{
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
		},
	)
}
