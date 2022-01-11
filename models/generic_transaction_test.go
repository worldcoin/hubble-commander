package models

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/stretchr/testify/require"
)

func TestOldestTransaction_EmptyArray(t *testing.T) {
	array := MakeTransferArray()
	oldest := FindOldestTransactionTime(array)
	require.Nil(t, oldest)
}

func TestOldestTransaction_NoTxnHasTime(t *testing.T) {
	tx := Transfer{
		TransactionBase: TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      MakeUint256(400),
			Fee:         MakeUint256(100),
			Nonce:       MakeUint256(0),
		},
		ToStateID: 1,
	}
	array := MakeTransferArray(tx)

	oldest := FindOldestTransactionTime(array)
	require.Nil(t, oldest)
}

func TestOldestTransaction_FindsOldestTime(t *testing.T) {
	oneSecondAgo := time.Now().Add(-time.Second)
	twoSecondAgo := time.Now().Add(-2 * time.Second)

	txs := []Transfer{
		{
			TransactionBase: TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      MakeUint256(400),
				Fee:         MakeUint256(100),
				Nonce:       MakeUint256(0),
				ReceiveTime: NewTimestamp(oneSecondAgo),
			},
			ToStateID: 1,
		},
		{
			TransactionBase: TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      MakeUint256(400),
				Fee:         MakeUint256(100),
				Nonce:       MakeUint256(0),
				ReceiveTime: NewTimestamp(twoSecondAgo),
			},
			ToStateID: 1,
		},
	}
	array := MakeTransferArray(txs...)

	oldest := FindOldestTransactionTime(array)
	require.Equal(t, twoSecondAgo, oldest.Time)
}
