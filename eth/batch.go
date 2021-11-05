package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

//TODO-sync: rename
type DecodedBatchInt interface {
	GetBatch() *models.Batch
	ToDecodedTxBatch() *DecodedBatch
	ToDecodedDepositBatch() *DecodedDepositBatch
	SetCalldata(calldata []byte) error
}

func newDecodedBatch(batch *models.Batch) DecodedBatchInt {
	switch batch.Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return &DecodedBatch{
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

type DecodedBatch struct {
	models.Batch
	Commitments []encoder.DecodedCommitment
}

func (b *DecodedBatch) SetCalldata(calldata []byte) error {
	commitments, err := encoder.DecodeBatchCalldata(calldata, &b.ID)
	if err != nil {
		return err
	}
	b.Commitments = commitments
	return nil
}

func (b *DecodedBatch) GetBatch() *models.Batch {
	return &b.Batch
}

func (b *DecodedBatch) ToDecodedDepositBatch() *DecodedDepositBatch {
	panic("ToDecodedDepositBatch cannot be invoked on DecodedBatch")
}

func (b *DecodedBatch) ToDecodedTxBatch() *DecodedBatch {
	return b
}

type DecodedDepositBatch struct {
	models.Batch
	PathAtDepth uint32
}

func (b *DecodedDepositBatch) SetCalldata(calldata []byte) error {
	pathAtDepth, err := encoder.DecodeDepositBatchCalldata(calldata)
	if err != nil {
		return err
	}
	b.PathAtDepth = *pathAtDepth
	return nil
}

func (b *DecodedDepositBatch) GetBatch() *models.Batch {
	return &b.Batch
}

func (b *DecodedDepositBatch) ToDecodedDepositBatch() *DecodedDepositBatch {
	return b
}

func (b *DecodedDepositBatch) ToDecodedTxBatch() *DecodedBatch {
	panic("ToDecodedTxBatch cannot be invoked on DecodedDepositBatch")
}
