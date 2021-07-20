package executor

import (
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

func (t *TransactionExecutor) ApplyCreate2Transfers(
	transfers []models.Create2Transfer,
	maxApplied uint32,
	feeReceiver *FeeReceiver,
) (*AppliedC2Transfers, error) {
	if len(transfers) == 0 {
		return &AppliedC2Transfers{}, nil
	}

	syncedBlock, err := t.storage.GetSyncedBlock(t.client.ChainState.ChainID)
	if err != nil {
		return nil, err
	}
	events, unsubscribe, err := t.client.WatchRegistrations(&bind.WatchOpts{
		Start: syncedBlock,
	})
	if err != nil {
		return nil, err
	}
	defer unsubscribe()

	returnStruct := &AppliedC2Transfers{}
	returnStruct.appliedTransfers = make([]models.Create2Transfer, 0, t.cfg.MaxTxsPerCommitment)
	returnStruct.addedPubKeyIDs = make([]uint32, 0, t.cfg.MaxTxsPerCommitment)

	combinedFee := models.NewUint256(0)

	for i := range transfers {
		if uint32(len(returnStruct.appliedTransfers)) == maxApplied {
			break
		}

		transfer := &transfers[i]
		var pubKeyID *uint32
		pubKeyID, err = t.getOrRegisterPubKeyID(events, transfer, feeReceiver.TokenID)
		if err != nil {
			return nil, err
		}

		appliedTransfer, transferError, appError := t.ApplyCreate2Transfer(transfer, *pubKeyID, feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			logAndSaveTransactionError(t.storage, &appliedTransfer.TransactionBase, transferError)
			returnStruct.invalidTransfers = append(returnStruct.invalidTransfers, *appliedTransfer)
			continue
		}

		returnStruct.appliedTransfers = append(returnStruct.appliedTransfers, *appliedTransfer)
		returnStruct.addedPubKeyIDs = append(returnStruct.addedPubKeyIDs, *pubKeyID)
		*combinedFee = *combinedFee.Add(&appliedTransfer.Fee)
	}

	if len(returnStruct.appliedTransfers) > 0 {
		_, err = t.ApplyFee(feeReceiver.StateID, *combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}

func (t *TransactionExecutor) ApplyCreate2TransfersForSync(
	transfers []models.Create2Transfer,
	pubKeyIDs []uint32,
	feeReceiver *FeeReceiver,
) ([]models.Create2Transfer, []models.StateMerkleProof, error) {
	if len(transfers) == 0 {
		return nil, nil, nil
	}
	if len(transfers) != len(pubKeyIDs) {
		return nil, nil, ErrInvalidSliceLength
	}

	appliedTransfers := make([]models.Create2Transfer, 0, t.cfg.MaxTxsPerCommitment)
	stateChangeProofs := make([]models.StateMerkleProof, 0, 2*len(transfers))
	combinedFee := models.NewUint256(0)

	for i := range transfers {
		transfer := &transfers[i]

		synced, transferError, appError := t.ApplyCreate2TransferForSync(transfer, pubKeyIDs[i], feeReceiver.TokenID)
		if appError != nil {
			return nil, nil, appError
		}
		stateChangeProofs = append(
			stateChangeProofs,
			synced.SenderStateProof,
			synced.ReceiverStateProof,
		)
		if transferError != nil {
			return nil, nil, NewDisputableTransferError(transferError, stateChangeProofs)
		}

		appliedTransfers = append(appliedTransfers, *synced.Transfer)
		*combinedFee = *combinedFee.Add(&synced.Transfer.Fee)
	}

	if len(appliedTransfers) > 0 {
		stateProof, err := t.ApplyFee(feeReceiver.StateID, *combinedFee)
		if err != nil {
			return nil, nil, err
		}
		stateChangeProofs = append(stateChangeProofs, *stateProof)
	}

	return appliedTransfers, stateChangeProofs, nil
}

func (t *TransactionExecutor) getOrRegisterPubKeyID(
	events chan *accountregistry.AccountRegistrySinglePubkeyRegistered,
	transfer *models.Create2Transfer,
	tokenID models.Uint256,
) (*uint32, error) {
	pubKeyID, err := t.storage.GetUnusedPubKeyID(&transfer.ToPublicKey, &tokenID)
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	} else if st.IsNotFoundError(err) {
		return t.client.RegisterAccount(&transfer.ToPublicKey, events)
	}
	return pubKeyID, nil
}
