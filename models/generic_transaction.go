package models

import "github.com/Worldcoin/hubble-commander/models/enums/txtype"

type GenericTransaction interface {
	Type() txtype.TransactionType
	GetBase() *TransactionBase
	GetFromStateID() uint32
	GetToStateID() *uint32
	GetAmount() Uint256
	GetFee() Uint256
	GetNonce() Uint256
	SetNonce(nonce Uint256)
}

type GenericTransactionArray []GenericTransaction

func MakeGenericTransferArray(transfers []Transfer) GenericTransactionArray {
	arr := make([]GenericTransaction, len(transfers))
	for i := range transfers {
		arr[i] = &transfers[i]
	}
	return arr
}

func MakeGenericCreate2TransferArray(transfers []Create2Transfer) GenericTransactionArray {
	arr := make([]GenericTransaction, len(transfers))
	for i := range transfers {
		arr[i] = &transfers[i]
	}
	return arr
}

func (a GenericTransactionArray) Type() txtype.TransactionType {
	return a[0].Type()
}

func (a GenericTransactionArray) ToTransferArray() []Transfer {
	arr := make([]Transfer, len(a))
	for i := range a {
		arr[i] = *a[i].(*Transfer)
	}
	return arr
}

func (a GenericTransactionArray) ToCreate2TransferArray() []Create2Transfer {
	arr := make([]Create2Transfer, len(a))
	for i := range a {
		arr[i] = *a[i].(*Create2Transfer)
	}
	return arr
}
