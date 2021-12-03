package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

const (
	storedTxBytesLength               = 213
	storedTxTransferBodyLength        = 4
	storedTxCreate2TransferBodyLength = PublicKeyLength
	storedTxMassMigrationBodyLength   = 32
)

var (
	StoredTxName                = getTypeName(StoredTx{})
	StoredTxPrefix              = getBadgerHoldPrefix(StoredTx{})
	errInvalidStoredTxIndexType = errors.New("invalid StoredTx index type")
)

type StoredTx struct {
	Hash        common.Hash
	TxType      txtype.TransactionType
	FromStateID uint32
	Amount      Uint256
	Fee         Uint256
	Nonce       Uint256
	Signature   Signature
	ReceiveTime *Timestamp

	Body TxBody
}

func NewStoredTxFromTransfer(t *Transfer) *StoredTx {
	return &StoredTx{
		Hash:        t.Hash,
		TxType:      t.TxType,
		FromStateID: t.FromStateID,
		Amount:      t.Amount,
		Fee:         t.Fee,
		Nonce:       t.Nonce,
		Signature:   t.Signature,
		ReceiveTime: t.ReceiveTime,
		Body: &StoredTxTransferBody{
			ToStateID: t.ToStateID,
		},
	}
}

func NewStoredTxFromCreate2Transfer(t *Create2Transfer) *StoredTx {
	return &StoredTx{
		Hash:        t.Hash,
		TxType:      t.TxType,
		FromStateID: t.FromStateID,
		Amount:      t.Amount,
		Fee:         t.Fee,
		Nonce:       t.Nonce,
		Signature:   t.Signature,
		ReceiveTime: t.ReceiveTime,
		Body: &StoredTxCreate2TransferBody{
			ToPublicKey: t.ToPublicKey,
		},
	}
}

func NewStoredTxFromMassMigration(m *MassMigration) *StoredTx {
	return &StoredTx{
		Hash:        m.Hash,
		TxType:      m.TxType,
		FromStateID: m.FromStateID,
		Amount:      m.Amount,
		Fee:         m.Fee,
		Nonce:       m.Nonce,
		Signature:   m.Signature,
		ReceiveTime: m.ReceiveTime,
		Body: &StoredTxMassMigrationBody{
			SpokeID: m.SpokeID,
		},
	}
}

func (t *StoredTx) Bytes() []byte {
	b := make([]byte, t.BytesLen())
	copy(b[0:32], t.Hash.Bytes())
	b[32] = byte(t.TxType)
	binary.BigEndian.PutUint32(b[33:37], t.FromStateID)
	copy(b[37:69], t.Amount.Bytes())
	copy(b[69:101], t.Fee.Bytes())
	copy(b[101:133], t.Nonce.Bytes())
	copy(b[133:197], t.Signature.Bytes())
	copy(b[197:213], encodeTimestampPointer(t.ReceiveTime))
	copy(b[213:], t.Body.Bytes())

	return b
}

func (t *StoredTx) SetBytes(data []byte) error {
	if len(data) < storedTxBytesLength {
		return ErrInvalidLength
	}
	err := t.Signature.SetBytes(data[133:197])
	if err != nil {
		return err
	}
	receiveTime, err := decodeTimestampPointer(data[197:213])
	if err != nil {
		return err
	}

	txType := txtype.TransactionType(data[32])
	body, err := txBody(data[213:], txType)
	if err != nil {
		return err
	}

	t.Hash.SetBytes(data[0:32])
	t.TxType = txType
	t.FromStateID = binary.BigEndian.Uint32(data[33:37])
	t.Amount.SetBytes(data[37:69])
	t.Fee.SetBytes(data[69:101])
	t.Nonce.SetBytes(data[101:133])
	t.ReceiveTime = receiveTime
	t.Body = body
	return nil
}

func (t *StoredTx) BytesLen() int {
	return storedTxBytesLength + t.Body.BytesLen()
}

func (t *StoredTx) ToTransfer(txReceipt *StoredTxReceipt) *Transfer {
	transferBody, ok := t.Body.(*StoredTxTransferBody)
	if !ok {
		panic("invalid transfer body type")
	}

	transfer := &Transfer{
		TransactionBase: TransactionBase{
			Hash:        t.Hash,
			TxType:      t.TxType,
			FromStateID: t.FromStateID,
			Amount:      t.Amount,
			Fee:         t.Fee,
			Nonce:       t.Nonce,
			Signature:   t.Signature,
			ReceiveTime: t.ReceiveTime,
		},
		ToStateID: transferBody.ToStateID,
	}

	if txReceipt != nil {
		transfer.CommitmentID = txReceipt.CommitmentID
		transfer.ErrorMessage = txReceipt.ErrorMessage
	}
	return transfer
}

