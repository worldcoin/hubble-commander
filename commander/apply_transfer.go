package commander

import (
	"errors"
	"math/big"
	"reflect"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

var (
	ErrStateTreeIsNil   = errors.New("state tree cannot be nil")
	ErrTransactionIsNil = errors.New("transaction cannot be nil")
	ErrUserStateIsNil   = errors.New("sender/receiver state cannot be nil")
	ErrNonceTooLow      = errors.New("nonce too low")
	ErrNonceTooHigh     = errors.New("nonce too high")
	ErrBalanceTooLow    = errors.New("amount exceeds balance")
)

func ApplyTransfer(stateTree *storage.StateTree, tx *models.Transaction) (txError, appError error) {
	if stateTree == nil {
		return nil, ErrStateTreeIsNil
	}
	if tx == nil {
		return nil, ErrTransactionIsNil
	}

	senderIndex := uint32(tx.FromIndex.Uint64())
	senderLeaf, err := stateTree.Leaf(senderIndex)
	if err != nil {
		return nil, err
	}
	receiverIndex := uint32(tx.ToIndex.Uint64())
	receiverLeaf, err := stateTree.Leaf(receiverIndex)
	if err != nil {
		return nil, err
	}

	if senderLeaf == nil || receiverLeaf == nil {
		return ErrUserStateIsNil, nil
	}

	senderState := senderLeaf.UserState
	receiverState := receiverLeaf.UserState

	newSenderState, newReceiverState, err := CalculateStateAfterTransfer(
		&senderState,
		&receiverState,
		tx,
	)
	if err != nil {
		return err, nil
	}
	if reflect.DeepEqual(newSenderState, senderState) && reflect.DeepEqual(newReceiverState, receiverState) {
		return nil, nil
	}

	err = stateTree.Set(senderIndex, &newSenderState)
	if err != nil {
		return nil, err
	}
	err = stateTree.Set(receiverIndex, &newReceiverState)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func CalculateStateAfterTransfer(
	senderState,
	receiverState *models.UserState,
	tx *models.Transaction,
) (
	newSenderState models.UserState,
	newReceiverState models.UserState,
	err error,
) {
	// TODO: Signature validation

	comparison := tx.Nonce.Cmp(&senderState.Nonce.Int)
	if comparison > 0 {
		err = ErrNonceTooHigh
		return
	} else if comparison < 0 {
		err = ErrNonceTooLow
		return
	}

	totalAmount := big.NewInt(0)
	totalAmount.Add(&tx.Amount.Int, &tx.Fee.Int)

	if senderState.Balance.Cmp(totalAmount) < 0 {
		err = ErrBalanceTooLow
		return
	}

	newSenderState = *senderState
	newReceiverState = *receiverState

	newSenderState.Nonce.Add(&senderState.Nonce.Int, big.NewInt(1))

	newSenderState.Balance.Sub(&senderState.Balance.Int, totalAmount)

	newReceiverState.Balance.Add(&receiverState.Balance.Int, &tx.Amount.Int)

	return newSenderState, newReceiverState, nil
}

func ApplyFee(stateTree *storage.StateTree, feeReceiverIndex uint32, fee models.Uint256) error {
	feeReceiverLeaf, err := stateTree.Leaf(feeReceiverIndex)
	if err != nil {
		return err
	}

	feeReceiverState := feeReceiverLeaf.UserState
	feeReceiverState.Balance.Add(&feeReceiverState.Balance.Int, &fee.Int)

	err = stateTree.Set(feeReceiverIndex, &feeReceiverState)
	if err != nil {
		return err
	}

	return nil
}
