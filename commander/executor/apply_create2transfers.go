package executor

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type AppliedC2Transfers struct {
	appliedTransfers []models.Create2Transfer
	invalidTransfers []models.Create2Transfer
	pendingAccounts  PendingAccounts
}

func NewAppliedC2Transfers(pendingAccounts PendingAccounts, capacity uint32) *AppliedC2Transfers {
	accounts := make(PendingAccounts, 0, uint32(len(pendingAccounts))+capacity)
	accounts = append(accounts, pendingAccounts...)

	return &AppliedC2Transfers{
		appliedTransfers: make([]models.Create2Transfer, 0, capacity),
		invalidTransfers: make([]models.Create2Transfer, 0),
		pendingAccounts:  accounts,
	}
}

func (t *TransactionExecutor) ApplyCreate2Transfers(
	pending *PendingC2Ts,
	maxApplied uint32,
	feeReceiver *FeeReceiver,
) (*AppliedC2Transfers, error) {
	if len(pending.Txs) == 0 {
		//TODO-reg: check if it's ok
		return &AppliedC2Transfers{}, nil
	}

	returnStruct := NewAppliedC2Transfers(pending.Accounts, t.cfg.MaxTxsPerCommitment)
	combinedFee := models.NewUint256(0)

	for i := range pending.Txs {
		if uint32(len(returnStruct.appliedTransfers)) == maxApplied {
			break
		}

		transfer := &pending.Txs[i]
		pubKeyID, err := t.getPubKeyID(returnStruct.pendingAccounts, transfer, feeReceiver.TokenID)
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

		returnStruct.pendingAccounts = append(returnStruct.pendingAccounts, models.AccountLeaf{
			PubKeyID:  *pubKeyID,
			PublicKey: transfer.ToPublicKey,
		})
		returnStruct.appliedTransfers = append(returnStruct.appliedTransfers, *appliedTransfer)
		*combinedFee = *combinedFee.Add(&appliedTransfer.Fee)
	}

	if len(returnStruct.appliedTransfers) > 0 {
		_, err := t.ApplyFee(feeReceiver.StateID, *combinedFee)
		if err != nil {
			return nil, err
		}
	}

	return returnStruct, nil
}

func (t *TransactionExecutor) ApplyCreate2TransfersForSync(
	transfers []models.Create2Transfer,
	pubKeyIDs []uint32,
	feeReceiverStateID uint32,
) ([]models.Create2Transfer, []models.StateMerkleProof, error) {
	transfersLen := len(transfers)
	if transfersLen != len(pubKeyIDs) {
		return nil, nil, ErrInvalidSlicesLength
	}

	appliedTransfers := make([]models.Create2Transfer, 0, transfersLen)
	stateChangeProofs := make([]models.StateMerkleProof, 0, 2*transfersLen+1)
	combinedFee := models.NewUint256(0)

	tokenID, err := t.getCommitmentTokenID(models.Create2TransferArray(transfers), feeReceiverStateID)
	if err != nil {
		return nil, nil, err
	}

	for i := range transfers {
		transfer := &transfers[i]

		synced, transferError, appError := t.ApplyCreate2TransferForSync(transfer, pubKeyIDs[i], *tokenID)
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

	stateProof, commitmentError, appError := t.ApplyFeeForSync(feeReceiverStateID, tokenID, combinedFee)
	if appError != nil {
		return nil, nil, appError
	}
	stateChangeProofs = append(stateChangeProofs, *stateProof)
	if commitmentError != nil {
		return nil, nil, NewDisputableErrorWithProofs(Transition, commitmentError.Error(), stateChangeProofs)
	}

	return appliedTransfers, stateChangeProofs, nil
}

func (t *TransactionExecutor) getPubKeyID(registrations PendingAccounts, transfer *models.Create2Transfer, tokenID models.Uint256) (
	*uint32, error,
) {
	pubKeyID, err := t.storage.GetUnusedPubKeyID(&transfer.ToPublicKey, &tokenID)
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	} else if st.IsNotFoundError(err) {
		return registrations.NextPubKeyID(t.client)
	}
	return pubKeyID, nil
}

type PendingC2Ts struct {
	Txs      []models.Create2Transfer
	Accounts PendingAccounts
}

type PendingAccounts []models.AccountLeaf

func (p PendingAccounts) ToPubKeyIDs() []uint32 {
	pubKeyIds := make([]uint32, 0, len(p))
	for i := range p {
		pubKeyIds = append(pubKeyIds, p[i].PubKeyID)
	}
	return pubKeyIds
}

func (p PendingAccounts) NextPubKeyID(client *eth.Client) (*uint32, error) {
	if len(p) == 0 {
		return client.GetNextSingleRegistrationPubKeyID()
	}
	nextPubKeyID := p[len(p)-1].PubKeyID + 1
	return &nextPubKeyID, nil
}
