package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func (t *TransactionExecutor) buildC2TCommitment(
	appliedTransfers []models.Create2Transfer,
	addedPubKeyIDs []uint32,
	feeReceiverStateID uint32,
	domain *bls.Domain,
) (
	*models.Commitment,
	error,
) {
	serializedTxs, err := encoder.SerializeCreate2Transfers(appliedTransfers, addedPubKeyIDs)
	if err != nil {
		return nil, err
	}

	combinedSignature, err := combineCreate2TransferSignatures(appliedTransfers, domain)
	if err != nil {
		return nil, err
	}

	commitment, err := t.createAndStoreCommitment(
		txtype.Create2Transfer,
		feeReceiverStateID,
		serializedTxs,
		combinedSignature,
	)
	if err != nil {
		return nil, err
	}

	err = t.markCreate2TransfersAsIncluded(appliedTransfers, commitment.ID)
	if err != nil {
		return nil, err
	}

	return commitment, nil
}
