package applier

import "github.com/Worldcoin/hubble-commander/models"

type SingleTxResult interface {
	AppliedTx() models.GenericTransaction
	AddedPubKeyID() uint32
}

type ApplySingleTransferResult struct {
	tx models.GenericTransaction
}

func (a *ApplySingleTransferResult) AppliedTx() models.GenericTransaction {
	return a.tx
}

func (a *ApplySingleTransferResult) AddedPubKeyID() uint32 {
	panic("AddedPubKeyID cannot be invoked on ApplySingleTransferResult")
}

type ApplySingleC2TResult struct {
	tx            *models.Create2Transfer
	addedPubKeyID uint32
}

func (a *ApplySingleC2TResult) AppliedTx() models.GenericTransaction {
	return a.tx
}

func (a *ApplySingleC2TResult) AddedPubKeyID() uint32 {
	return a.addedPubKeyID
}
