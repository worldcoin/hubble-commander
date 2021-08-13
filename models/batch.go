package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
)

const batchDataLength = 185

type Batch struct {
	ID                Uint256 `db:"batch_id"`
	Type              txtype.TransactionType
	TransactionHash   common.Hash  `db:"transaction_hash"`
	Hash              *common.Hash `db:"batch_hash" badgerhold:"index"` // root of tree containing all commitments included in this batch
	FinalisationBlock *uint32      `db:"finalisation_block"`            // nolint:misspell
	AccountTreeRoot   *common.Hash `db:"account_tree_root"`
	PrevStateRoot     *common.Hash `db:"prev_state_root"`
	SubmissionTime    *Timestamp   `db:"submission_time"`
}

func (b *Batch) Bytes() []byte {
	encoded := make([]byte, batchDataLength)
	copy(encoded[0:32], utils.PadLeft(b.ID.Bytes(), 32))
	encoded[32] = byte(b.Type)
	copy(encoded[33:65], b.TransactionHash[:])
	copy(encoded[65:98], encodeHashPointer(b.Hash))
	copy(encoded[98:103], encodeUint32Pointer(b.FinalisationBlock))
	copy(encoded[103:136], encodeHashPointer(b.AccountTreeRoot))
	copy(encoded[136:169], encodeHashPointer(b.PrevStateRoot))
	copy(encoded[169:185], encodePointer(15, b.SubmissionTime))

	return encoded
}

func (b *Batch) SetBytes(data []byte) error {
	if len(data) != batchDataLength {
		return ErrInvalidLength
	}

	b.ID.SetBytes(data[0:32])
	b.Type = txtype.TransactionType(data[32])
	b.TransactionHash.SetBytes(data[33:65])
	b.Hash = decodeHashPointer(data[65:98])
	b.FinalisationBlock = decodeUint32Pointer(data[98:103])
	b.AccountTreeRoot = decodeHashPointer(data[103:136])
	b.PrevStateRoot = decodeHashPointer(data[136:169])

	timestamp, err := decodeTimestampPointer(data[169:185])
	if err != nil {
		return err
	}
	b.SubmissionTime = timestamp
	return nil
}
