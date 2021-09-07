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
	commitmentID *models.CommitmentID,
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

	combinedSignature, err := CombineSignatures(models.Create2TransferArray(appliedTransfers), domain)
	if err != nil {
		return nil, err
	}

	commitment, err := t.createCommitment(
		commitmentID,
		txtype.Create2Transfer,
		feeReceiverStateID,
		serializedTxs,
		combinedSignature,
	)
	if err != nil {
		return nil, err
	}

	err = t.storage.MarkCreate2TransfersAsIncluded(appliedTransfers, commitmentID)
	if err != nil {
		return nil, err
	}

	return commitment, nil
}
