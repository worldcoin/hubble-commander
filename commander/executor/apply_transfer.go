package executor

import (
	"errors"

	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrNonceTooLow        = errors.New("nonce too low")
	ErrNonceTooHigh       = errors.New("nonce too high")
	ErrInvalidSliceLength = errors.New("invalid slices length")
	ErrNilReceiverStateID = errors.New("transfer receiver state id cannot be nil")

	ErrBalanceTooLow      = errors.New("not enough balance")
	ErrInvalidTokenID     = errors.New("invalid sender or receiver token ID")
	ErrInvalidTokenAmount = errors.New("amount cannot be equal to 0")
)

func (t *TransactionExecutor) ApplyTransfer(
	transfer models.GenericTransaction,
	commitmentTokenID models.Uint256,
) (transferError, appError error) {
	senderState, receiverState, appError := t.getParticipantsStates(transfer)
	if appError != nil {
		return nil, appError
	}

	appError = t.validateSenderTokenID(senderState, commitmentTokenID)
	if appError != nil {
		return nil, appError
	}

	appError = t.validateReceiverTokenID(receiverState, commitmentTokenID)
	if appError != nil {
		return nil, appError
	}

	if tErr := validateTransferNonce(&senderState.UserState, transfer.GetNonce()); tErr != nil {
		return tErr, nil
	}

	newSenderState, newReceiverState, tErr := calculateStateAfterTransfer(senderState.UserState, receiverState.UserState, transfer)
	if tErr != nil {
		return tErr, nil
	}

	_, appError = t.stateTree.Set(senderState.StateID, newSenderState)
	if appError != nil {
		return nil, appError
	}
	_, appError = t.stateTree.Set(receiverState.StateID, newReceiverState)
	if appError != nil {
		return nil, appError
	}

	return nil, nil
}

type SyncedTransfer struct {
	transfer           models.GenericTransaction
	senderStateProof   models.StateMerkleProof
	receiverStateProof models.StateMerkleProof
}

func (t *TransactionExecutor) ApplyTransferForSync(transfer models.GenericTransaction, commitmentTokenID models.Uint256) (
	syncedTransfer *SyncedTransfer,
	transferError, appError error,
) {
	senderState, receiverState, appError := t.getParticipantsStates(transfer)
	if appError != nil {
		return nil, nil, appError
	}

	// TODO-AFS we need to split this into two functions: validateSenderTokenID and validateReceiverTokenID
	//  the second one will need to be called after the SetReturningProof of senderState.
	if tErr := t.validateSenderTokenID(senderState, commitmentTokenID); tErr != nil {
		return nil, tErr, nil
	}

	newSenderState, newReceiverState, tErr := calculateStateAfterTransfer(senderState.UserState, receiverState.UserState, transfer)
	if tErr != nil {
		return nil, tErr, nil
	}

	senderWitness, appError := t.stateTree.Set(senderState.StateID, newSenderState)
	if appError != nil {
		return nil, nil, appError
	}

	// TODO-AFS validateReceiverTokenID here, on error we need to return senderStateProof
	if tErr := t.validateReceiverTokenID(receiverState, commitmentTokenID); tErr != nil {
		return nil, tErr, nil
	}

	receiverWitness, appError := t.stateTree.Set(receiverState.StateID, newReceiverState)
	if appError != nil {
		return nil, nil, appError
	}

	syncedTransfer = &SyncedTransfer{
		transfer: transfer.Copy(),
		senderStateProof: models.StateMerkleProof{
			UserState: &senderState.UserState,
			Witness:   senderWitness,
		},
		receiverStateProof: models.StateMerkleProof{
			UserState: &receiverState.UserState,
			Witness:   receiverWitness,
		},
	}
	syncedTransfer.transfer.SetNonce(senderState.Nonce)

	return syncedTransfer, nil, nil
}

func (t *TransactionExecutor) getParticipantsStates(transfer models.GenericTransaction) (
	senderState, receiverState *models.StateLeaf,
	err error,
) {
	receiverStateID := transfer.GetToStateID()
	if receiverStateID == nil {
		return nil, nil, ErrNilReceiverStateID
	}

	senderLeaf, err := t.storage.GetStateLeaf(transfer.GetFromStateID())
	if err != nil {
		return nil, nil, err
	}
	receiverLeaf, err := t.storage.GetStateLeaf(*receiverStateID)
	if err != nil {
		return nil, nil, err
	}

	return senderLeaf, receiverLeaf, nil
}

func (t *TransactionExecutor) validateSenderTokenID(senderState *models.StateLeaf, commitmentTokenID models.Uint256) error {
	if senderState.TokenID.Cmp(&commitmentTokenID) != 0 {
		return ErrInvalidTokenID
	}
	return nil
}

func (t *TransactionExecutor) validateReceiverTokenID(receiverState *models.StateLeaf, commitmentTokenID models.Uint256) error {
	if receiverState.TokenID.Cmp(&commitmentTokenID) != 0 {
		return ErrInvalidTokenID
	}
	return nil
}

func validateTransferNonce(senderState *models.UserState, transferNonce models.Uint256) error {
	comparison := transferNonce.Cmp(&senderState.Nonce)
	if comparison > 0 {
		return ErrNonceTooHigh
	} else if comparison < 0 {
		return ErrNonceTooLow
	}
	return nil
}

func calculateStateAfterTransfer(
	senderState, receiverState models.UserState, // nolint:gocritic
	transfer models.GenericTransaction,
) (
	newSenderState, newReceiverState *models.UserState,
	err error,
) {
	amount := transfer.GetAmount()
	fee := transfer.GetFee()

	if amount.CmpN(0) <= 0 {
		return nil, nil, ErrInvalidTokenAmount
	}

	totalAmount := amount.Add(&fee)
	if senderState.Balance.Cmp(totalAmount) < 0 {
		return nil, nil, ErrBalanceTooLow
	}

	newSenderState = &senderState
	newReceiverState = &receiverState

	newSenderState.Nonce = *newSenderState.Nonce.AddN(1)
	newSenderState.Balance = *newSenderState.Balance.Sub(totalAmount)
	newReceiverState.Balance = *newReceiverState.Balance.Add(&amount)

	return newSenderState, newReceiverState, nil
}
