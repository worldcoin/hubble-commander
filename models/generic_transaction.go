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
	GetSignature() Signature
	Copy() GenericTransaction
}

type GenericTransactionArray interface {
	Len() int
	At(index int) GenericTransaction
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
