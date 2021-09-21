package executor

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
)

type SyncedTxs interface {
	Txs() models.GenericTransactionArray
	PubKeyIDs() []uint32
	SyncedTxAt(index int) applier.SyncedTx
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

func (a *SyncedTransfers) SyncedTxAt(index int) applier.SyncedTx {
	return &applier.SyncedTransfer2{
		Tx: a.txs.At(index).ToTransfer(),
	}
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

func (a *SyncedC2Ts) SyncedTxAt(index int) applier.SyncedTx {
	return &applier.SyncedC2T2{
		Tx:       a.txs.At(index).ToCreate2Transfer(),
		PubKeyID: a.pubKeyIDs[index],
	}
}
