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

// TODO: test that this works, that we're not setting fields on a copy of the base
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
// TODO: Turn this into a constructor which returns (Option<BatchedTx>, Option<err>)
func (t *BatchedTx) SetBytes(data []byte) error {
	if len(data) < sizePendingTx {
		// TODO: What is the correct size to check for?
		return models.ErrInvalidLength
	}

	err := t.PendingTx.SetBytes(data)
	if err != nil {
		return err
	}

	// TODO: this code relies on there being a 1-to-1 mapping between internal
	//       states and serializations. This makes the code a little brittle!
	//       Better would be for `SetBytes` to return the remaining slice.
	//       ( it assumes len(x) == len(serialize(deserialize(x)) )
	_, rest := takeSlice(data, t.PendingTx.BytesLen())

	// see EncodeCommitmentIDPointer for why this length is chosen
	slice, _ := takeSlice(rest, models.CommitmentIDDataLength+1)
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
					// TODO: now that FailedTx is broken out,
					//       this never fires, right?
					return nil, nil
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
