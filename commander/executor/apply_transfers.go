package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type AppliedTransfers struct {
	appliedTransfers models.GenericTransactionArray
	invalidTransfers models.GenericTransactionArray
}

func (a *AppliedTransfers) AppliedTransfers() models.GenericTransactionArray {
	return a.appliedTransfers
}

func (a *AppliedTransfers) InvalidTransfers() models.GenericTransactionArray {
	return a.invalidTransfers
}

func (a *AppliedTransfers) AddedPubKeyIDs() models.GenericTransactionArray {
	return models.TransferArray{}
}

func (c *RollupContext) ApplyTransfers(
	transfers models.GenericTransactionArray,
	maxApplied uint32,
	feeReceiver *FeeReceiver,
) (ApplyTxsResult, error) {
	returnStruct := &AppliedTransfers{
		appliedTransfers: c.Executor.makeTransactionArray(0, 0),
		invalidTransfers: c.Executor.makeTransactionArray(0, 0),
	}

	if transfers.Len() == 0 {
		return returnStruct, nil
	}

	returnStruct.appliedTransfers = c.Executor.makeTransactionArray(0, c.cfg.MaxTxsPerCommitment)
	combinedFee := models.MakeUint256(0)

	for i := 0; i < transfers.Len(); i++ {
		if returnStruct.appliedTransfers.Len() == int(maxApplied) {
			break
		}

		transfer := transfers.At(i)
		receiverLeaf, err := c.Executor.beforeApplyTransaction(transfer)
		if err != nil {
			return nil, err
		}

		transferError, appError := c.ApplyTransfer(transfer, receiverLeaf, feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(c.storage, transfer.GetBase(), transferError)
			returnStruct.invalidTransfers = returnStruct.invalidTransfers.AppendOne(transfer)
			continue
		}

		returnStruct.appliedTransfers = returnStruct.appliedTransfers.AppendOne(transfer)
		combinedFee = *combinedFee.Add(&transfer.GetBase().Fee)
	}

	if returnStruct.appliedTransfers.Len() > 0 {
		_, err := c.ApplyFee(feeReceiver.StateID, combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}

func (c *ExecutionContext) ApplyTransfersForSync(transfers []models.Transfer, feeReceiverStateID uint32) (
	[]models.Transfer,
	[]models.StateMerkleProof,
	error,
) {
	transfersLen := len(transfers)
	appliedTransfers := make([]models.Transfer, 0, transfersLen)
	stateChangeProofs := make([]models.StateMerkleProof, 0, 2*transfersLen+1)
	combinedFee := models.NewUint256(0)

	tokenID, err := c.getCommitmentTokenID(models.TransferArray(transfers), feeReceiverStateID)
	if err != nil {
		return nil, nil, err
	}

	for i := range transfers {
		transfer := &transfers[i]
		synced, transferError, appError := c.ApplyTransferForSync(transfer, *tokenID)
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

	stateProof, commitmentError, appError := c.ApplyFeeForSync(feeReceiverStateID, tokenID, combinedFee)
	if appError != nil {
		return nil, nil, appError
	}
	stateChangeProofs = append(stateChangeProofs, *stateProof)
	if commitmentError != nil {
		return nil, nil, NewDisputableErrorWithProofs(Transition, commitmentError.Error(), stateChangeProofs)
	}

	return appliedTransfers, stateChangeProofs, nil
}

func (c *ExecutionContext) getCommitmentTokenID(transfers models.GenericTransactionArray, feeReceiverStateID uint32) (
	tokenID *models.Uint256,
	err error,
) {
	var leaf *models.StateLeaf

	if transfers.Len() > 0 {
		leaf, err = c.storage.StateTree.LeafOrEmpty(transfers.At(0).GetFromStateID())
	} else {
		leaf, err = c.storage.StateTree.LeafOrEmpty(feeReceiverStateID)
	}
	if err != nil {
		return nil, err
	}

	return &leaf.TokenID, nil
}
