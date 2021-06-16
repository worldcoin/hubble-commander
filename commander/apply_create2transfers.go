package commander

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type AppliedC2Transfers struct {
	appliedTransfers   []models.Create2Transfer
	invalidTransfers   []models.Create2Transfer
	addedPubKeyIDs     []uint32
	feeReceiverStateID *uint32
}

func (t *transactionExecutor) ApplyCreate2Transfers(
	transfers []models.Create2Transfer,
	maxAppliedTransfers uint64,
) (*AppliedC2Transfers, error) {
	if len(transfers) == 0 {
		return nil, nil
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

	commitmentTokenIndex, err := t.getTokenIndex(transfers[0].FromStateID)
	if err != nil {
		return nil, err
	}
	var pubKeyID *uint32
	var ok bool

	for i := range transfers {
		transfer := &transfers[i]

		pubKeyID, err = t.getOrRegisterPubKeyID(events, transfer, *commitmentTokenIndex)
		if err != nil {
			return nil, err
		}

		ok, err = t.handleApplyC2T(transfer, *pubKeyID, returnStruct, combinedFee, commitmentTokenIndex)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		returnStruct.addedPubKeyIDs = append(returnStruct.addedPubKeyIDs, *pubKeyID)
		if uint64(len(returnStruct.appliedTransfers)) == maxAppliedTransfers {
			break
		}
	}

	if len(returnStruct.appliedTransfers) > 0 {
		returnStruct.feeReceiverStateID, err = t.ApplyFee(*commitmentTokenIndex, *combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}

func (t *transactionExecutor) ApplyCreate2TransfersForSync(
	transfers []models.Create2Transfer,
	pubKeyIDs []uint32,
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

	commitmentTokenIndex, err := t.getTokenIndex(transfers[0].FromStateID)
	if err != nil {
		return nil, err
	}

	for i := range transfers {
		transfer := &transfers[i]

		_, err = t.handleApplyC2T(transfer, pubKeyIDs[i], returnStruct, combinedFee, commitmentTokenIndex)
		if err != nil {
			return nil, err
		}
	}

	if len(returnStruct.appliedTransfers) > 0 {
		_, err = t.ApplyFee(*commitmentTokenIndex, *combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}

func (t *transactionExecutor) getTokenIndex(stateID uint32) (*models.Uint256, error) {
	senderLeaf, err := t.storage.GetStateLeaf(stateID)
	if err != nil {
		return nil, err
	}
	return &senderLeaf.TokenIndex, nil
}

func (t *transactionExecutor) getOrRegisterPubKeyID(
	events chan *accountregistry.AccountRegistryPubkeyRegistered,
	transfer *models.Create2Transfer,
	tokenIndex models.Uint256,
) (*uint32, error) {
	pubKeyID, err := t.storage.GetUnusedPubKeyID(&transfer.ToPublicKey, &tokenIndex)
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	} else if st.IsNotFoundError(err) {
		return t.client.RegisterAccount(&transfer.ToPublicKey, events)
	}
	return pubKeyID, nil
}

func (t *transactionExecutor) handleApplyC2T(
	transfer *models.Create2Transfer,
	pubKeyID uint32,
	appliedTransfers *AppliedC2Transfers,
	combinedFee, tokenIndex *models.Uint256,
) (bool, error) {
	transferError, appError := t.ApplyCreate2Transfer(transfer, pubKeyID, *tokenIndex)
	if appError != nil {
		return false, appError
	}
	if transferError != nil {
		logAndSaveTransactionError(t.storage, &transfer.TransactionBase, transferError)
		appliedTransfers.invalidTransfers = append(appliedTransfers.invalidTransfers, *transfer)
		return false, nil
	}

	appliedTransfers.appliedTransfers = append(appliedTransfers.appliedTransfers, *transfer)
	*combinedFee = *combinedFee.Add(&transfer.Fee)
	return true, nil
}
