package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type SyncedTxs interface {
	Txs() models.GenericTransactionArray
	PubKeyIDs() []uint32
	TxAt(index int) SyncedTx
	SetTxs(txs models.GenericTransactionArray)
}

type SyncedTransfers struct {
	txs models.TransferArray
}

func (a *SyncedTransfers) Txs() models.GenericTransactionArray {
	return a.txs
}

func (a *SyncedTransfers) PubKeyIDs() []uint32 {
	panic("PubKeyIDs cannot be invoked on SyncedTransfers")
}

func (a *SyncedTransfers) TxAt(index int) SyncedTx {
	return &SyncedTransfer{
		tx: a.txs.At(index).ToTransfer(),
	}
}

func (a *SyncedTransfers) SetTxs(txs models.GenericTransactionArray) {
	a.txs = txs.ToTransferArray()
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

func (a *SyncedC2Ts) TxAt(index int) SyncedTx {
	return &SyncedC2T{
		tx:       a.txs.At(index).ToCreate2Transfer(),
		pubKeyID: a.pubKeyIDs[index],
	}
}

func (a *SyncedC2Ts) SetTxs(txs models.GenericTransactionArray) {
	a.txs = txs.ToCreate2TransferArray()
}

type SyncedTx interface {
	Tx() models.GenericTransaction
	PubKeyID() uint32
}

type SyncedTransfer struct {
	tx *models.Transfer
}

func (a *SyncedTransfer) Tx() models.GenericTransaction {
	return a.tx
}

func (a *SyncedTransfer) PubKeyID() uint32 {
	panic("PubKeyID cannot be invoked on SyncedTransfer")
}

type SyncedC2T struct {
	tx       *models.Create2Transfer
	pubKeyID uint32
}

func (a *SyncedC2T) Tx() models.GenericTransaction {
	return a.tx
}

func (a *SyncedC2T) PubKeyID() uint32 {
	return a.pubKeyID
}
