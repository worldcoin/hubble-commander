package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
)

func (c *TxsContext) BuildCommitment(
	executeResult ExecuteTxsForCommitmentResult,
	commitmentID *models.CommitmentID,
	feeReceiverStateID uint32,
) (models.CommitmentWithTxs, error) {
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

	commitment, err := c.Executor.NewCommitment(
		commitmentID,
		feeReceiverStateID,
		serializedTxs,
		combinedSignature,
	)
	if err != nil {
		return nil, err
	}

	err = c.storage.MarkTransactionsAsIncluded(executeResult.AppliedTxs(), commitmentID)
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
