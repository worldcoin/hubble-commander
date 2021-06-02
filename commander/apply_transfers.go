package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type AppliedTransfers struct {
	appliedTransfers   []models.Transfer
	invalidTransfers   []models.Transfer
	feeReceiverStateID *uint32
}

func (t *transactionExecutor) ApplyTransfers(
	transfers []models.Transfer,
) (returnStruct *AppliedTransfers, err error) {
	if len(transfers) == 0 {
		return
	}

	returnStruct.appliedTransfers = make([]models.Transfer, 0, t.cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	senderLeaf, err := t.storage.GetStateLeaf(transfers[0].FromStateID)
	if err != nil {
		return nil, err
	}

	commitmentTokenIndex := senderLeaf.TokenIndex

	for i := range transfers {
		transfer := &transfers[i]
		transferError, appError := ApplyTransfer(t.storage, transfer, commitmentTokenIndex)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(t.storage, &transfer.TransactionBase, transferError)
			returnStruct.invalidTransfers = append(returnStruct.invalidTransfers, *transfer)
			continue
		}

		returnStruct.appliedTransfers = append(returnStruct.appliedTransfers, *transfer)
		combinedFee = *combinedFee.Add(&transfer.Fee)

		if uint32(len(returnStruct.appliedTransfers)) == t.cfg.TxsPerCommitment {
			break
		}
	}

	if len(returnStruct.appliedTransfers) > 0 {
		returnStruct.feeReceiverStateID, err = t.ApplyFee(commitmentTokenIndex, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}
