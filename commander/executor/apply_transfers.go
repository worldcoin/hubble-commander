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
		receiverLeaf, err := t.storage.StateTree.Leaf(transfer.ToStateID)
		if err != nil {
			return nil, err
		}

		transferError, appError := t.ApplyTransfer(transfer, receiverLeaf, feeReceiver.TokenID)
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
		_, err := t.ApplyFee(feeReceiver.StateID, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}

func (t *TransactionExecutor) ApplyTransfersForSync(transfers []models.Transfer, feeReceiver *FeeReceiver) (
	[]models.Transfer,
	[]models.StateMerkleProof,
	error,
) {
	numTransfers := len(transfers)
	if numTransfers == 0 {
		// TODO state proofs should probably always contain at least fee receiver's proof for ErrInvalidCommitmentStateRoot disputes to work
		return []models.Transfer{}, nil, nil
	}

	appliedTransfers := make([]models.Transfer, 0, numTransfers)
	stateChangeProofs := make([]models.StateMerkleProof, 0, 2*numTransfers)
	combinedFee := models.NewUint256(0)

	tokenID, err := t.getCommitmentTokenID(models.TransferArray(transfers), &feeReceiver.TokenID)
	if err != nil {
		return nil, nil, err
	}

	for i := range transfers {
		transfer := &transfers[i]
		synced, transferError, appError := t.ApplyTransferForSync(transfer, *tokenID)
		if appError != nil {
			return nil, nil, appError
		}
		stateChangeProofs = append(
			stateChangeProofs,
			synced.SenderStateProof,
			synced.ReceiverStateProof,
		)
		if transferError != nil {
			return nil, nil, NewDisputableErrorWithProofs(Transition, transferError.Error(), stateChangeProofs)
		}
		appliedTransfers = append(appliedTransfers, *synced.Transfer)
		*combinedFee = *combinedFee.Add(&synced.Transfer.Fee)
	}

	stateProof, transferError, appError := t.ApplyFeeForSync(feeReceiver, tokenID, combinedFee)
	if appError != nil {
		return nil, nil, appError
	}
	stateChangeProofs = append(stateChangeProofs, *stateProof)
	if transferError != nil {
		return nil, nil, NewDisputableErrorWithProofs(Transition, transferError.Error(), stateChangeProofs)
	}

	return appliedTransfers, stateChangeProofs, nil
}

func (t *TransactionExecutor) getCommitmentTokenID(transfers models.GenericTransactionArray, feeReceiverTokenID *models.Uint256) (*models.Uint256, error) {
	if transfers.Len() == 0 {
		return feeReceiverTokenID, nil
	}
	leaf, err := t.storage.StateTree.Leaf(transfers.At(0).GetFromStateID())
	if err != nil {
		return nil, err
	}
	return &leaf.TokenID, nil
}
