package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
)

const TransactionBaseLength = 248

type TransactionBase struct {
	Hash         common.Hash
	TxType       txtype.TransactionType
	FromStateID  uint32
	Amount       Uint256
	Fee          Uint256
	Nonce        Uint256
	Signature    Signature
	ReceiveTime  *Timestamp
	BatchID      *Uint256 //TODO: use CommitmentID struct with badger
	IndexInBatch *uint8
	CommitmentID *CommitmentID
	ErrorMessage *string
}

func (t *TransactionBase) Bytes() []byte {
	b := make([]byte, t.BytesLen())
	copy(b[0:32], t.Hash.Bytes())
	b[32] = byte(t.TxType)
	binary.BigEndian.PutUint32(b[33:37], t.FromStateID)
	//TODO-tx: replace with only .Bytes() after merge
	copy(b[37:69], utils.PadLeft(t.Amount.Bytes(), 32))
	copy(b[69:101], utils.PadLeft(t.Fee.Bytes(), 32))
	copy(b[101:133], utils.PadLeft(t.Nonce.Bytes(), 32))
	copy(b[133:197], t.Signature.Bytes())
	copy(b[197:213], encodeTimestampPointer(t.ReceiveTime))
	copy(b[213:247], t.CommitmentID.PointerBytes())
	copy(b[247:], encodeStringPointer(t.ErrorMessage))

	return b
}

func (t *TransactionBase) SetBytes(data []byte) error {
	if len(data) < TransactionBaseLength {
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
	//TODO-tx: do the same in batch
	t.ReceiveTime, err = decodeTimestampPointer(data[197:213])
	if err != nil {
		return err
	}
	t.CommitmentID, err = decodeCommitmentIDPointer(data[213:247])
	if err != nil {
		return err
	}
	t.ErrorMessage = decodeStringPointer(data[247:])
	return nil
}

func (t *TransactionBase) BytesLen() int {
	if t.ErrorMessage != nil {
		return TransactionBaseLength + len(*t.ErrorMessage)
	}
	return TransactionBaseLength
}

type TransactionBaseForCommitment struct {
	Hash        common.Hash `db:"tx_hash"`
	FromStateID uint32      `db:"from_state_id"`
	Amount      Uint256
	Fee         Uint256
	Nonce       Uint256
	Signature   Signature
	ReceiveTime *Timestamp `db:"receive_time"`
}

func (t *TransactionBase) GetFromStateID() uint32 {
	return t.FromStateID
}

func (t *TransactionBase) GetAmount() Uint256 {
	return t.Amount
}

func (t *TransactionBase) GetFee() Uint256 {
	return t.Fee
}

func (t *TransactionBase) GetNonce() Uint256 {
	return t.Nonce
}

func (t *TransactionBase) SetNonce(nonce Uint256) {
	t.Nonce = nonce
}

func (t *TransactionBase) GetSignature() Signature {
	return t.Signature
}
