package executor

import "github.com/Worldcoin/hubble-commander/models"

type ApplyCommitmentResult interface {
	AppliedTransfers() models.GenericTransactionArray
	AddedPubKeyIDs() []uint32
}

func NewApplyCommitmentResult(appliedTxs models.GenericTransactionArray) ApplyCommitmentResult {
	return &AppliedTransferCommitments{
		appliedTransfers: appliedTxs.ToTransferArray(),
	}
}

type AppliedTransferCommitments struct {
	appliedTransfers models.TransferArray
}

func (a *AppliedTransferCommitments) AppliedTransfers() models.GenericTransactionArray {
	return a.appliedTransfers
}

func (a *AppliedTransferCommitments) AddedPubKeyIDs() []uint32 {
	panic("AddedPubKeyIDs cannot be invoked on AppliedTransferCommitments")
}

type ApplyTxsResult interface {
	AppliedTxs() models.GenericTransactionArray
	InvalidTxs() models.GenericTransactionArray
	AddInvalidTx(tx models.GenericTransaction)
	AddAppliedTx(tx models.GenericTransaction)
	AddedPubKeyIDs() []uint32
}

type AppliedTransfers struct {
	appliedTransfers models.GenericTransactionArray
	invalidTransfers models.GenericTransactionArray
}

func (a *AppliedTransfers) AppliedTxs() models.GenericTransactionArray {
	return a.appliedTransfers
}

func (a *AppliedTransfers) InvalidTxs() models.GenericTransactionArray {
	return a.invalidTransfers
}

func (a *AppliedTransfers) AddInvalidTx(tx models.GenericTransaction) {
	a.invalidTransfers = a.invalidTransfers.AppendOne(tx)
}

func (a *AppliedTransfers) AddAppliedTx(tx models.GenericTransaction) {
	a.appliedTransfers = a.appliedTransfers.AppendOne(tx)
}

func (a *AppliedTransfers) AddedPubKeyIDs() []uint32 {
	panic("AddedPubKeyIDs cannot be invoked on AppliedTxs")
}
