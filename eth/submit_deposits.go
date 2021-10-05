package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *Client) SubmitDeposits(previous *models.CommitmentInclusionProof, proof *models.SubtreeVacancyProof) (
	*types.Transaction,
	error,
) {
	transaction, err := c.rollup().
		WithValue(*c.config.StakeAmount.ToBig()).
		SubmitDeposits(
			*commitmentProofToCalldata(previous),
			*subtreeVacancyProofToCalldata(proof),
		)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (c *Client) SubmitDepositsAndWait(previous *models.CommitmentInclusionProof, proof *models.SubtreeVacancyProof) (
	*models.Batch, error,
) {
	return c.submitBatchAndWait(func() (*types.Transaction, error) {
		return c.SubmitDeposits(previous, proof)
	})
}

func subtreeVacancyProofToCalldata(proof *models.SubtreeVacancyProof) *rollup.TypesSubtreeVacancyProof {
	return &rollup.TypesSubtreeVacancyProof{
		PathAtDepth: new(big.Int).SetUint64(uint64(proof.PathAtDepth)),
		Witness:     proof.Witness.Bytes(),
	}
}