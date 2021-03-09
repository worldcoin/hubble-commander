package main

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
)

var (
	senderState = models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}
	receiverState = models.UserState{
		AccountIndex: models.MakeUint256(2),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(0),
		Nonce:        models.MakeUint256(0),
	}
)

func Test_CalculateStateAfterTransfer_UpdatesSenderAndReceiverStates(t *testing.T) {
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

	require.NotEqual(t, &newSenderState, &senderState)
	require.NotEqual(t, &newReceiverState, &receiverState)
}

func Test_CalculateStateAfterTransfer_Validation_Nonce(t *testing.T) {
	tx := models.Transaction{
		FromIndex: models.MakeUint256(1),
		ToIndex:   models.MakeUint256(2),
		Amount:    models.MakeUint256(100),
		Fee:       models.MakeUint256(10),
		Nonce:     models.MakeUint256(1),
	}

	_, _, err := CalculateStateAfterTransfer(&senderState, &receiverState, &tx)
	require.Error(t, err)
}

func Test_CalculateStateAfterTransfer_Validation_Balance(t *testing.T) {
	tx := models.Transaction{
		FromIndex: models.MakeUint256(1),
		ToIndex:   models.MakeUint256(2),
		Amount:    models.MakeUint256(400),
		Fee:       models.MakeUint256(50),
		Nonce:     models.MakeUint256(0),
	}

	_, _, err := CalculateStateAfterTransfer(&senderState, &receiverState, &tx)
	require.Error(t, err)
}

func Test_CalculateStateAfterTransfer_Validation_Account(t *testing.T) {
	tx := models.Transaction{
		Nonce: models.MakeUint256(0),
	}

	_, _, err := CalculateStateAfterTransfer(nil, &receiverState, &tx)
	require.Error(t, err)

	_, _, err = CalculateStateAfterTransfer(&senderState, nil, &tx)
	require.Error(t, err)

	_, _, err = CalculateStateAfterTransfer(&senderState, &receiverState, nil)
	require.Error(t, err)
}
