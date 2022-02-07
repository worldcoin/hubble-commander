package stored

import (
	"bytes"
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type PendingTx struct {
	Hash common.Hash

	TxType      txtype.TransactionType
	FromStateID uint32
	Amount      models.Uint256
	Fee         models.Uint256
	Nonce       models.Uint256
	Signature   models.Signature
	ReceiveTime *models.Timestamp

	Body TxBody
}

func NewPendingTx(tx models.GenericTransaction) *PendingTx {
	base := tx.GetBase()
	return &PendingTx{
		Hash:        base.Hash,
		TxType:      base.TxType,
		FromStateID: base.FromStateID,
		Amount:      base.Amount,
		Fee:         base.Fee,
		Nonce:       base.Nonce,
		Signature:   base.Signature,
		ReceiveTime: base.ReceiveTime,
		Body:        NewTxBody(tx),
	}
}

func (t *PendingTx) ToGenericTransaction() models.GenericTransaction {
	base := t.ToBase()
	return t.Body.ToGenericTransaction(base)
}

func (t *PendingTx) ToBase() *models.TransactionBase {
	return &models.TransactionBase{
		Hash:        t.Hash,
		TxType:      t.TxType,
		FromStateID: t.FromStateID,
		Amount:      t.Amount,
		Fee:         t.Fee,
		Nonce:       t.Nonce,
		Signature:   t.Signature,
		ReceiveTime: t.ReceiveTime,

		// These are set for BatchedTx, PendingTx are defined by not having these
		CommitmentID: nil,
		ErrorMessage: nil,
	}
}

func (t *PendingTx) ToTransfer() *models.Transfer {
	return t.ToGenericTransaction().ToTransfer()
}

func (t *PendingTx) ToCreate2Transfer() *models.Create2Transfer {
	return t.ToGenericTransaction().ToCreate2Transfer()
}

func (t *PendingTx) ToMassMigration() *models.MassMigration {
	return t.ToGenericTransaction().ToMassMigration()
}

func (t *PendingTx) Bytes() []byte {
	var buf bytes.Buffer

	bytesLen := t.BytesLen()
	buf.Grow(bytesLen)

	buf.Write(t.Hash.Bytes())
	buf.WriteByte(byte(t.TxType))
	buf.Write(EncodeUint32(t.FromStateID))
	buf.Write(t.Amount.Bytes())
	buf.Write(t.Fee.Bytes())
	buf.Write(t.Nonce.Bytes())
	buf.Write(t.Signature.Bytes())
	buf.Write(encodeTimestampPointer(t.ReceiveTime))
	buf.Write(t.Body.Bytes())

	return buf.Bytes()
}

func takeSlice(data []byte, count int) (slice, rest []byte) {
	return data[:count], data[count:]
}

// Careful: If there is a failure this will leave behind a partially-
//          populated PendinTx. If this gives you an error throw away the
//          PendingTx!
func (t *PendingTx) SetBytes(data []byte) error {
	if len(data) < sizePendingTxNoBody {
		return models.ErrInvalidLength
	}

	slice, rest := takeSlice(data, sizeHash)
	t.Hash.SetBytes(slice)

	slice, rest = takeSlice(rest, sizeTxType)
	t.TxType = txtype.TransactionType(slice[0])

	slice, rest = takeSlice(rest, sizeU32)
	t.FromStateID = binary.BigEndian.Uint32(slice)

	slice, rest = takeSlice(rest, sizeU256)
	t.Amount.SetBytes(slice)

	slice, rest = takeSlice(rest, sizeU256)
	t.Fee.SetBytes(slice)

	slice, rest = takeSlice(rest, sizeU256)
	t.Nonce.SetBytes(slice)

	slice, rest = takeSlice(rest, sizeSignature)
	err := t.Signature.SetBytes(slice)
	if err != nil {
		return err
	}

	slice, rest = takeSlice(rest, sizeTimestamp)
	receiveTime, err := decodeTimestampPointer(slice)
	if err != nil {
		return err
	}
	t.ReceiveTime = receiveTime

	bodyLen := expectedBodyBytesLen(t.TxType)
	slice, _ = takeSlice(rest, bodyLen)
	body, err := txBody(slice, t.TxType)
	if err != nil {
		return err
	}
	t.Body = body

	// Note that if we are given too many bytes then they are silently ignored.
	// BatchedTx.SetBytes and FailedTx.SetBytes depend on this behavior.

	return nil
}

func expectedBodyBytesLen(txType txtype.TransactionType) int {
	switch txType {
	case txtype.Transfer:
		return txTransferBodyLength
	case txtype.Create2Transfer:
		return txCreate2TransferBodyLength
	case txtype.MassMigration:
		return txMassMigrationBodyLength
	}

	panic("unexpected transaction type")
}

func (t *PendingTx) BytesLen() int {
	return sizePendingTxNoBody + t.Body.BytesLen()
}
