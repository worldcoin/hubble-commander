package applier

import "github.com/Worldcoin/hubble-commander/models"

type ApplySingleTxResult interface {
	AppliedTx() models.GenericTransaction
	AddedPubKeyID() uint32
	PendingAccount() *models.AccountLeaf
}

type ApplySingleTransferResult struct {
	tx *models.Transfer
}

func NewApplySingleTransferResult(tx *models.Transfer) *ApplySingleTransferResult {
	return &ApplySingleTransferResult{tx: tx}
}

func (a *ApplySingleTransferResult) AppliedTx() models.GenericTransaction {
	return a.tx
}

func (a *ApplySingleTransferResult) AddedPubKeyID() uint32 {
	panic("AddedPubKeyID cannot be invoked on ApplySingleTransferResult")
}

func (a *ApplySingleTransferResult) PendingAccount() *models.AccountLeaf {
	panic("PendingAccount cannot be invoked on ApplySingleTransferResult")
}

type ApplySingleC2TResult struct {
	tx             *models.Create2Transfer
	pubKeyID       uint32
	pendingAccount *models.AccountLeaf
}

func NewApplySingleC2TResult(tx *models.Create2Transfer, pubKeyID uint32) *ApplySingleC2TResult {
	return &ApplySingleC2TResult{tx: tx, pubKeyID: pubKeyID}
}

func (a *ApplySingleC2TResult) AppliedTx() models.GenericTransaction {
	return a.tx
}

func (a *ApplySingleC2TResult) AddedPubKeyID() uint32 {
	return a.pubKeyID
}

func (a *ApplySingleC2TResult) PendingAccount() *models.AccountLeaf {
	return a.pendingAccount
}

type ApplySingleMassMigrationResult struct {
	tx *models.MassMigration
}

func (a *ApplySingleMassMigrationResult) AppliedTx() models.GenericTransaction {
	return a.tx
}

func (a *ApplySingleMassMigrationResult) AddedPubKeyID() uint32 {
	panic("AddedPubKeyID cannot be invoked on ApplySingleMassMigrationResult")
}

func (a *ApplySingleMassMigrationResult) PendingAccount() *models.AccountLeaf {
	panic("PendingAccount cannot be invoked on ApplySingleMassMigrationResult")
}
