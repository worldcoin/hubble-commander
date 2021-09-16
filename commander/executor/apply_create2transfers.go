package executor

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

	syncedBlock, err := c.storage.GetSyncedBlock()
	if err != nil {
		return nil, err
	}
	events, unsubscribe, err := c.client.WatchRegistrations(&bind.WatchOpts{
		Start: syncedBlock,
	})
	if err != nil {
		return nil, err
	}
	defer unsubscribe()

	returnStruct := &AppliedC2Transfers{}
	returnStruct.appliedTransfers = make([]models.Create2Transfer, 0, c.cfg.MaxTxsPerCommitment)
	returnStruct.addedPubKeyIDs = make([]uint32, 0, c.cfg.MaxTxsPerCommitment)

	combinedFee := models.NewUint256(0)

	for i := range transfers {
		if uint32(len(returnStruct.appliedTransfers)) == maxApplied {
			break
		}

		transfer := &transfers[i]
		var pubKeyID *uint32
		pubKeyID, err = c.getOrRegisterPubKeyID(events, transfer, feeReceiver.TokenID)
		if err != nil {
			return nil, err
		}

		appliedTransfer, transferError, appError := c.ApplyCreate2Transfer(transfer, *pubKeyID, feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(c.storage, &appliedTransfer.TransactionBase, transferError)
			returnStruct.invalidTransfers = append(returnStruct.invalidTransfers, *appliedTransfer)
			continue
		}

		returnStruct.appliedTransfers = append(returnStruct.appliedTransfers, *appliedTransfer)
		returnStruct.addedPubKeyIDs = append(returnStruct.addedPubKeyIDs, *pubKeyID)
		*combinedFee = *combinedFee.Add(&appliedTransfer.Fee)
	}

	if len(returnStruct.appliedTransfers) > 0 {
		_, err = c.ApplyFee(feeReceiver.StateID, *combinedFee)
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

func (c *ExecutionContext) getOrRegisterPubKeyID(
	events chan *accountregistry.AccountRegistrySinglePubkeyRegistered,
	transfer *models.Create2Transfer,
	tokenID models.Uint256,
) (*uint32, error) {
	pubKeyID, err := c.storage.GetUnusedPubKeyID(&transfer.ToPublicKey, &tokenID)
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	} else if st.IsNotFoundError(err) {
		return c.client.RegisterAccount(&transfer.ToPublicKey, events)
	}
	return pubKeyID, nil
}
