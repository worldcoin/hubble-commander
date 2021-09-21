package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type SyncedTxs interface {
	Txs() models.GenericTransactionArray
	PubKeyIDs() []uint32
	SyncedTxAt(index int) SyncedTx
	SetTxs(txs models.GenericTransactionArray)
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

func (a *SyncedTransfers) SyncedTxAt(index int) SyncedTx {
	return &SyncedTransfer2{
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

func (a *SyncedC2Ts) SyncedTxAt(index int) SyncedTx {
	return &SyncedC2T2{
		tx:       a.txs.At(index).ToCreate2Transfer(),
		pubKeyID: a.pubKeyIDs[index],
	}
}

func (a *SyncedC2Ts) SetTxs(txs models.GenericTransactionArray) {
	a.txs = txs.ToCreate2TransferArray()
}

type SyncedTx interface {
	SyncedTx() models.GenericTransaction
	SyncedPubKeyID() uint32
}

type SyncedTransfer2 struct {
	tx *models.Transfer
}

func (a *SyncedTransfer2) SyncedTx() models.GenericTransaction {
	return a.tx
}

func (a *SyncedTransfer2) SyncedPubKeyID() uint32 {
	panic("SyncedPubKeyID cannot be invoked on SyncedTransfer2")
}

type SyncedC2T2 struct {
	tx       *models.Create2Transfer
	pubKeyID uint32
}

func (a *SyncedC2T2) SyncedTx() models.GenericTransaction {
	return a.tx
}

func (a *SyncedC2T2) SyncedPubKeyID() uint32 {
	return a.pubKeyID
}
