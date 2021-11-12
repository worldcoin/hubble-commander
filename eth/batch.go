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
	GetBatch() *models.Batch
	ToDecodedTxBatch() *DecodedTxBatch
	ToDecodedDepositBatch() *DecodedDepositBatch
	SetCalldata(calldata []byte) error
	GetCommitmentsLength() int
	verifyBatchHash() error
}

func newDecodedBatch(batch *models.Batch) DecodedBatch {
	switch batch.Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return &DecodedTxBatch{
			Batch: *batch,
		}
	case batchtype.Deposit:
		return &DecodedDepositBatch{
			Batch: *batch,
		}
	case batchtype.Genesis, batchtype.MassMigration:
		panic("batch type not supported")
	}
	return nil
}

type DecodedTxBatch struct {
	models.Batch
	Commitments []encoder.DecodedCommitment
}

func (b *DecodedTxBatch) GetID() models.Uint256 {
	return b.Batch.ID
}

func (b *DecodedTxBatch) GetBatch() *models.Batch {
	return &b.Batch
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
		leafHashes = append(leafHashes, b.Commitments[i].LeafHash(*b.AccountTreeRoot))
	}
	tree, err := merkletree.NewMerkleTree(leafHashes)
	if err != nil {
		return err
	}

	if tree.Root() != *b.Hash {
		return errBatchAlreadyRolledBack
	}
	return nil
}

type DecodedDepositBatch struct {
	models.Batch // TODO create and use eth.DecodedBatchBase type here
	PathAtDepth  uint32
}

func (b *DecodedDepositBatch) GetID() models.Uint256 {
	return b.Batch.ID
}

func (b *DecodedDepositBatch) GetBatch() *models.Batch {
	return &b.Batch
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
