package stored

import (
	"bytes"
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

var (
	BatchedTxName                = models.GetTypeName(BatchedTx{})
	errInvalidBatchedTxIndexType = fmt.Errorf("invalid stored.BatchedTx index type")
)

type BatchedTx struct {
	PendingTx
	CommitmentID *models.CommitmentID
}

func NewBatchedTx(tx models.GenericTransaction) *BatchedTx {
	base := tx.GetBase()

	if base.CommitmentID == nil {
		// this is a PendingTx or maybe a FailedTx
		return nil
	}

	return &BatchedTx{
		PendingTx:    *NewPendingTx(tx),
		CommitmentID: base.CommitmentID,
	}
}

func NewBatchedTxFromPendingAndCommitment(pendingTx *PendingTx, commitmentID *models.CommitmentID) *BatchedTx {
	return &BatchedTx{
		PendingTx:    *pendingTx,
		CommitmentID: commitmentID,
	}
}

func (t *BatchedTx) Bytes() []byte {
	var buf bytes.Buffer

	bytesLen := t.BytesLen()
	buf.Grow(bytesLen)

	buf.Write(t.PendingTx.Bytes())
	buf.Write(EncodeCommitmentIDPointer(t.CommitmentID))

	return buf.Bytes()
}

func (t *BatchedTx) ToGenericTransaction() models.GenericTransaction {
	txn := t.PendingTx.ToGenericTransaction()
	txn.GetBase().CommitmentID = t.CommitmentID
	return txn
}

func (t *BatchedTx) ToTransfer() *models.Transfer {
	return t.ToGenericTransaction().ToTransfer()
}

func (t *BatchedTx) ToCreate2Transfer() *models.Create2Transfer {
	return t.ToGenericTransaction().ToCreate2Transfer()
}

func (t *BatchedTx) ToMassMigration() *models.MassMigration {
	return t.ToGenericTransaction().ToMassMigration()
}

// Careful: If there is a failure this will leave behind a partially-
//          populated PendinTx. If this gives you an error throw away the
//          PendingTx!
func (t *BatchedTx) SetBytes(data []byte) error {
	if len(data) < sizeBatchedTxNoBody {
		// This prevents obvious errors but it is still possible for this []byte
		// to be too short: it might not include a BatchedTx.PendingTx.Body
		return models.ErrInvalidLength
	}

	err := t.PendingTx.SetBytes(data)
	if err != nil {
		return err
	}

	// This relies on PendingTx.BytesLen() correctly reporting exactly how many bytes
	// were read in the call to SetBytes()
	_, rest := takeSlice(data, t.PendingTx.BytesLen())

	slice, _ := takeSlice(rest, sizeCommitment)
	commitmentID, err := decodeCommitmentIDPointer(slice)
	if err != nil {
		return err
	}
	t.CommitmentID = commitmentID

	return nil
}

func (t *BatchedTx) BytesLen() int {
	return t.PendingTx.BytesLen() + sizeCommitment
}

// nolint:gocritic
// implement badgerhold.Storer
func (t BatchedTx) Type() string {
	return string(BatchedTxName)
}

// nolint:gocritic
// implement badgerhold.Storer
func (t BatchedTx) Indexes() map[string]bh.Index {
	return map[string]bh.Index{
		"CommitmentID": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToBatchedTx(value)
				if err != nil {
					return nil, err
				}
				if v.CommitmentID == nil {
					panic("Every BatchedTx must have a CommitmentID")
				}
				return v.CommitmentID.Bytes(), nil
			},
		},
	}
}

func interfaceToBatchedTx(value interface{}) (*BatchedTx, error) {
	p, ok := value.(*BatchedTx)
	if ok {
		return p, nil
	}
	v, ok := value.(BatchedTx)
	if ok {
		return &v, nil
	}
	return nil, errors.WithStack(errInvalidBatchedTxIndexType)
}
