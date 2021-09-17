package models

import "github.com/Worldcoin/hubble-commander/models/enums/txtype"

type GenericTransaction interface {
	Type() txtype.TransactionType
	GetBase() *TransactionBase
	GetFromStateID() uint32
	GetToStateID() *uint32
	GetAmount() Uint256
	GetFee() *Uint256
	GetNonce() Uint256
	SetNonce(nonce Uint256)
	GetSignature() Signature
	Copy() GenericTransaction
	ToTransfer() *Transfer
	ToCreate2Transfer() *Create2Transfer
}

type GenericTransactionArray interface {
	Len() int
	At(index int) GenericTransaction
	Set(index int, value GenericTransaction)
	Append(elems GenericTransactionArray) GenericTransactionArray
	AppendOne(elem GenericTransaction) GenericTransactionArray
	Slice(start, end int) GenericTransactionArray
	ToTransferArray() TransferArray
	ToCreate2TransferArray() Create2TransferArray
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

func (t TransferArray) ToTransferArray() TransferArray {
	return t
}

func (t TransferArray) ToCreate2TransferArray() Create2TransferArray {
	panic("TransferArray cannot be cast to Create2TransferArray")
}

type Create2TransferArray []Create2Transfer

func MakeCreate2TransferArray(create2Transfers ...Create2Transfer) Create2TransferArray {
	return create2Transfers
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

func (t Create2TransferArray) ToTransferArray() TransferArray {
	panic("Create2TransferArray cannot be cast to TransferArray")
}

func (t Create2TransferArray) ToCreate2TransferArray() Create2TransferArray {
	return t
}
