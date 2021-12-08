package disputer

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

func (c *Context) DisputeTransition(
	batch *eth.DecodedTxBatch,
	commitmentIndex int,
	merkleProofs []models.StateMerkleProof,
) error {
	previousCommitmentProof, err := c.proverCtx.PreviousCommitmentInclusionProof(batch, commitmentIndex-1)
	if err != nil {
		return err
	}

	switch batch.Type {
	case batchtype.Transfer:
		err = c.disputeTransitionTransfer(batch, commitmentIndex, merkleProofs, previousCommitmentProof)
	case batchtype.Create2Transfer:
		err = c.disputeTransitionC2T(batch, commitmentIndex, merkleProofs, previousCommitmentProof)
	case batchtype.MassMigration:
		err = c.disputeTransitionMM(batch, commitmentIndex, merkleProofs, previousCommitmentProof)
	}
	return err
}

func (c *Context) disputeTransitionTransfer(
	batch *eth.DecodedTxBatch,
	commitmentIndex int,
	merkleProofs []models.StateMerkleProof,
	previousCommitmentProof *models.CommitmentInclusionProof,
) error {
	targetCommitmentProof, err := c.proverCtx.TargetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return c.client.DisputeTransitionTransfer(&batch.ID, &batch.Hash, previousCommitmentProof, targetCommitmentProof, merkleProofs)
}

func (c *Context) disputeTransitionC2T(
	batch *eth.DecodedTxBatch,
	commitmentIndex int,
	merkleProofs []models.StateMerkleProof,
	previousCommitmentProof *models.CommitmentInclusionProof,
) error {
	targetCommitmentProof, err := c.proverCtx.TargetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return c.client.DisputeTransitionCreate2Transfer(&batch.ID, &batch.Hash, previousCommitmentProof, targetCommitmentProof, merkleProofs)
}

func (c *Context) disputeTransitionMM(
	batch *eth.DecodedTxBatch,
	commitmentIndex int,
	merkleProofs []models.StateMerkleProof,
	previousCommitmentProof *models.CommitmentInclusionProof,
) error {
	targetCommitmentProof, err := c.proverCtx.TargetMMCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}
	return c.client.DisputeTransitionMassMigration(&batch.ID, &batch.Hash, previousCommitmentProof, targetCommitmentProof, merkleProofs)
}
