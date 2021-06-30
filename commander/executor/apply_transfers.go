package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
)

type AppliedTransfers struct {
	appliedTransfers []models.Transfer
	invalidTransfers []models.Transfer
}

func (t *TransactionExecutor) ApplyTransfers(
	transfers []models.Transfer,
	maxApplied uint32,
	feeReceiver *FeeReceiver,
) (*AppliedTransfers, error) {
	if len(transfers) == 0 {
		return &AppliedTransfers{}, nil
	}

	returnStruct := &AppliedTransfers{}
	returnStruct.appliedTransfers = make([]models.Transfer, 0, t.cfg.TxsPerCommitment)

	combinedFee := models.MakeUint256(0)

	for i := range transfers {
		if len(returnStruct.appliedTransfers) == int(maxApplied) {
			break
		}

		transfer := &transfers[i]
		transferError, appError := t.ApplyTransfer(transfer, feeReceiver.TokenID)
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
	}

	if len(returnStruct.appliedTransfers) > 0 {
		err := t.ApplyFee(feeReceiver.StateID, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}

func (t *TransactionExecutor) ApplyTransfersForSync(transfers []models.Transfer, feeReceiver *FeeReceiver) (appliedTransfers []models.Transfer, err error) {
	if len(transfers) == 0 {
		return []models.Transfer{}, nil // TODO-AFS check if there can be commitment without transfers
	}

	appliedTransfers = make([]models.Transfer, 0, t.cfg.TxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	for i := range transfers {
		transfer := &transfers[i]
		synced, transferError, appError := t.ApplyTransferForSync(transfer, feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			return nil, errors.Errorf("invalid transfer") // TODO-AFS return proper error
		}

		appliedTransfers = append(appliedTransfers, *synced.transfer.(*models.Transfer))
		combinedFee = *combinedFee.Add(&transfer.Fee)
	}

	if combinedFee.CmpN(0) > 0 {
		err = t.ApplyFee(feeReceiver.StateID, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return appliedTransfers, nil
}
