package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v3"
)

const (
	storedTransactionLength   = 248
	transferBodyLength        = 4
	create2TransferBodyLength = 133
)

var (
	StoredTransactionPrefix     = getBadgerHoldPrefix(StoredTransaction{})
	errInvalidStoredTxIndexType = errors.New("invalid StoredTx index type")
)

type StoredTransaction struct {
	Hash         common.Hash
	TxType       txtype.TransactionType
	FromStateID  uint32
	Amount       Uint256
	Fee          Uint256
	Nonce        Uint256
	Signature    Signature
	ReceiveTime  *Timestamp
	CommitmentID *CommitmentID
	ErrorMessage *string

	Body TransactionBody
}

func MakeStoredTransactionFromTransfer(t *Transfer) StoredTransaction {
	return StoredTransaction{
		Hash:         t.Hash,
		TxType:       t.TxType,
		FromStateID:  t.FromStateID,
		Amount:       t.Amount,
		Fee:          t.Fee,
		Nonce:        t.Nonce,
		Signature:    t.Signature,
		ReceiveTime:  t.ReceiveTime,
		CommitmentID: t.CommitmentID,
		ErrorMessage: t.ErrorMessage,
		Body: &StoredTransferBody{
			ToStateID: t.ToStateID,
		},
	}
}

func MakeStoredTransactionFromCreate2Transfer(t *Create2Transfer) StoredTransaction {
	return StoredTransaction{
		Hash:         t.Hash,
		TxType:       t.TxType,
		FromStateID:  t.FromStateID,
		Amount:       t.Amount,
		Fee:          t.Fee,
		Nonce:        t.Nonce,
		Signature:    t.Signature,
		ReceiveTime:  t.ReceiveTime,
		CommitmentID: t.CommitmentID,
		ErrorMessage: t.ErrorMessage,
		Body: &StoredCreate2TransferBody{
			ToStateID:   t.ToStateID,
			ToPublicKey: t.ToPublicKey,
		},
	}
}

func (t *StoredTransaction) ToTransfer() *Transfer {
	transferBody, ok := t.Body.(*StoredTransferBody)
	if !ok {
		panic("invalid transfer body type")
	}

	return &Transfer{
		TransactionBase: TransactionBase{
			Hash:         t.Hash,
			TxType:       t.TxType,
			FromStateID:  t.FromStateID,
			Amount:       t.Amount,
			Fee:          t.Fee,
			Nonce:        t.Nonce,
			Signature:    t.Signature,
			ReceiveTime:  t.ReceiveTime,
			CommitmentID: t.CommitmentID,
			ErrorMessage: t.ErrorMessage,
		},
		ToStateID: transferBody.ToStateID,
	}
}

func (t *StoredTransaction) ToCreate2Transfer() *Create2Transfer {
	c2tBody, ok := t.Body.(*StoredCreate2TransferBody)
	if !ok {
		panic("invalid create2Transfer body type")
	}

	return &Create2Transfer{
		TransactionBase: TransactionBase{
			Hash:         t.Hash,
			TxType:       t.TxType,
			FromStateID:  t.FromStateID,
			Amount:       t.Amount,
			Fee:          t.Fee,
			Nonce:        t.Nonce,
			Signature:    t.Signature,
			ReceiveTime:  t.ReceiveTime,
			CommitmentID: t.CommitmentID,
			ErrorMessage: t.ErrorMessage,
		},
		ToStateID:   c2tBody.ToStateID,
		ToPublicKey: c2tBody.ToPublicKey,
	}
}

func (t *StoredTransaction) Bytes() []byte {
	b := make([]byte, t.BytesLen())
	copy(b[0:32], t.Hash.Bytes())
	b[32] = byte(t.TxType)
	binary.BigEndian.PutUint32(b[33:37], t.FromStateID)
	copy(b[37:69], t.Amount.Bytes())
	copy(b[69:101], t.Fee.Bytes())
	copy(b[101:133], t.Nonce.Bytes())
	copy(b[133:197], t.Signature.Bytes())
	copy(b[197:213], encodeTimestampPointer(t.ReceiveTime))
	copy(b[213:247], EncodeCommitmentIDPointer(t.CommitmentID))
	copy(b[247:], t.Body.Bytes())
	copy(b[247+t.Body.BytesLen():], encodeStringPointer(t.ErrorMessage))

	return b
}

