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
