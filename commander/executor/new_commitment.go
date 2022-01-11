package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

func (c *TxsContext) newCommitment(
	commitmentID *models.CommitmentID,
	feeReceiverStateID uint32,
	serializedTxs []byte,
	combinedSignature *models.Signature,
) (models.CommitmentWithTxs, error) {
	stateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return nil, err
	}

	if c.BatchType == batchtype.MassMigration {
		return &models.MMCommitmentWithTxs{
			MMCommitment: models.MMCommitment{
				CommitmentBase: models.CommitmentBase{
					ID:            *commitmentID,
					Type:          c.BatchType,
					PostStateRoot: *stateRoot,
				},
				FeeReceiver:       feeReceiverStateID,
				CombinedSignature: *combinedSignature,
			},
			Transactions: serializedTxs,
		}, nil
	}

	return &models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				ID:            *commitmentID,
				Type:          c.BatchType,
				PostStateRoot: *stateRoot,
			},
			FeeReceiver:       feeReceiverStateID,
			CombinedSignature: *combinedSignature,
		},
		Transactions: serializedTxs,
	}, nil
}

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
