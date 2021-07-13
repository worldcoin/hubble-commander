package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
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
	returnStruct.appliedTransfers = make([]models.Transfer, 0, t.cfg.MaxTxsPerCommitment)

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
			logAndSaveTransactionError(t.Storage, &transfer.TransactionBase, transferError)
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

func (t *TransactionExecutor) ApplyTransfersForSync(transfers []models.Transfer, feeReceiver *FeeReceiver) (
	appliedTransfers []models.Transfer,
	err error,
) {
	numTransfers := len(transfers)
	if numTransfers == 0 {
		return []models.Transfer{}, nil
	}

	stateChangeProofs := make([]models.StateMerkleProof, 0, 2*numTransfers)

	appliedTransfers = make([]models.Transfer, 0, numTransfers)
	combinedFee := models.MakeUint256(0)

	for i := range transfers {
		transfer := &transfers[i]
		synced, transferError, appError := t.ApplyTransferForSync(transfer, feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		stateChangeProofs = append(
			stateChangeProofs,
			synced.SenderStateProof,
			synced.ReceiverStateProof,
		)
		if transferError != nil {
			return nil, NewDisputableTransferError(transferError, stateChangeProofs)
		}
		appliedTransfers = append(appliedTransfers, *synced.Transfer)
		combinedFee = *combinedFee.Add(&synced.Transfer.Fee)
	}

	if combinedFee.CmpN(0) > 0 {
		err = t.ApplyFee(feeReceiver.StateID, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return appliedTransfers, nil
}
