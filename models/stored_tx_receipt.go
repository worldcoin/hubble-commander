package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

const storedTxReceiptBytesLength = 72

var (
	StoredTxReceiptName                = getTypeName(StoredTxReceipt{})
	StoredTxReceiptPrefix              = getBadgerHoldPrefix(StoredTxReceipt{})
	errInvalidStoredTxReceiptIndexType = errors.New("invalid StoredTxReceipt index type")
)

type StoredTxReceipt struct {
	Hash         common.Hash
	CommitmentID *CommitmentID
	ToStateID    *uint32 // specified for C2Ts, nil for Transfers and MassMigrations
	ErrorMessage *string
}

func NewStoredTxReceiptFromTransfer(t *Transfer) *StoredTxReceipt {
	return &StoredTxReceipt{
		Hash:         t.Hash,
		CommitmentID: t.CommitmentID,
		ErrorMessage: t.ErrorMessage,
	}
}

func NewStoredTxReceiptFromCreate2Transfer(t *Create2Transfer) *StoredTxReceipt {
	return &StoredTxReceipt{
		Hash:         t.Hash,
		CommitmentID: t.CommitmentID,
		ToStateID:    t.ToStateID,
		ErrorMessage: t.ErrorMessage,
	}
}

func NewStoredTxReceiptFromMassMigration(m *MassMigration) *StoredTxReceipt {
	return &StoredTxReceipt{
		Hash:         m.Hash,
		CommitmentID: m.CommitmentID,
		ErrorMessage: m.ErrorMessage,
	}
}

func (t *StoredTxReceipt) Bytes() []byte {
	b := make([]byte, t.BytesLen())
	copy(b[0:32], t.Hash.Bytes())
	copy(b[32:66], EncodeCommitmentIDPointer(t.CommitmentID))
	copy(b[66:71], EncodeUint32Pointer(t.ToStateID))
	copy(b[71:], encodeStringPointer(t.ErrorMessage))
	return b
}

func (t *StoredTxReceipt) SetBytes(data []byte) error {
	if len(data) < storedTxReceiptBytesLength {
		return ErrInvalidLength
	}
	commitmentID, err := decodeCommitmentIDPointer(data[32:66])
	if err != nil {
		return err
	}

	t.Hash.SetBytes(data[0:32])
	t.CommitmentID = commitmentID
	t.ToStateID = decodeUint32Pointer(data[66:71])
	t.ErrorMessage = decodeStringPointer(data[71:])
	return nil
}

func (t *StoredTxReceipt) BytesLen() int {
	if t.ErrorMessage != nil {
		return storedTxReceiptBytesLength + len(*t.ErrorMessage)
	}
	return storedTxReceiptBytesLength
}

// Type implements badgerhold.Storer
func (t StoredTxReceipt) Type() string {
	return string(StoredTxReceiptName)
}

// Indexes implements badgerhold.Storer
func (t StoredTxReceipt) Indexes() map[string]bh.Index {
	return map[string]bh.Index{
		"CommitmentID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredTxReceipt(value)
				if err != nil {
					return nil, err
				}
				if v.CommitmentID == nil {
					return nil, nil
				}
				return v.CommitmentID.Bytes(), nil
			},
		},
		"ToStateID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredTxReceipt(value)
				if err != nil {
					return nil, err
				}
				if v.ToStateID == nil {
					return nil, nil
				}
				return EncodeUint32(*v.ToStateID), nil
			},
		},
	}
}

func interfaceToStoredTxReceipt(value interface{}) (*StoredTxReceipt, error) {
	p, ok := value.(*StoredTxReceipt)
	if ok {
		return p, nil
	}
	v, ok := value.(StoredTxReceipt)
	if ok {
		return &v, nil
	}
	return nil, errInvalidStoredTxReceiptIndexType
}
