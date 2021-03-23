package models

type TransactionReceipt struct {
	Transaction
	Status      TransactionStatus
}
