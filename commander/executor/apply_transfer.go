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
	ErrBalanceTooLow      = NewDisputableTransferError(TransitionError, "not enough balance")
	ErrInvalidTokenID     = NewDisputableTransferError(TransitionError, "invalid sender or receiver token ID")
	ErrInvalidTokenAmount = NewDisputableTransferError(TransitionError, "amount cannot be equal to 0")
)

func (t *TransactionExecutor) ApplyTransfer(
	transfer models.GenericTransfer,
	commitmentTokenID models.Uint256,
) (transferError, appError error) {
	receiverStateID, err := getReceiverStateID(transfer)
	if err != nil {
		return nil, err
	}

	senderLeaf, err := t.storage.GetStateLeaf(transfer.GetFromStateID())
	if err != nil {
		return nil, err
	}
	receiverLeaf, err := t.storage.GetStateLeaf(*receiverStateID)
	if err != nil {
		return nil, err
	}

	senderState := senderLeaf.UserState
	receiverState := receiverLeaf.UserState

	if senderState.TokenID.Cmp(&commitmentTokenID) != 0 && receiverState.TokenID.Cmp(&commitmentTokenID) != 0 {
		return nil, ErrInvalidTokenID
	}

	if t.opts.AssumeNonces {
		transfer.SetNonce(senderState.Nonce)
	} else {
		nonce := transfer.GetNonce()
		err = validateTransferNonce(&senderState, &nonce)
		if err != nil {
			return err, nil
		}
	}

	newSenderState, newReceiverState, err := CalculateStateAfterTransfer(senderState, receiverState, transfer)
	if err != nil {
		return err, nil
	}

	err = t.stateTree.Set(transfer.GetFromStateID(), newSenderState)
	if err != nil {
		return nil, err
	}
	err = t.stateTree.Set(*receiverStateID, newReceiverState)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func CalculateStateAfterTransfer(
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

func getReceiverStateID(transfer models.GenericTransfer) (*uint32, error) {
	stateID := transfer.GetToStateID()
	if stateID == nil {
		return nil, ErrNilReceiverStateID
	}
	return stateID, nil
}

func validateTransferNonce(senderState *models.UserState, transferNonce *models.Uint256) error {
	comparison := transferNonce.Cmp(&senderState.Nonce)
	if comparison > 0 {
		return ErrNonceTooHigh
	} else if comparison < 0 {
		return ErrNonceTooLow
	}
	return nil
}
