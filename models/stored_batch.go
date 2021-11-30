package models

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

const storedBatchDataLength = 185

var (
	StoredBatchName                = getTypeName(StoredBatch{})
	StoredBatchPrefix              = getBadgerHoldPrefix(StoredBatch{})
	errInvalidStoredBatchIndexType = fmt.Errorf("invalid StoredBatch index type")
)

type StoredBatch struct {
	ID                Uint256
	BType             batchtype.BatchType
	TransactionHash   common.Hash
	Hash              *common.Hash // root of tree containing all commitments included in this batch
	FinalisationBlock *uint32
	AccountTreeRoot   *common.Hash
	PrevStateRoot     *common.Hash
	SubmissionTime    *Timestamp
}

func NewStoredBatchFromBatch(b *Batch) *StoredBatch {
	return &StoredBatch{
		ID:                b.ID,
		BType:             b.Type,
		TransactionHash:   b.TransactionHash,
		Hash:              b.Hash,
		FinalisationBlock: b.FinalisationBlock,
		AccountTreeRoot:   b.AccountTreeRoot,
		PrevStateRoot:     b.PrevStateRoot,
		SubmissionTime:    b.SubmissionTime,
	}
}

func (b *StoredBatch) ToBatch() *Batch {
	return &Batch{
		ID:                b.ID,
		Type:              b.BType,
		TransactionHash:   b.TransactionHash,
		Hash:              b.Hash,
		FinalisationBlock: b.FinalisationBlock,
		AccountTreeRoot:   b.AccountTreeRoot,
		PrevStateRoot:     b.PrevStateRoot,
		SubmissionTime:    b.SubmissionTime,
	}
}

func (b *StoredBatch) Bytes() []byte {
	encoded := make([]byte, storedBatchDataLength)
	copy(encoded[0:32], b.ID.Bytes())
	encoded[32] = byte(b.BType)
	copy(encoded[33:65], b.TransactionHash.Bytes())
	copy(encoded[65:98], EncodeHashPointer(b.Hash))
	copy(encoded[98:103], EncodeUint32Pointer(b.FinalisationBlock))
	copy(encoded[103:136], EncodeHashPointer(b.AccountTreeRoot))
	copy(encoded[136:169], EncodeHashPointer(b.PrevStateRoot))
	copy(encoded[169:185], encodeTimestampPointer(b.SubmissionTime))

	return encoded
}

func (b *StoredBatch) SetBytes(data []byte) error {
	if len(data) != storedBatchDataLength {
		return ErrInvalidLength
	}
	timestamp, err := decodeTimestampPointer(data[169:185])
	if err != nil {
		return err
	}

	b.ID.SetBytes(data[0:32])
	b.BType = batchtype.BatchType(data[32])
	b.TransactionHash.SetBytes(data[33:65])
	b.Hash = decodeHashPointer(data[65:98])
	b.FinalisationBlock = decodeUint32Pointer(data[98:103])
	b.AccountTreeRoot = decodeHashPointer(data[103:136])
	b.PrevStateRoot = decodeHashPointer(data[136:169])
	b.SubmissionTime = timestamp
	return nil
}

// nolint:gocritic
// Type implements badgerhold.Storer
func (b StoredBatch) Type() string {
	return string(StoredBatchName)
}

// nolint:gocritic
// Indexes implements badgerhold.Storer
func (b StoredBatch) Indexes() map[string]bh.Index {
	return map[string]bh.Index{
		"Hash": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToStoredBatch(value)
				if err != nil {
					return nil, err
				}
				if v.Hash == nil {
					return nil, nil
				}
				return v.Hash.Bytes(), nil
			},
		},
	}
}

func interfaceToStoredBatch(value interface{}) (*StoredBatch, error) {
	p, ok := value.(*StoredBatch)
	if ok {
		return p, nil
	}
	v, ok := value.(StoredBatch)
	if ok {
		return &v, nil
	}
	return nil, errors.WithStack(errInvalidStoredBatchIndexType)
}
