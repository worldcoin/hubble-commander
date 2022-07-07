package eth

import (
	"context"
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (c *Client) SubmitDeposits(
	batchID *models.Uint256,
	previous *models.CommitmentInclusionProof,
	proof *models.SubtreeVacancyProof,
) (*types.Transaction, error) {
	_, span := clientTracer.Start(context.Background(), "Rollup.SubmitDeposits")
	defer span.End()

	if batchID != nil {
		span.SetAttributes(attribute.String("batchID", batchID.String()))
	}

	tx, err := c.rollup().
		WithValue(c.config.StakeAmount).
		WithGasLimit(*c.config.DepositBatchSubmissionGasLimit).
		SubmitDeposits(
			batchID.ToBig(),
			*commitmentProofToCalldata(previous),
			*subtreeVacancyProofToCalldata(proof),
		)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
		return nil, err
	}
	if tx != nil {
		span.SetAttributes(attribute.String("txHash", tx.Hash().Hex()))
	}
	span.SetStatus(codes.Ok, "")
	return tx, nil
}

func (c *Client) SubmitDepositsAndWait(
	batchID *models.Uint256,
	previous *models.CommitmentInclusionProof,
	proof *models.SubtreeVacancyProof,
) (*models.Batch, error) {
	return c.submitBatchAndWait(func() (*types.Transaction, error) {
		return c.SubmitDeposits(batchID, previous, proof)
	})
}

func subtreeVacancyProofToCalldata(proof *models.SubtreeVacancyProof) *rollup.TypesSubtreeVacancyProof {
	return &rollup.TypesSubtreeVacancyProof{
		PathAtDepth: new(big.Int).SetUint64(uint64(proof.PathAtDepth)),
		Witness:     proof.Witness.Bytes(),
	}
}
