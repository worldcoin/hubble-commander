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
	targetCommitmentProof, err := c.proverCtx.TargetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	if batch.Type == batchtype.Transfer {
		err = c.client.DisputeTransitionTransfer(&batch.ID, batch.Hash, previousCommitmentProof, targetCommitmentProof, merkleProofs)
	} else {
		err = c.client.DisputeTransitionCreate2Transfer(&batch.ID, batch.Hash, previousCommitmentProof, targetCommitmentProof, merkleProofs)
	}
	return err
}
