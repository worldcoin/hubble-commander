package executor

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
)

type ExecuteTxsForCommitmentResult interface {
	AppliedTxs() models.GenericTransactionArray
	AddedPubKeyIDs() []uint32
	PendingAccounts() []models.AccountLeaf
}

type ExecuteTransfersForCommitmentResult struct {
	appliedTxs models.TransferArray
}

func (a *ExecuteTransfersForCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return a.appliedTxs
}

func (a *ExecuteTransfersForCommitmentResult) AddedPubKeyIDs() []uint32 {
	panic("AddedPubKeyIDs cannot be invoked on ExecuteTransfersForCommitmentResult")
}

func (a *ExecuteTransfersForCommitmentResult) PendingAccounts() []models.AccountLeaf {
	panic("PendingAccounts cannot be invoked on ExecuteTransfersForCommitmentResult")
}

type ExecuteC2TForCommitmentResult struct {
	appliedTxs      models.Create2TransferArray
	addedPubKeyIDs  []uint32
	pendingAccounts []models.AccountLeaf
}

func (a *ExecuteC2TForCommitmentResult) AppliedTxs() models.GenericTransactionArray {
	return a.appliedTxs
}

func (a *ExecuteC2TForCommitmentResult) AddedPubKeyIDs() []uint32 {
	return a.addedPubKeyIDs
}

func (a *ExecuteC2TForCommitmentResult) PendingAccounts() []models.AccountLeaf {
	return a.pendingAccounts
}

type ExecuteTxsResult interface {
	AppliedTxs() models.GenericTransactionArray
	InvalidTxs() models.GenericTransactionArray
	AddedPubKeyIDs() []uint32
	PendingAccounts() []models.AccountLeaf
	AllTxs() models.GenericTransactionArray
	AddApplied(singleTxResult applier.ApplySingleTxResult)
	AddInvalidTx(tx models.GenericTransaction)
	AddApplyResult(other ExecuteTxsResult)
}

type ExecuteTransfersResult struct {
	appliedTxs models.GenericTransactionArray
	invalidTxs models.GenericTransactionArray
}

func (a *ExecuteTransfersResult) AppliedTxs() models.GenericTransactionArray {
	return a.appliedTxs
}

func (a *ExecuteTransfersResult) InvalidTxs() models.GenericTransactionArray {
	return a.invalidTxs
}

func (a *ExecuteTransfersResult) AddedPubKeyIDs() []uint32 {
	panic("AddedPubKeyIDs cannot be invoked on ExecuteTransfersResult")
}

func (a *ExecuteTransfersResult) PendingAccounts() []models.AccountLeaf {
	panic("PendingAccounts cannot be invoked on ExecuteTransfersResult")
}

func (a *ExecuteTransfersResult) AllTxs() models.GenericTransactionArray {
	return a.appliedTxs.Append(a.invalidTxs)
}

func (a *ExecuteTransfersResult) AddApplied(singleTxResult applier.ApplySingleTxResult) {
	a.appliedTxs = a.appliedTxs.AppendOne(singleTxResult.AppliedTx())
}

func (a *ExecuteTransfersResult) AddInvalidTx(tx models.GenericTransaction) {
	a.invalidTxs = a.invalidTxs.AppendOne(tx)
}

func (a *ExecuteTransfersResult) AddApplyResult(other ExecuteTxsResult) {
	a.appliedTxs = a.appliedTxs.Append(other.AppliedTxs())
	a.invalidTxs = a.invalidTxs.Append(other.InvalidTxs())
}

type ExecuteC2TResult struct {
	appliedTxs      models.GenericTransactionArray
	invalidTxs      models.GenericTransactionArray
	addedPubKeyIDs  []uint32
	pendingAccounts []models.AccountLeaf
}

func (a *ExecuteC2TResult) AppliedTxs() models.GenericTransactionArray {
	return a.appliedTxs
}

func (a *ExecuteC2TResult) InvalidTxs() models.GenericTransactionArray {
	return a.invalidTxs
}

func (a *ExecuteC2TResult) AddedPubKeyIDs() []uint32 {
	return a.addedPubKeyIDs
}

func (a *ExecuteC2TResult) PendingAccounts() []models.AccountLeaf {
	return a.pendingAccounts
}

func (a *ExecuteC2TResult) AllTxs() models.GenericTransactionArray {
	return a.appliedTxs.Append(a.invalidTxs)
}

func (a *ExecuteC2TResult) AddApplied(singleTxResult applier.ApplySingleTxResult) {
	a.appliedTxs = a.appliedTxs.AppendOne(singleTxResult.AppliedTx())
	a.addedPubKeyIDs = append(a.addedPubKeyIDs, singleTxResult.AddedPubKeyID())

	if singleTxResult.PendingAccount() != nil {
		a.pendingAccounts = append(a.pendingAccounts, *singleTxResult.PendingAccount())
	}
}

func (a *ExecuteC2TResult) AddInvalidTx(tx models.GenericTransaction) {
	a.invalidTxs = a.invalidTxs.AppendOne(tx)
}

func (a *ExecuteC2TResult) AddApplyResult(other ExecuteTxsResult) {
	a.appliedTxs = a.appliedTxs.Append(other.AppliedTxs())
	a.invalidTxs = a.invalidTxs.Append(other.InvalidTxs())
	a.addedPubKeyIDs = append(a.addedPubKeyIDs, other.AddedPubKeyIDs()...)
	a.pendingAccounts = append(a.pendingAccounts, other.PendingAccounts()...)
}
