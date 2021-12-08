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
	batch *eth.DecodedTxBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	switch batch.Type {
	case batchtype.Transfer:
		return c.disputeTransferSignature(batch, commitmentIndex, stateProofs)
	case batchtype.Create2Transfer:
		return c.disputeCreate2TransferSignature(batch, commitmentIndex, stateProofs)
	case batchtype.MassMigration:
		return c.disputeMassMigrationSignature(batch, commitmentIndex, stateProofs)
	case batchtype.Genesis, batchtype.Deposit:
		return errors.WithStack(ErrUnsupportedBatchType)
	}
	return nil
}

func (c *Context) disputeTransferSignature(
	batch *eth.DecodedTxBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	signatureProof, err := c.proverCtx.SignatureProof(stateProofs)
	if err != nil {
		return err
	}

	targetCommitmentProof, err := c.proverCtx.TargetTransferCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return c.client.DisputeSignatureTransfer(&batch.ID, &batch.Hash, targetCommitmentProof, signatureProof)
}

func (c *Context) disputeCreate2TransferSignature(
	batch *eth.DecodedTxBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	signatureProof, err := c.proverCtx.SignatureProofWithReceiver(batch.Commitments[commitmentIndex].ToDecodedCommitment(), stateProofs)
	if err != nil {
		return err
	}

	targetCommitmentProof, err := c.proverCtx.TargetTransferCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return c.client.DisputeSignatureCreate2Transfer(&batch.ID, &batch.Hash, targetCommitmentProof, signatureProof)
}

func (c *Context) disputeMassMigrationSignature(
	batch *eth.DecodedTxBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	signatureProof, err := c.proverCtx.SignatureProof(stateProofs)
	if err != nil {
		return err
	}

	targetCommitmentProof, err := c.proverCtx.TargetMMCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return c.client.DisputeSignatureMassMigration(&batch.ID, &batch.Hash, targetCommitmentProof, signatureProof)
}
