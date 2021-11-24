package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

const batchDataLength = 185

var BatchPrefix = getBadgerHoldPrefix(Batch{})

type Batch struct {
	ID                Uint256
	Type              batchtype.BatchType
	TransactionHash   common.Hash
	Hash              *common.Hash `badgerhold:"index"` // root of tree containing all commitments included in this batch
	FinalisationBlock *uint32
	AccountTreeRoot   *common.Hash
	PrevStateRoot     *common.Hash
	SubmissionTime    *Timestamp
}

func (b *Batch) Bytes() []byte {
	encoded := make([]byte, batchDataLength)
	copy(encoded[0:32], b.ID.Bytes())
	encoded[32] = byte(b.Type)
	copy(encoded[33:65], b.TransactionHash.Bytes())
	copy(encoded[65:98], EncodeHashPointer(b.Hash))
	copy(encoded[98:103], EncodeUint32Pointer(b.FinalisationBlock))
	copy(encoded[103:136], EncodeHashPointer(b.AccountTreeRoot))
	copy(encoded[136:169], EncodeHashPointer(b.PrevStateRoot))
	copy(encoded[169:185], encodeTimestampPointer(b.SubmissionTime))

	return encoded
}

func (b *Batch) SetBytes(data []byte) error {
	if len(data) != batchDataLength {
		return ErrInvalidLength
	}
	timestamp, err := decodeTimestampPointer(data[169:185])
	if err != nil {
		return err
	}

	b.ID.SetBytes(data[0:32])
	b.Type = batchtype.BatchType(data[32])
	b.TransactionHash.SetBytes(data[33:65])
	b.Hash = decodeHashPointer(data[65:98])
	b.FinalisationBlock = decodeUint32Pointer(data[98:103])
	b.AccountTreeRoot = decodeHashPointer(data[103:136])
	b.PrevStateRoot = decodeHashPointer(data[136:169])
	b.SubmissionTime = timestamp
	return nil
}