func (t *StoredTransaction) SetBytes(data []byte) error {
	if len(data) < storedTransactionLength {
		return ErrInvalidLength
	}

	t.Hash.SetBytes(data[0:32])
	t.TxType = txtype.TransactionType(data[32])
	t.FromStateID = binary.BigEndian.Uint32(data[33:37])
	t.Amount.SetBytes(data[37:69])
	t.Fee.SetBytes(data[69:101])
	t.Nonce.SetBytes(data[101:133])
	err := t.Signature.SetBytes(data[133:197])
	if err != nil {
		return err
	}
	t.ReceiveTime, err = decodeTimestampPointer(data[197:213])
	if err != nil {
		return err
	}
	t.CommitmentID, err = decodeCommitmentIDPointer(data[213:247])
	if err != nil {
		return err
	}
	t.Body, err = transactionBody(data[247:], t.TxType)
	if err != nil {
		return err
	}
	t.ErrorMessage = decodeStringPointer(data[247+t.Body.BytesLen():])
	return nil
}

func (t *StoredTransaction) BytesLen() int {
	length := storedTransactionLength + t.Body.BytesLen()
	if t.ErrorMessage != nil {
		length += len(*t.ErrorMessage)
	}
	return length
}

// nolint:gocritic
// Type implements badgerhold.Storer
func (t StoredTransaction) Type() string {
	return string(StoredTransactionPrefix[3:])
}

// nolint:gocritic
// Indexes implements badgerhold.Storer
func (t StoredTransaction) Indexes() map[string]bh.Index {
	return map[string]bh.Index{
		"FromStateID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredTransaction(value)
				if err != nil {
					return nil, err
				}
				return EncodeUint32(&v.FromStateID)
			},
		},
		"CommitmentID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredTransaction(value)
				if err != nil {
					return nil, err
				}
				return EncodeCommitmentIDPointer(v.CommitmentID), nil
			},
		},
		"ToStateID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredTransaction(value)
				if err != nil {
					return nil, err
				}
				return EncodeUint32Pointer(v.Body.GetToStateID()), nil
			},
		},
	}
}

func transactionBody(data []byte, transactionType txtype.TransactionType) (TransactionBody, error) {
	switch transactionType {
	case txtype.Transfer:
		body := new(StoredTransferBody)
		err := body.SetBytes(data)
		return body, err
	case txtype.Create2Transfer:
		body := new(StoredCreate2TransferBody)
		err := body.SetBytes(data)
		return body, err
	case txtype.Genesis, txtype.MassMigration:
		return nil, errors.Errorf("unsupported tx type: %s", transactionType)
	}
	return nil, nil
}

func interfaceToStoredTransaction(value interface{}) (*StoredTransaction, error) {
	p, ok := value.(*StoredTransaction)
	if ok {
		return p, nil
	}
	v, ok := value.(StoredTransaction)
	if ok {
		return &v, nil
	}
	return nil, errInvalidStoredTxIndexType
}

type TransactionBody interface {
	ByteEncoder
	BytesLen() int
	GetToStateID() *uint32
}

type StoredTransferBody struct {
	ToStateID uint32
}

func (t *StoredTransferBody) Bytes() []byte {
	b := make([]byte, transferBodyLength)
	binary.BigEndian.PutUint32(b, t.ToStateID)
	return b
}

func (t *StoredTransferBody) SetBytes(data []byte) error {
	t.ToStateID = binary.BigEndian.Uint32(data)
	return nil
}

func (t *StoredTransferBody) BytesLen() int {
	return transferBodyLength
}

func (t *StoredTransferBody) GetToStateID() *uint32 {
	return &t.ToStateID
}

type StoredCreate2TransferBody struct {
	ToStateID   *uint32
	ToPublicKey PublicKey
}

func (t *StoredCreate2TransferBody) Bytes() []byte {
	b := make([]byte, create2TransferBodyLength)
	copy(b[:128], t.ToPublicKey.Bytes())
	copy(b[128:133], EncodeUint32Pointer(t.ToStateID))
	return b
}

func (t *StoredCreate2TransferBody) SetBytes(data []byte) error {
	err := t.ToPublicKey.SetBytes(data[:128])
	if err != nil {
		return err
	}
	t.ToStateID = decodeUint32Pointer(data[128:133])
	return nil
}

func (t *StoredCreate2TransferBody) BytesLen() int {
	return create2TransferBodyLength
}

func (t *StoredCreate2TransferBody) GetToStateID() *uint32 {
	return t.ToStateID
}
