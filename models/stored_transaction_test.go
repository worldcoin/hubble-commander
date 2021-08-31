package models

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
)

func TestStoredTransaction_Bytes_Transfer(t *testing.T) {
	transfer := &Transfer{
		TransactionBase: TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 11,
			Amount:      MakeUint256(10),
			Fee:         MakeUint256(111),
			Nonce:       MakeUint256(1),
			Signature:   Signature{1, 2, 3, 4, 5},
			ReceiveTime: NewTimestamp(time.Unix(10, 0).UTC()),
			CommitmentID: &CommitmentID{
				BatchID:      MakeUint256(10),
				IndexInBatch: 2,
			},
			ErrorMessage: ref.String("some error message"),
		},
		ToStateID: 5,
	}

	storedTransaction := MakeStoredTransactionFromTransfer(transfer)
	bytes := storedTransaction.Bytes()

	decodedStoredTransaction := StoredTransaction{}
	err := decodedStoredTransaction.SetBytes(bytes)
	require.NoError(t, err)
	require.EqualValues(t, storedTransaction, decodedStoredTransaction)

	decodedTransfer := decodedStoredTransaction.ToTransfer()
	require.Equal(t, *transfer, *decodedTransfer)
}

func TestStoredTransaction_ToTransfer_InvalidType(t *testing.T) {
	tx := MakeStoredTransactionFromCreate2Transfer(&Create2Transfer{})

	require.Panics(t, func() {
		tx.ToTransfer()
	})
}

func TestStoredTransaction_Bytes_Create2Transfer(t *testing.T) {
	transfer := &Create2Transfer{
		TransactionBase: TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Create2Transfer,
			FromStateID: 11,
			Amount:      MakeUint256(10),
			Fee:         MakeUint256(111),
			Nonce:       MakeUint256(1),
			Signature:   Signature{1, 2, 3, 4, 5},
			ReceiveTime: NewTimestamp(time.Unix(10, 0).UTC()),
			CommitmentID: &CommitmentID{
				BatchID:      MakeUint256(10),
				IndexInBatch: 2,
			},
			ErrorMessage: ref.String("some error message"),
		},
		ToStateID:   ref.Uint32(5),
		ToPublicKey: PublicKey{1, 2, 3, 4},
	}

	storedTransaction := MakeStoredTransactionFromCreate2Transfer(transfer)
	bytes := storedTransaction.Bytes()

	decodedStoredTransaction := StoredTransaction{}
	err := decodedStoredTransaction.SetBytes(bytes)
	require.NoError(t, err)
	require.EqualValues(t, storedTransaction, decodedStoredTransaction)

	decodedTransfer := decodedStoredTransaction.ToCreate2Transfer()
	require.Equal(t, *transfer, *decodedTransfer)
}

func TestStoredTransaction_ToCreate2Transfer_InvalidType(t *testing.T) {
	tx := MakeStoredTransactionFromTransfer(&Transfer{})

	require.Panics(t, func() {
		tx.ToCreate2Transfer()
	})
}
