package main

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
)

func CalculateStateAfterTransfer(senderState, receiverState *models.UserState, tx *models.Transaction) (*models.UserState, *models.UserState, error) {
	// TODO: VALIDATION

	senderState.Nonce.Add(&senderState.Nonce.Int, big.NewInt(1))

	totalAmount := big.NewInt(0)
	totalAmount.Add(&tx.Amount.Int, &tx.Fee.Int)
	senderState.Balance.Sub(&senderState.Balance.Int, totalAmount)

	receiverState.Balance.Add(&receiverState.Balance.Int, &tx.Amount.Int)

	return senderState, receiverState, nil
}
