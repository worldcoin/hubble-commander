package stored

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
)

func TestTx_Bytes_Transfer(t *testing.T) {
	transfer := &models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 11,
			Amount:      models.MakeUint256(10),
			Fee:         models.MakeUint256(111),
			Nonce:       models.MakeUint256(1),
			Signature:   models.Signature{1, 2, 3, 4, 5},
			ReceiveTime: models.NewTimestamp(time.Unix(10, 0).UTC()),
			CommitmentID: &models.CommitmentID{
				BatchID:      models.MakeUint256(10),
				IndexInBatch: 2,
			},
			ErrorMessage: ref.String("some error message"),
		},
		ToStateID: 5,
	}

	storedTransaction := NewTxFromTransfer(transfer)
	bytes := storedTransaction.Bytes()

	decodedStoredTx := Tx{}
	err := decodedStoredTx.SetBytes(bytes)
	require.NoError(t, err)
	require.EqualValues(t, *storedTransaction, decodedStoredTx)

	storedTxReceipt := NewTxReceiptFromTransfer(transfer)
	decodedTransfer := decodedStoredTx.ToTransfer(storedTxReceipt)
	require.Equal(t, *transfer, *decodedTransfer)
}

func TestTx_Bytes_Create2Transfer(t *testing.T) {
	transfer := &models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Create2Transfer,
			FromStateID: 11,
			Amount:      models.MakeUint256(10),
			Fee:         models.MakeUint256(111),
			Nonce:       models.MakeUint256(1),
			Signature:   models.Signature{1, 2, 3, 4, 5},
			ReceiveTime: models.NewTimestamp(time.Unix(10, 0).UTC()),
			CommitmentID: &models.CommitmentID{
				BatchID:      models.MakeUint256(10),
				IndexInBatch: 2,
			},
			ErrorMessage: ref.String("some error message"),
		},
		ToStateID:   ref.Uint32(5),
		ToPublicKey: models.PublicKey{1, 2, 3, 4},
	}

	storedTransaction := NewTxFromCreate2Transfer(transfer)
	bytes := storedTransaction.Bytes()

	decodedStoredTx := Tx{}
	err := decodedStoredTx.SetBytes(bytes)
	require.NoError(t, err)
	require.EqualValues(t, *storedTransaction, decodedStoredTx)

	storedTxReceipt := NewTxReceiptFromCreate2Transfer(transfer)
	decodedTransfer := decodedStoredTx.ToCreate2Transfer(storedTxReceipt)
	require.Equal(t, *transfer, *decodedTransfer)
}

func TestTx_Bytes_MassMigration(t *testing.T) {
	massMigration := &models.MassMigration{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.MassMigration,
			FromStateID: 11,
			Amount:      models.MakeUint256(10),
			Fee:         models.MakeUint256(111),
			Nonce:       models.MakeUint256(1),
			Signature:   models.Signature{1, 2, 3, 4, 5},
			ReceiveTime: models.NewTimestamp(time.Unix(10, 0).UTC()),
			CommitmentID: &models.CommitmentID{
				BatchID:      models.MakeUint256(10),
				IndexInBatch: 2,
			},
			ErrorMessage: ref.String("some error message"),
		},
		SpokeID: 5,
	}

	storedTransaction := NewTxFromMassMigration(massMigration)
	bytes := storedTransaction.Bytes()

	decodedStoredTx := Tx{}
	err := decodedStoredTx.SetBytes(bytes)
	require.NoError(t, err)
	require.EqualValues(t, *storedTransaction, decodedStoredTx)

	storedTxReceipt := NewTxReceiptFromMassMigration(massMigration)
	decodedTransfer := decodedStoredTx.ToMassMigration(storedTxReceipt)
	require.Equal(t, *massMigration, *decodedTransfer)
}

func TestTx_ToTransfer_InvalidType(t *testing.T) {
	tx := NewTxFromCreate2Transfer(&models.Create2Transfer{})
	txReceipt := NewTxReceiptFromCreate2Transfer(&models.Create2Transfer{})

	require.Panics(t, func() {
		tx.ToTransfer(txReceipt)
	})
}

func TestTx_ToCreate2Transfer_InvalidType(t *testing.T) {
	tx := NewTxFromTransfer(&models.Transfer{})
	txReceipt := NewTxReceiptFromTransfer(&models.Transfer{})

	require.Panics(t, func() {
		tx.ToCreate2Transfer(txReceipt)
	})
}

func TestTx_ToMassMigration_InvalidType(t *testing.T) {
	tx := NewTxFromMassMigration(&models.MassMigration{})
	txReceipt := NewTxReceiptFromMassMigration(&models.MassMigration{})

	require.Panics(t, func() {
		tx.ToTransfer(txReceipt)
	})
}
