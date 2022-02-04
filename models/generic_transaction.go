package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type GenericTransaction interface {
	Type() txtype.TransactionType
	GetBase() *TransactionBase
	GetFromStateID() uint32
	GetToStateID() *uint32
	GetAmount() Uint256
	GetFee() Uint256
	GetNonce() Uint256
	SetNonce(nonce Uint256)
	GetSignature() Signature
	Copy() GenericTransaction
	ToTransfer() *Transfer
	ToCreate2Transfer() *Create2Transfer
	ToMassMigration() *MassMigration
}

type GenericTransactionArray interface {
	Len() int
	At(index int) GenericTransaction
	Set(index int, value GenericTransaction)
	Append(elems GenericTransactionArray) GenericTransactionArray
	AppendOne(elem GenericTransaction) GenericTransactionArray
	Slice(start, end int) GenericTransactionArray
	Type() txtype.TransactionType
	ToTransferArray() TransferArray
	ToCreate2TransferArray() Create2TransferArray
	ToMassMigrationArray() MassMigrationArray
}

type TransactionWithBatchDetails struct {
	Transaction interface{}
	BatchHash   *common.Hash
	MinedTime   *Timestamp
}

func NewGenericTransactionArray(txType txtype.TransactionType, size, capacity int) GenericTransactionArray {
	switch txType {
	case txtype.Transfer:
		return make(TransferArray, size, capacity)
	case txtype.Create2Transfer:
		return make(Create2TransferArray, size, capacity)
	case txtype.MassMigration:
		return make(MassMigrationArray, size, capacity)
	}
	return nil
}

type TransferArray []Transfer

func MakeTransferArray(transfers ...Transfer) TransferArray {
	return transfers
}

func (t TransferArray) Len() int {
	return len(t)
}

func (t TransferArray) At(index int) GenericTransaction {
	return &t[index]
}

func (t TransferArray) Set(index int, value GenericTransaction) {
	t[index] = *value.ToTransfer()
}

func (t TransferArray) Append(elems GenericTransactionArray) GenericTransactionArray {
	return append(t, elems.ToTransferArray()...)
}

func (t TransferArray) AppendOne(elem GenericTransaction) GenericTransactionArray {
	return append(t, *elem.ToTransfer())
}

func (t TransferArray) Slice(start, end int) GenericTransactionArray {
	return t[start:end]
}

func (t TransferArray) Type() txtype.TransactionType {
	return txtype.Transfer
}

func (t TransferArray) ToTransferArray() TransferArray {
	return t
}

func (t TransferArray) ToCreate2TransferArray() Create2TransferArray {
	panic("TransferArray cannot be cast to Create2TransferArray")
}

func (t TransferArray) ToMassMigrationArray() MassMigrationArray {
	panic("TransferArray cannot be cast to MassMigrationArray")
}

type Create2TransferArray []Create2Transfer

func MakeCreate2TransferArray(txns ...Create2Transfer) Create2TransferArray {
	return txns
}

func (t Create2TransferArray) Len() int {
	return len(t)
}

func (t Create2TransferArray) At(index int) GenericTransaction {
	return &t[index]
}

func (t Create2TransferArray) Set(index int, value GenericTransaction) {
	t[index] = *value.ToCreate2Transfer()
}

func (t Create2TransferArray) Append(elems GenericTransactionArray) GenericTransactionArray {
	return append(t, elems.ToCreate2TransferArray()...)
}

func (t Create2TransferArray) AppendOne(elem GenericTransaction) GenericTransactionArray {
	return append(t, *elem.ToCreate2Transfer())
}

func (t Create2TransferArray) Slice(start, end int) GenericTransactionArray {
	return t[start:end]
}

func (t Create2TransferArray) Type() txtype.TransactionType {
	return txtype.Create2Transfer
}

func (t Create2TransferArray) ToTransferArray() TransferArray {
	panic("Create2TransferArray cannot be cast to TransferArray")
}

func (t Create2TransferArray) ToCreate2TransferArray() Create2TransferArray {
	return t
}

func (t Create2TransferArray) ToMassMigrationArray() MassMigrationArray {
	panic("Create2TransferArray cannot be cast to MassMigrationArray")
}

type MassMigrationArray []MassMigration

func MakeMassMigrationArray(txns ...MassMigration) MassMigrationArray {
	return txns
}

func (m MassMigrationArray) Len() int {
	return len(m)
}

func (m MassMigrationArray) At(index int) GenericTransaction {
	return &m[index]
}

func (m MassMigrationArray) Set(index int, value GenericTransaction) {
	m[index] = *value.ToMassMigration()
}

func (m MassMigrationArray) Append(elems GenericTransactionArray) GenericTransactionArray {
	return append(m, elems.ToMassMigrationArray()...)
}

func (m MassMigrationArray) AppendOne(elem GenericTransaction) GenericTransactionArray {
	return append(m, *elem.ToMassMigration())
}

func (m MassMigrationArray) Slice(start, end int) GenericTransactionArray {
	return m[start:end]
}

func (m MassMigrationArray) Type() txtype.TransactionType {
	return txtype.MassMigration
}

func (m MassMigrationArray) ToTransferArray() TransferArray {
	panic("MassMigrationArray cannot be cast to TransferArray")
}

func (m MassMigrationArray) ToCreate2TransferArray() Create2TransferArray {
	panic("MassMigrationArray cannot be cast to Create2TransferArray")
}

func (m MassMigrationArray) ToMassMigrationArray() MassMigrationArray {
	return m
}
