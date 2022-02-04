package stored

import (
	"bytes"

	"github.com/Worldcoin/hubble-commander/models"
)

type FailedTx struct {
	PendingTx

	ErrorMessage string
}

func NewFailedTx(tx models.GenericTransaction) *FailedTx {
	errorMessage := tx.GetBase().ErrorMessage
	if errorMessage == nil {
		panic("missing ErrorMessage in param passed to NewFailedTx")
	}

	return &FailedTx{
		PendingTx:    *NewPendingTx(tx),
		ErrorMessage: *errorMessage,
	}
}

func NewFailedTxFromError(pendingTx *PendingTx, errorMessage string) *FailedTx {
	return &FailedTx{
		PendingTx:    *pendingTx,
		ErrorMessage: errorMessage,
	}
}

func (t *FailedTx) ToGenericTransaction() models.GenericTransaction {
	txn := t.PendingTx.ToGenericTransaction()
	txn.GetBase().ErrorMessage = &t.ErrorMessage
	return txn
}

func (t *FailedTx) Bytes() []byte {
	var buf bytes.Buffer

	bytesLen := t.BytesLen()
	buf.Grow(bytesLen)

	buf.Write(t.PendingTx.Bytes())
	buf.Write([]byte(t.ErrorMessage))

	return buf.Bytes()
}

func (t *FailedTx) SetBytes(data []byte) error {
	if len(data) < sizePendingTxNoBody {
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
	t.ErrorMessage = string(rest)

	return nil
}

func (t *FailedTx) BytesLen() int {
	return t.PendingTx.BytesLen() + len(t.ErrorMessage)
}
