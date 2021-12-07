package stored

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

const batchDataLength = 185

var (
	BatchName                = models.GetTypeName(Batch{})
	BatchPrefix              = models.GetBadgerHoldPrefix(Batch{})
	errInvalidBatchIndexType = fmt.Errorf("invalid stored.Batch index type")
)

type Batch struct {
	ID                models.Uint256
	BType             batchtype.BatchType
	TransactionHash   common.Hash
	Hash              *common.Hash // root of tree containing all commitments included in this batch
	FinalisationBlock *uint32
	AccountTreeRoot   *common.Hash
	PrevStateRoot     *common.Hash
	SubmissionTime    *models.Timestamp
}

func NewBatchFromModelsBatch(b *models.Batch) *Batch {
	return &Batch{
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

func (b *Batch) ToModelsBatch() *models.Batch {
	return &models.Batch{
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

func (b *Batch) Bytes() []byte {
	encoded := make([]byte, batchDataLength)
	copy(encoded[0:32], b.ID.Bytes())
	encoded[32] = byte(b.BType)
	copy(encoded[33:65], b.TransactionHash.Bytes())
	copy(encoded[65:98], models.EncodeHashPointer(b.Hash))
	copy(encoded[98:103], models.EncodeUint32Pointer(b.FinalisationBlock))
	copy(encoded[103:136], models.EncodeHashPointer(b.AccountTreeRoot))
	copy(encoded[136:169], models.EncodeHashPointer(b.PrevStateRoot))
	copy(encoded[169:185], models.EncodeTimestampPointer(b.SubmissionTime))

	return encoded
}

func (b *Batch) SetBytes(data []byte) error {
	if len(data) != batchDataLength {
		return models.ErrInvalidLength
	}
	timestamp, err := models.DecodeTimestampPointer(data[169:185])
	if err != nil {
		return err
	}

	b.ID.SetBytes(data[0:32])
	b.BType = batchtype.BatchType(data[32])
	b.TransactionHash.SetBytes(data[33:65])
	b.Hash = models.DecodeHashPointer(data[65:98])
	b.FinalisationBlock = models.DecodeUint32Pointer(data[98:103])
	b.AccountTreeRoot = models.DecodeHashPointer(data[103:136])
	b.PrevStateRoot = models.DecodeHashPointer(data[136:169])
	b.SubmissionTime = timestamp
	return nil
}

// nolint:gocritic
// Type implements badgerhold.Storer
func (b Batch) Type() string {
	return string(BatchName)
}

// nolint:gocritic
// Indexes implements badgerhold.Storer
func (b Batch) Indexes() map[string]bh.Index {
	return map[string]bh.Index{
		"Hash": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToBatch(value)
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

func interfaceToBatch(value interface{}) (*Batch, error) {
	p, ok := value.(*Batch)
	if ok {
		return p, nil
	}
	v, ok := value.(Batch)
	if ok {
		return &v, nil
	}
	return nil, errors.WithStack(errInvalidBatchIndexType)
}
