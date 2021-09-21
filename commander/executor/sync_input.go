package executor

import "github.com/Worldcoin/hubble-commander/models"

type SyncedTxs interface {
	Txs() models.GenericTransactionArray
	PubKeyIDs() []uint32
}

type SyncedTransfers struct {
	txs models.TransferArray
}

func (a *SyncedTransfers) Txs() models.GenericTransactionArray {
	return a.txs
}

func (a *SyncedTransfers) PubKeyIDs() []uint32 {
	panic("PubKeyIDs cannot be invoked on ApplyTransfersForCommitmentResult")
}

type SyncedC2Ts struct {
	txs       models.Create2TransferArray
	pubKeyIDs []uint32
}

func (a *SyncedC2Ts) Txs() models.GenericTransactionArray {
	return a.txs
}

func (a *SyncedC2Ts) PubKeyIDs() []uint32 {
	return a.pubKeyIDs
}
