package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type AppliedTransfers struct {
	appliedTransfers   []models.Transfer
	invalidTransfers   []models.Transfer
	feeReceiverStateID *uint32
}

func (t *TransactionExecutor) ApplyTransfers(
	transfers []models.Transfer,
	maxApplied uint32,
	isSyncing bool,
) (*AppliedTransfers, error) {
	if len(transfers) == 0 {
		return &AppliedTransfers{}, nil
	}

	commitmentTokenIndex, err := t.getTokenIndex(transfers[0].FromStateID)
	if err != nil {
		return nil, err
	}

	returnStruct := &AppliedTransfers{}
	returnStruct.appliedTransfers = make([]models.Transfer, 0, t.cfg.TxsPerCommitment)

	combinedFee := models.MakeUint256(0)

	for i := range transfers {
		if len(returnStruct.appliedTransfers) == int(maxApplied) {
			break
		}

		transfer := &transfers[i]
		transferError, appError := t.ApplyTransfer(transfer, *commitmentTokenIndex)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(t.storage, &transfer.TransactionBase, transferError)
			returnStruct.invalidTransfers = append(returnStruct.invalidTransfers, *transfer)
			if isSyncing {
				return returnStruct, nil
			}
			continue
		}

		returnStruct.appliedTransfers = append(returnStruct.appliedTransfers, *transfer)
		combinedFee = *combinedFee.Add(&transfer.Fee)
	}

	if len(returnStruct.appliedTransfers) > 0 {
		returnStruct.feeReceiverStateID, err = t.ApplyFee(*commitmentTokenIndex, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}
