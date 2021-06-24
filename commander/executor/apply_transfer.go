package executor

import (
	"errors"
	"reflect"

	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrNonceTooLow           = errors.New("nonce too low")
	ErrNonceTooHigh          = errors.New("nonce too high")
	ErrInvalidSliceLength    = errors.New("invalid slices length")
	ErrNilReceiverStateID    = errors.New("transfer receiver state id cannot be nil")
	ErrBalanceTooLow         = NewDisputableTransferError("not enough balance")
	ErrIncorrectTokenIndices = NewDisputableTransferError("sender's, receiver's and fee receiver's token indices are not the same")
	ErrInvalidTokenAmount    = NewDisputableTransferError("amount cannot be equal to 0")
)

func (t *TransactionExecutor) ApplyTransfer(
	transfer models.GenericTransfer,
	commitmentTokenIndex models.Uint256,
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

	if senderState.TokenIndex.Cmp(&commitmentTokenIndex) != 0 && receiverState.TokenIndex.Cmp(&commitmentTokenIndex) != 0 {
		return nil, ErrIncorrectTokenIndices
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

	newSenderState, newReceiverState, err := CalculateStateAfterTransfer(
		&senderState,
		&receiverState,
		transfer,
	)
	if err != nil {
		return err, nil
	}

	if !reflect.DeepEqual(newSenderState, senderState) {
		err = t.stateTree.Set(transfer.GetFromStateID(), &newSenderState)
		if err != nil {
			return nil, err
		}
	}
	if !reflect.DeepEqual(newReceiverState, receiverState) {
		err = t.stateTree.Set(*receiverStateID, &newReceiverState)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func CalculateStateAfterTransfer(
	senderState,
	receiverState *models.UserState,
	transfer models.GenericTransfer,
) (
	newSenderState models.UserState,
	newReceiverState models.UserState,
	err error,
) {
	amount := transfer.GetAmount()
	fee := transfer.GetFee()

	if amount.CmpN(0) <= 0 {
		err = ErrInvalidTokenAmount
		return
	}

	totalAmount := amount.Add(&fee)
	if senderState.Balance.Cmp(totalAmount) < 0 {
		err = ErrBalanceTooLow
		return
	}

	newSenderState = *senderState
	newReceiverState = *receiverState

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
