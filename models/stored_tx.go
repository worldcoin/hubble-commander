package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

const (
	immutableStoredTxBytesLength = 213
)

type StoredTx struct {
	Hash        common.Hash
	TxType      txtype.TransactionType
	FromStateID uint32
	Amount      Uint256
	Fee         Uint256
	Nonce       Uint256
	Signature   Signature
	ReceiveTime *Timestamp

	Body TxBody // ToStateID for transfer, ToPublicKey for C2T
}

func (t *StoredTx) Bytes() []byte {
	b := make([]byte, t.BytesLen())
	copy(b[0:32], t.Hash.Bytes())
	b[32] = byte(t.TxType)
	binary.BigEndian.PutUint32(b[33:37], t.FromStateID)
	copy(b[37:69], t.Amount.Bytes())
	copy(b[69:101], t.Fee.Bytes())
	copy(b[101:133], t.Nonce.Bytes())
	copy(b[133:197], t.Signature.Bytes())
	copy(b[197:213], encodeTimestampPointer(t.ReceiveTime))
	copy(b[213:], t.Body.Bytes())

	return b
}

func (t *StoredTx) SetBytes(data []byte) error {
	if len(data) < immutableStoredTxBytesLength {
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
	t.Body, err = transactionBody(data[213:], t.TxType)
	return err
}

func (t *StoredTx) BytesLen() int {
	return storedTransactionLength + t.Body.BytesLen()
}

type StoredTxReceipt struct {
	Hash         common.Hash
	TxType       txtype.TransactionType
	CommitmentID *CommitmentID
	ErrorMessage *string
	ToStateID    *uint32 // only for C2T
}

type TxBody interface {
	ByteEncoder
	BytesLen() int
}

type ImmutableTransferBody struct {
	ToStateID uint32
}

type ImmutableC2TBody struct {
	ToPublicKey PublicKey
}

type MutableC2TBody struct {
	ToStateID uint32
}
