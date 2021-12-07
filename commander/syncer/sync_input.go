package syncer

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

func NewSyncedTransfers(txs models.TransferArray) *SyncedTransfers {
	return &SyncedTransfers{txs: txs}
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

func NewSyncedC2Ts(txs models.Create2TransferArray, pubKeyIDs []uint32) *SyncedC2Ts {
	return &SyncedC2Ts{txs: txs, pubKeyIDs: pubKeyIDs}
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

type SyncedMMs struct {
	txs models.MassMigrationArray
}

func NewSyncedMM(txs models.MassMigrationArray) *SyncedMMs {
	return &SyncedMMs{txs: txs}
}

func (a *SyncedMMs) Txs() models.GenericTransactionArray {
	return a.txs
}

func (a *SyncedMMs) PubKeyIDs() []uint32 {
	panic("PubKeyIDs cannot be invoked on SyncedMMs")
}

func (a *SyncedMMs) TxAt(index int) SyncedTx {
	return &SyncedMM{
		tx: a.txs.At(index).ToMassMigration(),
	}
}

func (a *SyncedMMs) SetTxs(txs models.GenericTransactionArray) {
	a.txs = txs.ToMassMigrationArray()
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

type SyncedMM struct {
	tx *models.MassMigration
}

func (a *SyncedMM) Tx() models.GenericTransaction {
	return a.tx
}

func (a *SyncedMM) PubKeyID() uint32 {
	panic("PubKeyID cannot be invoked on SyncedMM")
}
