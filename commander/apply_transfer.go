package commander

import (
	"errors"
	"reflect"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

var (
	ErrIncorrectTokenIndices = errors.New("sender's, receiver's and fee receiver's token indices are not the same")
	ErrNonceTooLow           = errors.New("nonce too low")
	ErrNonceTooHigh          = errors.New("nonce too high")
	ErrBalanceTooLow         = errors.New("not enough balance")
)

func ApplyTransfer(
	stateTree *storage.StateTree,
	transfer *models.Transfer,
	commitmentTokenIndex models.Uint256,
) (transferError, appError error) {
	senderLeaf, err := stateTree.Leaf(transfer.FromStateID)
	if err != nil {
		return nil, err
	}
	receiverLeaf, err := stateTree.Leaf(transfer.ToStateID)
	if err != nil {
		return nil, err
	}

	senderState := senderLeaf.UserState
	receiverState := receiverLeaf.UserState

	if senderState.TokenIndex.Cmp(&commitmentTokenIndex) != 0 && receiverState.TokenIndex.Cmp(&commitmentTokenIndex) != 0 {
		return nil, ErrIncorrectTokenIndices
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
		err = stateTree.Set(transfer.FromStateID, &newSenderState)
		if err != nil {
			return nil, err
		}
	}

	if !reflect.DeepEqual(newReceiverState, receiverState) {
		err = stateTree.Set(transfer.ToStateID, &newReceiverState)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func CalculateStateAfterTransfer(
	senderState,
	receiverState *models.UserState,
	transfer *models.Transfer,
) (
	newSenderState models.UserState,
	newReceiverState models.UserState,
	err error,
) {
	comparison := transfer.Nonce.Cmp(&senderState.Nonce)
	if comparison > 0 {
		err = ErrNonceTooHigh
		return
	} else if comparison < 0 {
		err = ErrNonceTooLow
		return
	}

	totalAmount := transfer.Amount.Add(&transfer.Fee)
	if senderState.Balance.Cmp(totalAmount) < 0 {
		err = ErrBalanceTooLow
		return
	}

	newSenderState = *senderState
	newReceiverState = *receiverState

	newSenderState.Nonce = *newSenderState.Nonce.AddN(1)
	newSenderState.Balance = *newSenderState.Balance.Sub(totalAmount)
	newReceiverState.Balance = *newReceiverState.Balance.Add(&transfer.Amount)

	return newSenderState, newReceiverState, nil
}
