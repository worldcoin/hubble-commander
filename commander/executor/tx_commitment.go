package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
)

func (c *RollupContext) BuildCommitment(
	executeResult ExecuteTxsForCommitmentResult,
	commitmentID *models.CommitmentID,
	feeReceiverStateID uint32,
) (*models.Commitment, error) {
	serializedTxs, err := c.Executor.SerializeTxs(executeResult)
	if err != nil {
		return nil, err
	}

	domain, err := c.client.GetDomain()
	if err != nil {
		return nil, err
	}

	combinedSignature, err := CombineSignatures(executeResult.AppliedTxs(), domain)
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

	err = c.Executor.MarkTxsAsIncluded(executeResult.AppliedTxs(), commitmentID)
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