package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
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

	combinedSignature, err := CombineSignatures(models.Create2TransferArray(appliedTransfers), domain)
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

	err = t.markCreate2TransfersAsIncluded(appliedTransfers, commitment.IndexInBatch)
	if err != nil {
		return nil, err
	}

	return commitment, nil
}

func (t *TransactionExecutor) markCreate2TransfersAsIncluded(transfers []models.Create2Transfer, commitmentID int32) error {
	hashes := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		hashes = append(hashes, transfers[i].Hash)

		err := t.storage.SetCreate2TransferToStateID(transfers[i].Hash, *transfers[i].ToStateID)
		if err != nil {
			return err
		}
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, &commitmentID)
}
