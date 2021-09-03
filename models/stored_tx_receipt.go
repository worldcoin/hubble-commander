package models

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v3"
)

const storedReceiptBytesLength = 72

var (
	StoredReceiptPrefix              = getBadgerHoldPrefix(StoredReceipt{})
	errInvalidStoredReceiptIndexType = errors.New("invalid StoredReceipt index type")
)

type StoredReceipt struct {
	Hash         common.Hash
	CommitmentID *CommitmentID
	ToStateID    *uint32 // only for C2T
	ErrorMessage *string
}

func MakeStoredReceiptFromTransfer(t *Transfer) StoredReceipt {
	return StoredReceipt{
		Hash:         t.Hash,
		CommitmentID: t.CommitmentID,
		ErrorMessage: t.ErrorMessage,
	}
}

func MakeStoredReceiptFromCreate2Transfer(t *Create2Transfer) StoredReceipt {
	return StoredReceipt{
		Hash:         t.Hash,
		CommitmentID: t.CommitmentID,
		ToStateID:    t.ToStateID,
		ErrorMessage: t.ErrorMessage,
	}
}

func (t *StoredReceipt) Bytes() []byte {
	b := make([]byte, t.BytesLen())
	copy(b[0:32], t.Hash.Bytes())
	copy(b[32:66], EncodeCommitmentIDPointer(t.CommitmentID))
	copy(b[66:71], EncodeUint32Pointer(t.ToStateID))
	copy(b[71:], encodeStringPointer(t.ErrorMessage))
	return b
}

func (t *StoredReceipt) SetBytes(data []byte) (err error) {
	if len(data) < storedReceiptBytesLength {
		return ErrInvalidLength
	}

	t.Hash.SetBytes(data[0:32])
	t.CommitmentID, err = decodeCommitmentIDPointer(data[32:66])
	if err != nil {
		return err
	}
	t.ToStateID = decodeUint32Pointer(data[66:71])
	t.ErrorMessage = decodeStringPointer(data[71:])
	return nil
}

func (t *StoredReceipt) BytesLen() int {
	if t.ErrorMessage != nil {
		return storedReceiptBytesLength + len(*t.ErrorMessage)
	}
	return storedReceiptBytesLength
}

// Type implements badgerhold.Storer
func (t StoredReceipt) Type() string {
	return string(StoredReceiptPrefix[3:])
}

// Indexes implements badgerhold.Storer
func (t StoredReceipt) Indexes() map[string]bh.Index {
	return map[string]bh.Index{
		"CommitmentID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredReceipt(value)
				if err != nil {
					return nil, err
				}
				return EncodeCommitmentIDPointer(v.CommitmentID), nil
			},
		},
		"ToStateID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredReceipt(value)
				if err != nil {
					return nil, err
				}

				if v.ToStateID == nil {
					return nil, nil
				}
				return EncodeUint32(v.ToStateID)
			},
		},
	}
}

func interfaceToStoredReceipt(value interface{}) (*StoredReceipt, error) {
	p, ok := value.(*StoredReceipt)
	if ok {
		return p, nil
	}
	v, ok := value.(StoredReceipt)
	if ok {
		return &v, nil
	}
	return nil, errInvalidStoredReceiptIndexType
}
