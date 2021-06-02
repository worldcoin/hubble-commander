package commander

import (
	"github.com/Worldcoin/hubble-commander/models"
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
) (*AppliedC2Transfers, error) {
	if len(transfers) == 0 {
		return nil, nil
	}
	events, unsubscribe, err := t.client.WatchRegistrations(&bind.WatchOpts{})
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

		pubKeyID, err = getOrRegisterPubKeyID(t.storage, t.client, events, transfer, *commitmentTokenIndex)
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
		if uint32(len(returnStruct.appliedTransfers)) == t.cfg.TxsPerCommitment {
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
		return nil, nil // ! FIX ME
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

		if uint32(len(returnStruct.appliedTransfers)) == t.cfg.TxsPerCommitment {
			break
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

func (t *transactionExecutor) handleApplyC2T(
	transfer *models.Create2Transfer,
	pubKeyID uint32,
	transactions *AppliedC2Transfers,
	combinedFee, tokenIndex *models.Uint256,
) (bool, error) {
	transferError, appError := ApplyCreate2Transfer(t.storage, transfer, pubKeyID, *tokenIndex)
	if appError != nil {
		return false, appError
	}
	if transferError != nil {
		logAndSaveTransactionError(t.storage, &transfer.TransactionBase, transferError)
		transactions.invalidTransfers = append(transactions.invalidTransfers, *transfer)
		return false, nil
	}

	transactions.appliedTransfers = append(transactions.appliedTransfers, *transfer)
	*combinedFee = *combinedFee.Add(&transfer.Fee)
	return true, nil
}
