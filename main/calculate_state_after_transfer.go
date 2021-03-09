package main

import (
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
)

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
	// TODO: Do we check if sender and receiver states exist?

	if tx == nil {
		err = fmt.Errorf("transaction cannot be nil")
		return
	}

	if senderState == nil || receiverState == nil {
		err = fmt.Errorf("sender/receiver cannot be nil")
		return
	}

	if senderState.Nonce.Cmp(&tx.Nonce.Int) != 0 {
		err = fmt.Errorf("incorrect nonce")
		return
	}

	totalAmount := big.NewInt(0)
	totalAmount.Add(&tx.Amount.Int, &tx.Fee.Int)

	if senderState.Balance.Cmp(totalAmount) < 0 {
		err = fmt.Errorf("amount exceeds balance")
		return
	}

	newSenderState = *senderState
	newReceiverState = *receiverState

	newSenderState.Nonce.Add(&senderState.Nonce.Int, big.NewInt(1))

	newSenderState.Balance.Sub(&senderState.Balance.Int, totalAmount)

	newReceiverState.Balance.Add(&receiverState.Balance.Int, &tx.Amount.Int)

	return newSenderState, newReceiverState, nil
}
