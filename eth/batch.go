package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type DecodedBatch interface {
	GetID() models.Uint256
	GetBase() *DecodedBatchBase
	ToDecodedTxBatch() *DecodedTxBatch
	ToDecodedDepositBatch() *DecodedDepositBatch
	ToBatch(prevStateRoot common.Hash) *models.Batch
	GetCommitmentsLength() int
}

type DecodedBatchBase struct {
	ID                models.Uint256
	Type              batchtype.BatchType
	TransactionHash   common.Hash
	Hash              common.Hash
	FinalisationBlock uint32
	AccountTreeRoot   common.Hash
	SubmissionTime    models.Timestamp
	Committer         common.Address
}

func NewDecodedBatchBase(
	batch *ContractBatch,
	transactionHash, accountRoot common.Hash,
	submissionTime *models.Timestamp,
) *DecodedBatchBase {
	return &DecodedBatchBase{
		ID:                batch.ID,
		Type:              batch.BatchType,
		TransactionHash:   transactionHash,
		Hash:              batch.Hash,
		FinalisationBlock: batch.FinaliseOn,
		AccountTreeRoot:   accountRoot,
		SubmissionTime:    *submissionTime,
		Committer:         batch.Committer,
	}
}

func (b *DecodedBatchBase) ToBatch(prevStateRoot common.Hash) *models.Batch {
	return &models.Batch{
		ID:                b.ID,
		Type:              b.Type,
		TransactionHash:   b.TransactionHash,
		Hash:              &b.Hash,
		FinalisationBlock: &b.FinalisationBlock,
		AccountTreeRoot:   &b.AccountTreeRoot,
		SubmissionTime:    &b.SubmissionTime,
		PrevStateRoot:     &prevStateRoot,
	}
}

type DecodedTxBatch struct {
	DecodedBatchBase
	Commitments []encoder.Commitment
}

func (b *DecodedTxBatch) GetID() models.Uint256 {
	return b.DecodedBatchBase.ID
}

func (b *DecodedTxBatch) GetBase() *DecodedBatchBase {
	return &b.DecodedBatchBase
}

func (b *DecodedTxBatch) ToDecodedTxBatch() *DecodedTxBatch {
	return b
}

func (b *DecodedTxBatch) ToDecodedDepositBatch() *DecodedDepositBatch {
	panic("ToDecodedDepositBatch cannot be invoked on DecodedTxBatch")
}

func (b *DecodedTxBatch) GetCommitmentsLength() int {
	return len(b.Commitments)
}

type DecodedDepositBatch struct {
	DecodedBatchBase
	SubtreeID   models.Uint256
	PathAtDepth uint32
}

func (b *DecodedDepositBatch) GetID() models.Uint256 {
	return b.DecodedBatchBase.ID
}

func (b *DecodedDepositBatch) GetBase() *DecodedBatchBase {
	return &b.DecodedBatchBase
}

func (b *DecodedDepositBatch) ToDecodedTxBatch() *DecodedTxBatch {
	panic("ToDecodedTxBatch cannot be invoked on DecodedDepositBatch")
}

func (b *DecodedDepositBatch) ToDecodedDepositBatch() *DecodedDepositBatch {
	return b
}

func (b *DecodedDepositBatch) GetCommitmentsLength() int {
	return 1
}
