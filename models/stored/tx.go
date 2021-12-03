package stored

import (
	"encoding/binary"
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

const (
	txBytesLength               = 213
	txTransferBodyLength        = 4
	txCreate2TransferBodyLength = models.PublicKeyLength
	txMassMigrationBodyLength   = 4
)

var (
	TxName                = models.GetTypeName(Tx{})
	TxPrefix              = models.GetBadgerHoldPrefix(Tx{})
	errInvalidTxIndexType = fmt.Errorf("invalid stored.Tx index type")
)

type Tx struct {
	Hash        common.Hash
	TxType      txtype.TransactionType
	FromStateID uint32
	Amount      models.Uint256
	Fee         models.Uint256
	Nonce       models.Uint256
	Signature   models.Signature
	ReceiveTime *models.Timestamp

	Body TxBody
}

func NewTxFromTransfer(t *models.Transfer) *Tx {
	return &Tx{
		Hash:        t.Hash,
		TxType:      t.TxType,
		FromStateID: t.FromStateID,
		Amount:      t.Amount,
		Fee:         t.Fee,
		Nonce:       t.Nonce,
		Signature:   t.Signature,
		ReceiveTime: t.ReceiveTime,
		Body: &TxTransferBody{
			ToStateID: t.ToStateID,
		},
	}
}

func NewTxFromCreate2Transfer(t *models.Create2Transfer) *Tx {
	return &Tx{
		Hash:        t.Hash,
		TxType:      t.TxType,
		FromStateID: t.FromStateID,
		Amount:      t.Amount,
		Fee:         t.Fee,
		Nonce:       t.Nonce,
		Signature:   t.Signature,
		ReceiveTime: t.ReceiveTime,
		Body: &TxCreate2TransferBody{
			ToPublicKey: t.ToPublicKey,
		},
	}
}

func NewTxFromMassMigration(m *models.MassMigration) *Tx {
	return &Tx{
		Hash:        m.Hash,
		TxType:      m.TxType,
		FromStateID: m.FromStateID,
		Amount:      m.Amount,
		Fee:         m.Fee,
		Nonce:       m.Nonce,
		Signature:   m.Signature,
		ReceiveTime: m.ReceiveTime,
		Body: &TxMassMigrationBody{
			SpokeID: m.SpokeID,
		},
	}
}

func (t *Tx) Bytes() []byte {
	b := make([]byte, t.BytesLen())
	copy(b[0:32], t.Hash.Bytes())
	b[32] = byte(t.TxType)
	binary.BigEndian.PutUint32(b[33:37], t.FromStateID)
	copy(b[37:69], t.Amount.Bytes())
	copy(b[69:101], t.Fee.Bytes())
	copy(b[101:133], t.Nonce.Bytes())
	copy(b[133:197], t.Signature.Bytes())
	copy(b[197:213], models.EncodeTimestampPointer(t.ReceiveTime))
	copy(b[213:], t.Body.Bytes())

	return b
}

func (t *Tx) SetBytes(data []byte) error {
	if len(data) < txBytesLength {
		return models.ErrInvalidLength
	}
	err := t.Signature.SetBytes(data[133:197])
	if err != nil {
		return err
	}
	receiveTime, err := models.DecodeTimestampPointer(data[197:213])
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

func (t *Tx) BytesLen() int {
	return txBytesLength + t.Body.BytesLen()
}

func (t *Tx) ToTransfer(txReceipt *TxReceipt) *models.Transfer {
	transferBody, ok := t.Body.(*TxTransferBody)
	if !ok {
		panic("invalid transfer body type")
	}

	transfer := &models.Transfer{
		TransactionBase: models.TransactionBase{
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

func (t *Tx) ToCreate2Transfer(txReceipt *TxReceipt) *models.Create2Transfer {
	c2tBody, ok := t.Body.(*TxCreate2TransferBody)
	if !ok {
		panic("invalid create2Transfer body type")
	}

	transfer := &models.Create2Transfer{
		TransactionBase: models.TransactionBase{
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

func (t *Tx) ToMassMigration(txReceipt *TxReceipt) *models.MassMigration {
	massMigrationBody, ok := t.Body.(*TxMassMigrationBody)
	if !ok {
		panic("invalid mass migration body type")
	}

	massMigration := &models.MassMigration{
		TransactionBase: models.TransactionBase{
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

// nolint:gocritic
// Type implements badgerhold.Storer
func (t Tx) Type() string {
	return string(TxName)
}

// nolint:gocritic
// Indexes implements badgerhold.Storer
func (t Tx) Indexes() map[string]bh.Index {
	return map[string]bh.Index{
		"FromStateID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredTx(value)
				if err != nil {
					return nil, err
				}
				return models.EncodeUint32(v.FromStateID), nil
			},
		},
		"ToStateID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredTx(value)
				if err != nil {
					return nil, err
				}

				transferBody, ok := v.Body.(*TxTransferBody)
				if !ok {
					return nil, nil
				}
				return models.EncodeUint32(transferBody.ToStateID), nil
			},
		},
	}
}

func interfaceToStoredTx(value interface{}) (*Tx, error) {
	p, ok := value.(*Tx)
	if ok {
		return p, nil
	}
	v, ok := value.(Tx)
	if ok {
		return &v, nil
	}
	return nil, errors.WithStack(errInvalidTxIndexType)
}

type TxBody interface {
	models.ByteEncoder
	BytesLen() int
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

type TxCreate2TransferBody struct {
	ToPublicKey models.PublicKey
}

func (t *TxCreate2TransferBody) Bytes() []byte {
	return t.ToPublicKey.Bytes()
}

func (t *TxCreate2TransferBody) SetBytes(data []byte) error {
	return t.ToPublicKey.SetBytes(data)
}

func (t *TxCreate2TransferBody) BytesLen() int {
	return txCreate2TransferBodyLength
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
