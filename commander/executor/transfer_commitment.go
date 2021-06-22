package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func (t *TransactionExecutor) prepareTransferCommitment(
	appliedTransfers []models.Transfer,
	feeReceiverStateID uint32,
	domain *bls.Domain,
) (*models.Commitment, error) {
	serializedTxs, err := encoder.SerializeTransfers(appliedTransfers)
	if err != nil {
		return nil, err
	}

	combinedSignature, err := combineTransferSignatures(appliedTransfers, domain)
	if err != nil {
		return nil, err
	}

	commitment, err := t.createAndStoreCommitment(
		txtype.Transfer,
		feeReceiverStateID,
		serializedTxs,
		combinedSignature,
	)
	if err != nil {
		return nil, err
	}

	return commitment, nil
}
