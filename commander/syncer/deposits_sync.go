package syncer

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
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

	err = c.storage.DeletePendingDepositSubtrees(depositSubtree.ID)
	if err != nil {
		return err
	}

	return c.addCommitment(batch.ID, depositSubtree)
}

func (c *DepositsContext) UpdateExistingBatch(batch eth.DecodedBatch, prevStateRoot common.Hash) error {
	return c.storage.UpdateBatch(batch.ToBatch(prevStateRoot))
}

func (c *DepositsContext) getDepositSubtree(batch *eth.DecodedDepositBatch) (uint32, *models.PendingDepositSubtree, error) {
	subtreeDepth, err := c.client.GetMaxSubtreeDepthParam()
	if err != nil {
		return 0, nil, err
	}
	startStateID := batch.PathAtDepth << *subtreeDepth

	depositSubtree, err := c.storage.GetPendingDepositSubtree(batch.SubtreeID)
	if err != nil {
		return 0, nil, err
	}
	return startStateID, depositSubtree, nil
}

func (c *DepositsContext) addCommitment(batchID models.Uint256, depositSubtree *models.PendingDepositSubtree) error {
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
			SubtreeID:   depositSubtree.ID,
			SubtreeRoot: depositSubtree.Root,
			Deposits:    depositSubtree.Deposits,
		},
	)
}