func (t *StoredTx) ToCreate2Transfer(txReceipt *StoredTxReceipt) *Create2Transfer {
	c2tBody, ok := t.Body.(*StoredTxCreate2TransferBody)
	if !ok {
		panic("invalid create2Transfer body type")
	}

	transfer := &Create2Transfer{
		TransactionBase: TransactionBase{
			Hash:        t.Hash,
			TxType:      t.TxType,
			FromStateID: t.FromStateID,
			Amount:      t.Amount,
			Fee:         t.Fee,
			Nonce:       t.Nonce,
			Signature:   t.Signature,
			ReceiveTime: t.ReceiveTime,
		},
		ToPublicKey: c2tBody.ToPublicKey,
	}

	if txReceipt != nil {
		transfer.CommitmentID = txReceipt.CommitmentID
		transfer.ErrorMessage = txReceipt.ErrorMessage
		transfer.ToStateID = txReceipt.ToStateID
	}
	return transfer
}

func (t *StoredTx) ToMassMigration(txReceipt *StoredTxReceipt) *MassMigration {
	massMigrationBody, ok := t.Body.(*StoredTxMassMigrationBody)
	if !ok {
		panic("invalid mass migration body type")
	}

	massMigration := &MassMigration{
		TransactionBase: TransactionBase{
			Hash:        t.Hash,
			TxType:      t.TxType,
			FromStateID: t.FromStateID,
			Amount:      t.Amount,
			Fee:         t.Fee,
			Nonce:       t.Nonce,
			Signature:   t.Signature,
			ReceiveTime: t.ReceiveTime,
		},
		SpokeID: massMigrationBody.SpokeID,
	}

	if txReceipt != nil {
		massMigration.CommitmentID = txReceipt.CommitmentID
		massMigration.ErrorMessage = txReceipt.ErrorMessage
	}
	return massMigration
}

func txBody(data []byte, transactionType txtype.TransactionType) (TxBody, error) {
	switch transactionType {
	case txtype.Transfer:
		body := new(StoredTxTransferBody)
		err := body.SetBytes(data)
		return body, err
	case txtype.Create2Transfer:
		body := new(StoredTxCreate2TransferBody)
		err := body.SetBytes(data)
		return body, err
	case txtype.MassMigration:
		body := new(StoredTxMassMigrationBody)
		err := body.SetBytes(data)
		return body, err
	}
	return nil, nil
}

// nolint:gocritic
// Type implements badgerhold.Storer
func (t StoredTx) Type() string {
	return string(StoredTxName)
}

// nolint:gocritic
// Indexes implements badgerhold.Storer
func (t StoredTx) Indexes() map[string]bh.Index {
	return map[string]bh.Index{
		"FromStateID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredTx(value)
				if err != nil {
					return nil, err
				}
				return EncodeUint32(v.FromStateID), nil
			},
		},
		"ToStateID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredTx(value)
				if err != nil {
					return nil, err
				}

				transferBody, ok := v.Body.(*StoredTxTransferBody)
				if !ok {
					return nil, nil
				}
				return EncodeUint32(transferBody.ToStateID), nil
			},
		},
	}
}

func interfaceToStoredTx(value interface{}) (*StoredTx, error) {
	p, ok := value.(*StoredTx)
	if ok {
		return p, nil
	}
	v, ok := value.(StoredTx)
	if ok {
		return &v, nil
	}
	return nil, errInvalidStoredTxIndexType
}

type TxBody interface {
	ByteEncoder
	BytesLen() int
}

type StoredTxTransferBody struct {
	ToStateID uint32
}

func (t *StoredTxTransferBody) Bytes() []byte {
	b := make([]byte, storedTxTransferBodyLength)
	binary.BigEndian.PutUint32(b, t.ToStateID)
	return b
}

func (t *StoredTxTransferBody) SetBytes(data []byte) error {
	t.ToStateID = binary.BigEndian.Uint32(data)
	return nil
}

func (t *StoredTxTransferBody) BytesLen() int {
	return storedTxTransferBodyLength
}

type StoredTxCreate2TransferBody struct {
	ToPublicKey PublicKey
}

func (t *StoredTxCreate2TransferBody) Bytes() []byte {
	return t.ToPublicKey.Bytes()
}

func (t *StoredTxCreate2TransferBody) SetBytes(data []byte) error {
	return t.ToPublicKey.SetBytes(data)
}

func (t *StoredTxCreate2TransferBody) BytesLen() int {
	return storedTxCreate2TransferBodyLength
}

type StoredTxMassMigrationBody struct {
	SpokeID uint32
}

func (t *StoredTxMassMigrationBody) Bytes() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b[0:], t.SpokeID)
	return b
}

func (t *StoredTxMassMigrationBody) SetBytes(data []byte) error {
	t.SpokeID = binary.BigEndian.Uint32(data)
	return nil
}

func (t *StoredTxMassMigrationBody) BytesLen() int {
	return storedTxMassMigrationBodyLength
}
