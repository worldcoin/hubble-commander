package eth

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

type DecodedBatch interface {
	GetBatch() *models.Batch
	ToDecodedTxBatch() *DecodedTxBatch
	ToDecodedDepositBatch() *DecodedDepositBatch
	SetCalldata(calldata []byte) error
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

func (b *DecodedTxBatch) SetCalldata(calldata []byte) error {
	commitments, err := encoder.DecodeBatchCalldata(calldata, &b.ID)
	if err != nil {
		return err
	}
	b.Commitments = commitments
	return nil
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

func (b *DecodedDepositBatch) ToDecodedTxBatch() *DecodedTxBatch {
	panic("ToDecodedTxBatch cannot be invoked on DecodedDepositBatch")
}
