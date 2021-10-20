package disputer

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/pkg/errors"
)

var ErrUnsupportedBatchType = fmt.Errorf("unsupported batch type")

func (c *Context) DisputeSignature(
	batch *eth.DecodedBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	switch batch.Type {
	case batchtype.Transfer:
		return c.disputeTransferSignature(batch, commitmentIndex, stateProofs)
	case batchtype.Create2Transfer:
		return c.disputeCreate2TransferSignature(batch, commitmentIndex, stateProofs)
	case batchtype.Genesis, batchtype.MassMigration, batchtype.Deposit:
		return errors.WithStack(ErrUnsupportedBatchType)
	}
	return nil
}

func (c *Context) disputeTransferSignature(
	batch *eth.DecodedBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	signatureProof, err := c.prooferCtx.SignatureProof(stateProofs)
	if err != nil {
		return err
	}

	targetCommitmentProof, err := c.prooferCtx.TargetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return c.client.DisputeSignatureTransfer(&batch.ID, batch.Hash, targetCommitmentProof, signatureProof)
}

func (c *Context) disputeCreate2TransferSignature(
	batch *eth.DecodedBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	signatureProof, err := c.prooferCtx.SignatureProofWithReceiver(&batch.Commitments[commitmentIndex], stateProofs)
	if err != nil {
		return err
	}

	targetCommitmentProof, err := c.prooferCtx.TargetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return c.client.DisputeSignatureCreate2Transfer(&batch.ID, batch.Hash, targetCommitmentProof, signatureProof)
}
