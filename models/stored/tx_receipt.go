package stored

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

const txReceiptBytesLength = 72

var (
	TxReceiptName                = models.GetTypeName(TxReceipt{})
	TxReceiptPrefix              = models.GetBadgerHoldPrefix(TxReceipt{})
	errInvalidTxReceiptIndexType = fmt.Errorf("invalid stored.TxReceipt index type")
)

type TxReceipt struct {
	Hash         common.Hash
	CommitmentID *models.CommitmentID
	ToStateID    *uint32 // specified for C2Ts, nil for Transfers and MassMigrations
	ErrorMessage *string
}

func NewTxReceiptFromTransfer(t *models.Transfer) *TxReceipt {
	return &TxReceipt{
		Hash:         t.Hash,
		CommitmentID: t.CommitmentID,
		ErrorMessage: t.ErrorMessage,
	}
}

func NewTxReceiptFromCreate2Transfer(t *models.Create2Transfer) *TxReceipt {
	return &TxReceipt{
		Hash:         t.Hash,
		CommitmentID: t.CommitmentID,
		ToStateID:    t.ToStateID,
		ErrorMessage: t.ErrorMessage,
	}
}

func NewTxReceiptFromMassMigration(m *models.MassMigration) *TxReceipt {
	return &TxReceipt{
		Hash:         m.Hash,
		CommitmentID: m.CommitmentID,
		ErrorMessage: m.ErrorMessage,
	}
}

func (t *TxReceipt) Bytes() []byte {
	b := make([]byte, t.BytesLen())
	copy(b[0:32], t.Hash.Bytes())
	copy(b[32:66], models.EncodeCommitmentIDPointer(t.CommitmentID))
	copy(b[66:71], models.EncodeUint32Pointer(t.ToStateID))
	copy(b[71:], models.EncodeStringPointer(t.ErrorMessage))
	return b
}

func (t *TxReceipt) SetBytes(data []byte) error {
	if len(data) < txReceiptBytesLength {
		return models.ErrInvalidLength
	}
	commitmentID, err := models.DecodeCommitmentIDPointer(data[32:66])
	if err != nil {
		return err
	}

	t.Hash.SetBytes(data[0:32])
	t.CommitmentID = commitmentID
	t.ToStateID = models.DecodeUint32Pointer(data[66:71])
	t.ErrorMessage = models.DecodeStringPointer(data[71:])
	return nil
}

func (t *TxReceipt) BytesLen() int {
	if t.ErrorMessage != nil {
		return txReceiptBytesLength + len(*t.ErrorMessage)
	}
	return txReceiptBytesLength
}

// Type implements badgerhold.Storer
func (t TxReceipt) Type() string {
	return string(TxReceiptName)
}

// Indexes implements badgerhold.Storer
func (t TxReceipt) Indexes() map[string]bh.Index {
	return map[string]bh.Index{
		"CommitmentID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToTxReceipt(value)
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
				v, err := interfaceToTxReceipt(value)
				if err != nil {
					return nil, err
				}
				if v.ToStateID == nil {
					return nil, nil
				}
				return models.EncodeUint32(*v.ToStateID), nil
			},
		},
	}
}

func interfaceToTxReceipt(value interface{}) (*TxReceipt, error) {
	p, ok := value.(*TxReceipt)
	if ok {
		return p, nil
	}
	v, ok := value.(TxReceipt)
	if ok {
		return &v, nil
	}
	return nil, errors.WithStack(errInvalidTxReceiptIndexType)
}
