package executor

import "github.com/Worldcoin/hubble-commander/models"

type ApplyTxsForCommitmentResult interface {
	AppliedTransfers() models.GenericTransactionArray
	AddedPubKeyIDs() []uint32
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
	AddedPubKeyIDs() []uint32
	AllTxs() models.GenericTransactionArray
	AddAppliedTx(tx models.GenericTransaction)
	AddInvalidTx(tx models.GenericTransaction)
	AddTxs(other ApplyTxsResult)
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

func (a *ApplyTransfersResult) AddedPubKeyIDs() []uint32 {
	panic("AddedPubKeyIDs cannot be invoked on AppliedTxs")
}

func (a *ApplyTransfersResult) AllTxs() models.GenericTransactionArray {
	return a.appliedTransfers.Append(a.invalidTransfers)
}

func (a *ApplyTransfersResult) AddAppliedTx(tx models.GenericTransaction) {
	a.appliedTransfers = a.appliedTransfers.AppendOne(tx)
}

func (a *ApplyTransfersResult) AddInvalidTx(tx models.GenericTransaction) {
	a.invalidTransfers = a.invalidTransfers.AppendOne(tx)
}

func (a *ApplyTransfersResult) AddTxs(other ApplyTxsResult) {
	a.appliedTransfers = a.appliedTransfers.Append(other.AppliedTxs())
	a.invalidTransfers = a.invalidTransfers.Append(other.InvalidTxs())
}
