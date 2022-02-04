package stored

import (
	"bytes"

	"github.com/Worldcoin/hubble-commander/models"
)

var (
	BatchedTxName   = models.GetTypeName(BatchedTx{})
	BatchedTxPrefix = models.GetBadgerHoldPrefix(BatchedTx{})
)

type BatchedTx struct {
	PendingTx
	CommitmentID models.CommitmentID `badgerhold:"index"`
}

func NewBatchedTx(tx models.GenericTransaction) *BatchedTx {
	base := tx.GetBase()

	if base.CommitmentID == nil {
		// this is a PendingTx or maybe a FailedTx
		return nil
	}

	return &BatchedTx{
		PendingTx:    *NewPendingTx(tx),
		CommitmentID: *base.CommitmentID,
	}
}

func NewBatchedTxFromPendingAndCommitment(pendingTx *PendingTx, commitmentID *models.CommitmentID) *BatchedTx {
	return &BatchedTx{
		PendingTx:    *pendingTx,
		CommitmentID: *commitmentID,
	}
}

func (t *BatchedTx) Bytes() []byte {
	var buf bytes.Buffer

	bytesLen := t.BytesLen()
	buf.Grow(bytesLen)

	buf.Write(t.PendingTx.Bytes())
	buf.Write(t.CommitmentID.Bytes())

	return buf.Bytes()
}

func (t *BatchedTx) ToGenericTransaction() models.GenericTransaction {
	txn := t.PendingTx.ToGenericTransaction()
	txn.GetBase().CommitmentID = &t.CommitmentID
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
	return t.CommitmentID.SetBytes(slice)
}

func (t *BatchedTx) BytesLen() int {
	return t.PendingTx.BytesLen() + sizeCommitment
}
