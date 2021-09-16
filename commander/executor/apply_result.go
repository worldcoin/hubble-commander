package executor

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
)

type ApplyTxsForCommitmentResult interface {
	AppliedTxs() models.GenericTransactionArray
	AddedPubKeyIDs() []uint32
}

type ApplyTransfersForCommitmentResult struct {
	appliedTxs models.TransferArray
}

func (a *ApplyTransfersForCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return a.appliedTxs
}

func (a *ApplyTransfersForCommitmentResult) AddedPubKeyIDs() []uint32 {
	panic("AddedPubKeyIDs cannot be invoked on ApplyTransfersForCommitmentResult")
}

type ApplyC2TForCommitmentResult struct {
	appliedTxs     models.Create2TransferArray
	addedPubKeyIDs []uint32
}

func (a *ApplyC2TForCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return a.appliedTxs
}

func (a *ApplyC2TForCommitmentResult) AddedPubKeyIDs() []uint32 {
	return a.addedPubKeyIDs
}

type ApplyTxsResult interface {
	AppliedTxs() models.GenericTransactionArray
	InvalidTxs() models.GenericTransactionArray
	AddedPubKeyIDs() []uint32
	AllTxs() models.GenericTransactionArray
	AddApplied(singleTxResult applier.SingleTxResult)
	AddInvalidTx(tx models.GenericTransaction)
	AddTxs(other ApplyTxsResult)
}

type ApplyTransfersResult struct {
	appliedTxs models.GenericTransactionArray
	invalidTxs models.GenericTransactionArray
}

func (a *ApplyTransfersResult) AppliedTxs() models.GenericTransactionArray {
	return a.appliedTxs
}

func (a *ApplyTransfersResult) InvalidTxs() models.GenericTransactionArray {
	return a.invalidTxs
}

func (a *ApplyTransfersResult) AddedPubKeyIDs() []uint32 {
	panic("AddedPubKeyIDs cannot be invoked on AppliedTxs")
}

func (a *ApplyTransfersResult) AllTxs() models.GenericTransactionArray {
	return a.appliedTxs.Append(a.invalidTxs)
}

func (a *ApplyTransfersResult) AddApplied(singleTxResult applier.SingleTxResult) {
	a.appliedTxs = a.appliedTxs.AppendOne(singleTxResult.AppliedTx())
}

func (a *ApplyTransfersResult) AddInvalidTx(tx models.GenericTransaction) {
	a.invalidTxs = a.invalidTxs.AppendOne(tx)
}

func (a *ApplyTransfersResult) AddTxs(other ApplyTxsResult) {
	a.appliedTxs = a.appliedTxs.Append(other.AppliedTxs())
	a.invalidTxs = a.invalidTxs.Append(other.InvalidTxs())
}

type ApplyC2TResult struct {
	appliedTxs     models.GenericTransactionArray
	invalidTxs     models.GenericTransactionArray
	addedPubKeyIDs []uint32
}

func (a *ApplyC2TResult) AppliedTxs() models.GenericTransactionArray {
	return a.appliedTxs
}

func (a *ApplyC2TResult) InvalidTxs() models.GenericTransactionArray {
	return a.invalidTxs
}

func (a *ApplyC2TResult) AddedPubKeyIDs() []uint32 {
	return a.addedPubKeyIDs
}

func (a *ApplyC2TResult) AllTxs() models.GenericTransactionArray {
	return a.appliedTxs.Append(a.invalidTxs)
}

func (a *ApplyC2TResult) AddApplied(singleTxResult applier.SingleTxResult) {
	a.appliedTxs = a.appliedTxs.AppendOne(singleTxResult.AppliedTx())
	a.addedPubKeyIDs = append(a.addedPubKeyIDs, singleTxResult.AddedPubKeyID())
}

func (a *ApplyC2TResult) AddInvalidTx(tx models.GenericTransaction) {
	a.invalidTxs = a.invalidTxs.AppendOne(tx)
}

func (a *ApplyC2TResult) AddTxs(other ApplyTxsResult) {
	a.appliedTxs = a.appliedTxs.Append(other.AppliedTxs())
	a.invalidTxs = a.invalidTxs.Append(other.InvalidTxs())
	a.addedPubKeyIDs = append(a.addedPubKeyIDs, other.AddedPubKeyIDs()...)
}
