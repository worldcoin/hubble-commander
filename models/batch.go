package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

const batchDataLength = 153

type Batch struct {
	ID                Uint256 `db:"batch_id"`
	Type              txtype.TransactionType
	TransactionHash   common.Hash  `db:"transaction_hash"`
	Hash              *common.Hash `db:"batch_hash"`         // root of tree containing all commitments included in this batch
	FinalisationBlock *uint32      `db:"finalisation_block"` // nolint:misspell
	AccountTreeRoot   *common.Hash `db:"account_tree_root"`
	PrevStateRoot     *common.Hash `db:"prev_state_root"`
	SubmissionTime    *Timestamp   `db:"submission_time"`
}

func (b *Batch) Bytes() []byte {
	encoded := make([]byte, batchDataLength)
	encoded[0] = byte(b.Type)
	copy(encoded[1:33], b.TransactionHash[:])
	copy(encoded[33:66], encodeHashPointer(b.Hash))
	copy(encoded[66:71], encodeUint32Pointer(b.FinalisationBlock))
	copy(encoded[71:104], encodeHashPointer(b.AccountTreeRoot))
	copy(encoded[104:137], encodeHashPointer(b.PrevStateRoot))
	copy(encoded[137:153], encodePointer(15, b.SubmissionTime))

	return encoded
}

func (b *Batch) SetBytes(data []byte) error {
	if len(data) != batchDataLength {
		return ErrInvalidLength
	}

	b.Type = txtype.TransactionType(data[0])
	b.TransactionHash.SetBytes(data[1:33])
	b.Hash = decodeHashPointer(data[33:66])
	b.FinalisationBlock = decodeUint32Pointer(data[66:71])
	b.AccountTreeRoot = decodeHashPointer(data[71:104])
	b.PrevStateRoot = decodeHashPointer(data[104:137])

	timestamp, err := decodeTimestampPointer(data[137:153])
	if err != nil {
		return err
	}
	b.SubmissionTime = timestamp
	return nil
}
