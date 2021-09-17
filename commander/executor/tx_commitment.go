package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
)

func (c *RollupContext) buildCommitment(
	applyResult ApplyTxsForCommitmentResult,
	commitmentID *models.CommitmentID,
	feeReceiverStateID uint32,
) (*models.Commitment, error) {
	serializedTxs, err := c.Executor.SerializeTxs(applyResult)
	if err != nil {
		return nil, err
	}

	domain, err := c.client.GetDomain()
	if err != nil {
		return nil, err
	}

	combinedSignature, err := CombineSignatures(applyResult.AppliedTxs(), domain)
	if err != nil {
		return nil, err
	}

	commitment, err := c.newCommitment(
		commitmentID,
		c.BatchType,
		feeReceiverStateID,
		serializedTxs,
		combinedSignature,
	)
	if err != nil {
		return nil, err
	}

	err = c.Executor.MarkTxsAsIncluded(applyResult.AppliedTxs(), commitmentID)
	if err != nil {
		return nil, err
	}

	return commitment, nil
}

func CombineSignatures(txs models.GenericTransactionArray, domain *bls.Domain) (*models.Signature, error) {
	signatures := make([]*bls.Signature, 0, txs.Len())
	for i := 0; i < txs.Len(); i++ {
		sig, err := bls.NewSignatureFromBytes(txs.At(i).GetSignature().Bytes(), *domain)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sig)
	}
	return bls.NewAggregatedSignature(signatures).ModelsSignature(), nil
}
