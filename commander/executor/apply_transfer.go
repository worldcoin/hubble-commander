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
	transfer models.GenericTransfer,
	commitmentTokenID models.Uint256,
) (transferError, appError error) {
	senderState, receiverState, err := t.getParticipantsStates(transfer)
	if err != nil {
		return nil, err
	}

	if appErr := t.validateTokenIDs(senderState, receiverState, commitmentTokenID); appErr != nil {
		return nil, appErr
	}

	if tErr := validateTransferNonce(&senderState.UserState, transfer.GetNonce()); tErr != nil {
		return tErr, nil
	}

	newSenderState, newReceiverState, tErr := calculateStateAfterTransfer(senderState.UserState, receiverState.UserState, transfer)
	if tErr != nil {
		return tErr, nil
	}

	err = t.stateTree.Set(senderState.StateID, newSenderState)
	if err != nil {
		return nil, err
	}
	err = t.stateTree.Set(receiverState.StateID, newReceiverState)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

type SyncedTransfer struct {
	transfer             models.GenericTransfer
	senderStateWitness   models.Witness
	receiverStateWitness models.Witness
}

func (t *TransactionExecutor) ApplyTransferForSync(transfer models.GenericTransfer, commitmentTokenID models.Uint256) (
	syncedTransfer *SyncedTransfer,
	transferError, appError error,
) {
	senderState, receiverState, err := t.getParticipantsStates(transfer)
	if err != nil {
		return nil, nil, err
	}

	if tErr := t.validateTokenIDs(senderState, receiverState, commitmentTokenID); tErr != nil {
		return nil, tErr, nil
	}

	newSenderState, newReceiverState, tErr := calculateStateAfterTransfer(senderState.UserState, receiverState.UserState, transfer)
	if tErr != nil {
		return nil, tErr, nil
	}

	senderWitness, err := t.stateTree.SetReturningWitness(senderState.StateID, newSenderState)
	if err != nil {
		return nil, nil, err
	}
	receiverWitness, err := t.stateTree.SetReturningWitness(receiverState.StateID, newReceiverState)
	if err != nil {
		return nil, nil, err
	}

	syncedTransfer = &SyncedTransfer{
		transfer:             transfer.Copy(),
		senderStateWitness:   senderWitness,
		receiverStateWitness: receiverWitness,
	}
	syncedTransfer.transfer.SetNonce(senderState.Nonce)

	return syncedTransfer, nil, nil
}

func (t *TransactionExecutor) getParticipantsStates(transfer models.GenericTransfer) (
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

func (t *TransactionExecutor) validateTokenIDs(senderState, receiverState *models.StateLeaf, commitmentTokenID models.Uint256) error {
	if senderState.TokenID.Cmp(&commitmentTokenID) != 0 || receiverState.TokenID.Cmp(&commitmentTokenID) != 0 {
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
	transfer models.GenericTransfer,
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
