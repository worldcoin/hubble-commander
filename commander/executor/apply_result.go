package executor

import "github.com/Worldcoin/hubble-commander/models"

type ApplyTxsForCommitmentResult interface {
	AppliedTransfers() models.GenericTransactionArray
	AddedPubKeyIDs() []uint32
}

func NewApplyTxsForCommitmentResult(appliedTxs models.GenericTransactionArray) ApplyTxsForCommitmentResult {
	return &ApplyTransfersForCommitmentResult{
		appliedTransfers: appliedTxs.ToTransferArray(),
	}
}

type ApplyTransfersForCommitmentResult struct {
	appliedTransfers models.TransferArray
}

func (a *ApplyTransfersForCommitmentResult) AppliedTransfers() models.GenericTransactionArray {
	return a.appliedTransfers
}

func (a *ApplyTransfersForCommitmentResult) AddedPubKeyIDs() []uint32 {
	panic("AddedPubKeyIDs cannot be invoked on ApplyTransfersForCommitmentResult")
}

type ApplyTxsResult interface {
	AppliedTxs() models.GenericTransactionArray
	InvalidTxs() models.GenericTransactionArray
	AddInvalidTx(tx models.GenericTransaction)
	AddAppliedTx(tx models.GenericTransaction)
	AddedPubKeyIDs() []uint32
}

type ApplyTransfersResult struct {
	appliedTransfers models.GenericTransactionArray
	invalidTransfers models.GenericTransactionArray
}

func (a *ApplyTransfersResult) AppliedTxs() models.GenericTransactionArray {
	return a.appliedTransfers
}

func (a *ApplyTransfersResult) InvalidTxs() models.GenericTransactionArray {
	return a.invalidTransfers
}

func (a *ApplyTransfersResult) AddInvalidTx(tx models.GenericTransaction) {
	a.invalidTransfers = a.invalidTransfers.AppendOne(tx)
}

func (a *ApplyTransfersResult) AddAppliedTx(tx models.GenericTransaction) {
	a.appliedTransfers = a.appliedTransfers.AppendOne(tx)
}

func (a *ApplyTransfersResult) AddedPubKeyIDs() []uint32 {
	panic("AddedPubKeyIDs cannot be invoked on AppliedTxs")
}
