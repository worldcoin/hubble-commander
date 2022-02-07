package stored

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

const (
	txTransferBodyLength        = 4
	txCreate2TransferBodyLength = models.PublicKeyLength + 5
	txMassMigrationBodyLength   = 4
)

func txBody(data []byte, transactionType txtype.TransactionType) (TxBody, error) {
	switch transactionType {
	case txtype.Transfer:
		body := new(TxTransferBody)
		err := body.SetBytes(data)
		return body, err
	case txtype.Create2Transfer:
		body := new(TxCreate2TransferBody)
		err := body.SetBytes(data)
		return body, err
	case txtype.MassMigration:
		body := new(TxMassMigrationBody)
		err := body.SetBytes(data)
		return body, err
	}
	return nil, nil
}

type TxBody interface {
	ByteEncoder
	BytesLen() int
	ToGenericTransaction(base *models.TransactionBase) models.GenericTransaction
}

func NewTxBody(tx models.GenericTransaction) TxBody {
	switch tx.Type() {
	case txtype.Transfer:
		return &TxTransferBody{
			ToStateID: tx.ToTransfer().ToStateID,
		}
	case txtype.Create2Transfer:
		return &TxCreate2TransferBody{
			ToPublicKey: tx.ToCreate2Transfer().ToPublicKey,
			ToStateID:   tx.ToCreate2Transfer().ToStateID,
		}
	case txtype.MassMigration:
		return &TxMassMigrationBody{
			SpokeID: tx.ToMassMigration().SpokeID,
		}
	}
	panic("unknown transaction type")
}

type TxTransferBody struct {
	ToStateID uint32
}

func (t *TxTransferBody) Bytes() []byte {
	b := make([]byte, txTransferBodyLength)
	binary.BigEndian.PutUint32(b, t.ToStateID)
	return b
}

func (t *TxTransferBody) SetBytes(data []byte) error {
	t.ToStateID = binary.BigEndian.Uint32(data)
	return nil
}

func (t *TxTransferBody) BytesLen() int {
	return txTransferBodyLength
}

func (t *TxTransferBody) ToGenericTransaction(base *models.TransactionBase) models.GenericTransaction {
	return &models.Transfer{
		TransactionBase: *base,
		ToStateID:       t.ToStateID,
	}
}

type TxCreate2TransferBody struct {
	ToPublicKey models.PublicKey
	ToStateID   *uint32
}

func (t *TxCreate2TransferBody) Bytes() []byte {
	b := make([]byte, txCreate2TransferBodyLength)
	copy(b[0:5], EncodeUint32Pointer(t.ToStateID))
	copy(b[5:], t.ToPublicKey.Bytes())
	return b
}

func (t *TxCreate2TransferBody) SetBytes(data []byte) error {
	if len(data) < txCreate2TransferBodyLength {
		return models.ErrInvalidLength
	}
	t.ToStateID = decodeUint32Pointer(data[0:5])
	return t.ToPublicKey.SetBytes(data[5:])
}

func (t *TxCreate2TransferBody) BytesLen() int {
	return txCreate2TransferBodyLength
}

func (t *TxCreate2TransferBody) ToGenericTransaction(base *models.TransactionBase) models.GenericTransaction {
	return &models.Create2Transfer{
		TransactionBase: *base,
		ToStateID:       t.ToStateID,
		ToPublicKey:     t.ToPublicKey,
	}
}

type TxMassMigrationBody struct {
	SpokeID uint32
}

func (t *TxMassMigrationBody) Bytes() []byte {
	b := make([]byte, txMassMigrationBodyLength)
	binary.BigEndian.PutUint32(b, t.SpokeID)
	return b
}

func (t *TxMassMigrationBody) SetBytes(data []byte) error {
	t.SpokeID = binary.BigEndian.Uint32(data)
	return nil
}

func (t *TxMassMigrationBody) BytesLen() int {
	return txMassMigrationBodyLength
}

func (t *TxMassMigrationBody) ToGenericTransaction(base *models.TransactionBase) models.GenericTransaction {
	return &models.MassMigration{
		TransactionBase: *base,
		SpokeID:         t.SpokeID,
	}
}
