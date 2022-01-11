package executor

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/stretchr/testify/require"
)

func TestFindOldestTransactionTime_EmptyArray(t *testing.T) {
	array := models.MakeTransferArray()
	oldest := findOldestTransactionTime(array)
	require.Nil(t, oldest)
}

func TestFindOldestTransactionTime_NoTxHasTime(t *testing.T) {
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 1,
	}
	array := models.MakeTransferArray(tx)

	oldest := findOldestTransactionTime(array)
	require.Nil(t, oldest)
}

func TestFindOldestTransactionTime_FindsOldestTime(t *testing.T) {
	oneSecondAgo := time.Now().Add(-time.Second)
	twoSecondAgo := time.Now().Add(-2 * time.Second)

	txs := []models.Transfer{
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(400),
				Fee:         models.MakeUint256(100),
				Nonce:       models.MakeUint256(0),
				ReceiveTime: models.NewTimestamp(oneSecondAgo),
			},
			ToStateID: 1,
		},
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(400),
				Fee:         models.MakeUint256(100),
				Nonce:       models.MakeUint256(0),
				ReceiveTime: models.NewTimestamp(twoSecondAgo),
			},
			ToStateID: 1,
		},
	}
	array := models.MakeTransferArray(txs...)

	oldest := findOldestTransactionTime(array)
	require.Equal(t, twoSecondAgo, oldest.Time)
}
