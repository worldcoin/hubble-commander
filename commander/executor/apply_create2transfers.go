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
	returnStruct.appliedTransfers = make([]models.Create2Transfer, 0, t.cfg.TxsPerCommitment)
	returnStruct.addedPubKeyIDs = make([]uint32, 0, t.cfg.TxsPerCommitment)

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

		_, err = t.handleApplyC2T(transfer, *pubKeyID, returnStruct, combinedFee, &feeReceiver.TokenID)
		if err != nil {
			return nil, err
		}
	}

	if len(returnStruct.appliedTransfers) > 0 {
		err = t.ApplyFee(feeReceiver.StateID, *combinedFee)
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
) (*AppliedC2Transfers, error) {
	if len(transfers) == 0 {
		return nil, nil
	}
	if len(transfers) != len(pubKeyIDs) {
		return nil, ErrInvalidSliceLength
	}

	returnStruct := &AppliedC2Transfers{}

	returnStruct.appliedTransfers = make([]models.Create2Transfer, 0, t.cfg.TxsPerCommitment)
	combinedFee := models.NewUint256(0)

	for i := range transfers {
		transfer := &transfers[i]

		transferError, appError := t.handleApplyC2T(transfer, pubKeyIDs[i], returnStruct, combinedFee, &feeReceiver.TokenID)
		if appError != nil {
			return nil, appError
		}
		if transferError != nil {
			return returnStruct, nil
		}
	}

	if len(returnStruct.appliedTransfers) > 0 {
		err := t.ApplyFee(feeReceiver.StateID, *combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
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

func (t *TransactionExecutor) handleApplyC2T(
	transfer *models.Create2Transfer,
	pubKeyID uint32,
	appliedTransfers *AppliedC2Transfers,
	combinedFee, tokenID *models.Uint256,
) (create2TransferError, appError error) {
	transferError, appError := t.ApplyCreate2Transfer(transfer, pubKeyID, *tokenID)
	if appError != nil {
		return nil, appError
	}
	if transferError != nil {
		logAndSaveTransactionError(t.storage, &transfer.TransactionBase, transferError)
		appliedTransfers.invalidTransfers = append(appliedTransfers.invalidTransfers, *transfer)
		return transferError, nil
	}

	appliedTransfers.appliedTransfers = append(appliedTransfers.appliedTransfers, *transfer)
	appliedTransfers.addedPubKeyIDs = append(appliedTransfers.addedPubKeyIDs, pubKeyID)
	*combinedFee = *combinedFee.Add(&transfer.Fee)
	return nil, nil
}
