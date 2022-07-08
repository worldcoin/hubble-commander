package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
	"go.opentelemetry.io/otel/attribute"
)

func (c *Client) SubmitDeposits(
	batchID *models.Uint256,
	previous *models.CommitmentInclusionProof,
	proof *models.SubtreeVacancyProof,
) (*types.Transaction, error) {
	builder := c.rollup().
		WithValue(c.config.StakeAmount).
		WithGasLimit(*c.config.DepositBatchSubmissionGasLimit)

	if batchID != nil {
		builder = builder.WithAttribute(
			attribute.String("batchID", batchID.String()),
		)
	}

	return builder.SubmitDeposits(
		batchID.ToBig(),
		*commitmentProofToCalldata(previous),
		*subtreeVacancyProofToCalldata(proof),
	)
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
