package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

// TODO move this function / rename file

func (t *TransactionExecutor) buildTransferCommitment(
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

	err = t.markTransfersAsIncluded(appliedTransfers, commitment.ID)
	if err != nil {
		return nil, err
	}

	return commitment, nil
}
