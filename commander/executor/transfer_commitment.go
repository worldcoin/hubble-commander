package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func (t *TransactionExecutor) buildTransferCommitment(
	appliedTransfers []models.Transfer,
	commitmentID *models.CommitmentID,
	feeReceiverStateID uint32,
	domain *bls.Domain,
) (*models.Commitment, error) {
	serializedTxs, err := encoder.SerializeTransfers(appliedTransfers)
	if err != nil {
		return nil, err
	}

	combinedSignature, err := CombineSignatures(models.TransferArray(appliedTransfers), domain)
	if err != nil {
		return nil, err
	}

	commitment, err := t.createCommitment(
		commitmentID,
		txtype.Transfer,
		feeReceiverStateID,
		serializedTxs,
		combinedSignature,
	)
	if err != nil {
		return nil, err
	}

	err = t.storage.MarkTransfersAsIncluded(appliedTransfers, commitmentID)
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
