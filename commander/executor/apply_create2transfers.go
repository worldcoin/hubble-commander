package executor

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
)

type AppliedC2Transfers struct {
	appliedTransfers []models.Create2Transfer
	invalidTransfers []models.Create2Transfer
	addedPubKeyIDs   []uint32
}

func (c *RollupContext) ApplyCreate2Transfers(
	transfers []models.Create2Transfer,
	maxApplied uint32,
	feeReceiver *FeeReceiver,
) (*AppliedC2Transfers, error) {
	if len(transfers) == 0 {
		return &AppliedC2Transfers{}, nil
	}

	returnStruct := &AppliedC2Transfers{}
	returnStruct.appliedTransfers = make([]models.Create2Transfer, 0, c.cfg.MaxTxsPerCommitment)
	returnStruct.addedPubKeyIDs = make([]uint32, 0, c.cfg.MaxTxsPerCommitment)

	combinedFee := models.NewUint256(0)

	for i := range transfers {
		if uint32(len(returnStruct.appliedTransfers)) == maxApplied {
			break
		}

		transfer := &transfers[i]
		applyResult, transferError, appError := c.ApplyCreate2Transfer(transfer, feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(c.storage, applyResult.AppliedTx().GetBase(), transferError)
			returnStruct.invalidTransfers = append(returnStruct.invalidTransfers, *applyResult.AppliedTx().ToCreate2Transfer())
			continue
		}

		returnStruct.appliedTransfers = append(returnStruct.appliedTransfers, *applyResult.AppliedTx().ToCreate2Transfer())
		returnStruct.addedPubKeyIDs = append(returnStruct.addedPubKeyIDs, applyResult.AddedPubKeyID())
		//TODO: change GetFee func to return pointer
		*combinedFee = *combinedFee.Add(&applyResult.AppliedTx().GetBase().Fee)
	}

	if len(returnStruct.appliedTransfers) > 0 {
		_, err := c.ApplyFee(feeReceiver.StateID, *combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}

func (c *ExecutionContext) ApplyCreate2TransfersForSync(
	transfers []models.Create2Transfer,
	pubKeyIDs []uint32,
	feeReceiverStateID uint32,
) ([]models.Create2Transfer, []models.StateMerkleProof, error) {
	transfersLen := len(transfers)
	if transfersLen != len(pubKeyIDs) {
		return nil, nil, applier.ErrInvalidSlicesLength
	}

	appliedTransfers := make([]models.Create2Transfer, 0, transfersLen)
	stateChangeProofs := make([]models.StateMerkleProof, 0, 2*transfersLen+1)
	combinedFee := models.NewUint256(0)

	tokenID, err := c.getCommitmentTokenID(models.Create2TransferArray(transfers), feeReceiverStateID)
	if err != nil {
		return nil, nil, err
	}

	for i := range transfers {
		transfer := &transfers[i]

		synced, transferError, appError := c.ApplyCreate2TransferForSync(transfer, pubKeyIDs[i], *tokenID)
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
