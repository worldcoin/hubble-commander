package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
)

type DecodedBatch interface {
	GetID() models.Uint256
	GetBatch() *DecodedBatchBase
	ToDecodedTxBatch() *DecodedTxBatch
	ToDecodedDepositBatch() *DecodedDepositBatch
	SetCalldata(calldata []byte) error
	GetCommitmentsLength() int
	verifyBatchHash() error
}

func newDecodedBatch(batch *models.Batch, transactionHash, accountRoot common.Hash) DecodedBatch {
	switch batch.Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return &DecodedTxBatch{
			DecodedBatchBase: *NewDecodedBatchBase(batch, transactionHash, accountRoot),
		}
	case batchtype.Deposit:
		return &DecodedDepositBatch{
			DecodedBatchBase: *NewDecodedBatchBase(batch, transactionHash, accountRoot),
		}
	case batchtype.Genesis, batchtype.MassMigration:
		panic("batch type not supported")
	}
	return nil
}

type DecodedBatchBase struct {
	ID                models.Uint256
	Type              batchtype.BatchType
	TransactionHash   common.Hash
	Hash              common.Hash
	FinalisationBlock uint32
	AccountTreeRoot   common.Hash
	SubmissionTime    models.Timestamp
}

func NewDecodedBatchBase(batch *models.Batch, transactionHash, accountRoot common.Hash) *DecodedBatchBase {
	return &DecodedBatchBase{
		ID:                batch.ID,
		Type:              batch.Type,
		TransactionHash:   transactionHash,
		Hash:              *batch.Hash,
		FinalisationBlock: *batch.FinalisationBlock,
		AccountTreeRoot:   accountRoot,
		SubmissionTime:    models.Timestamp{},
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
	Commitments []encoder.DecodedCommitment
}

func (b *DecodedTxBatch) GetID() models.Uint256 {
	return b.DecodedBatchBase.ID
}

func (b *DecodedTxBatch) GetBatch() *DecodedBatchBase {
	return &b.DecodedBatchBase
}

func (b *DecodedTxBatch) ToDecodedDepositBatch() *DecodedDepositBatch {
	panic("ToDecodedDepositBatch cannot be invoked on DecodedTxBatch")
}

func (b *DecodedTxBatch) ToDecodedTxBatch() *DecodedTxBatch {
	return b
}

func (b *DecodedTxBatch) SetCalldata(calldata []byte) error {
	commitments, err := encoder.DecodeBatchCalldata(calldata, &b.ID)
	if err != nil {
		return err
	}
	b.Commitments = commitments
	return nil
}

func (b *DecodedTxBatch) GetCommitmentsLength() int {
	return len(b.Commitments)
}

func (b *DecodedTxBatch) verifyBatchHash() error {
	leafHashes := make([]common.Hash, 0, len(b.Commitments))
	for i := range b.Commitments {
		leafHashes = append(leafHashes, b.Commitments[i].LeafHash(b.AccountTreeRoot))
	}
	tree, err := merkletree.NewMerkleTree(leafHashes)
	if err != nil {
		return err
	}

	if tree.Root() != b.Hash {
		return errBatchAlreadyRolledBack
	}
	return nil
}

type DecodedDepositBatch struct {
	DecodedBatchBase
	PathAtDepth uint32
}

func (b *DecodedDepositBatch) GetID() models.Uint256 {
	return b.DecodedBatchBase.ID
}

func (b *DecodedDepositBatch) GetBatch() *DecodedBatchBase {
	return &b.DecodedBatchBase
}

func (b *DecodedDepositBatch) ToDecodedDepositBatch() *DecodedDepositBatch {
	return b
}

func (b *DecodedDepositBatch) ToDecodedTxBatch() *DecodedTxBatch {
	panic("ToDecodedTxBatch cannot be invoked on DecodedDepositBatch")
}

func (b *DecodedDepositBatch) SetCalldata(calldata []byte) error {
	pathAtDepth, err := encoder.DecodeDepositBatchCalldata(calldata)
	if err != nil {
		return err
	}
	b.PathAtDepth = *pathAtDepth
	return nil
}

func (b *DecodedDepositBatch) GetCommitmentsLength() int {
	return 1
}

func (b *DecodedDepositBatch) verifyBatchHash() error {
	return nil
}
