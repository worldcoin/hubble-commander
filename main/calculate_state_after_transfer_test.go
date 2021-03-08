package main

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
)

func Test_CalculateStateAfterTransfer_UpdatesSenderAndReceiverStates(t *testing.T) {
	senderState := models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}
	receiverState := models.UserState{
		AccountIndex: models.MakeUint256(2),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(0),
		Nonce:        models.MakeUint256(0),
	}

	tx := models.Transaction{
		FromIndex: models.MakeUint256(1),
		ToIndex:   models.MakeUint256(2),
		Amount:    models.MakeUint256(100),
		Fee:       models.MakeUint256(10),
		Nonce:     models.MakeUint256(0),
	}

	newSenderState, newReceiverState, err := CalculateStateAfterTransfer(&senderState, &receiverState, &tx)
	require.NoError(t, err)

	require.Equal(t, newSenderState.Nonce, models.MakeUint256(1))
	require.Equal(t, newSenderState.Balance, models.MakeUint256(310))

	require.Equal(t, newReceiverState.Nonce, models.MakeUint256(0))
	require.Equal(t, newReceiverState.Balance, models.MakeUint256(100))
}
