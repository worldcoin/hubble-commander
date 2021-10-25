package models

import "github.com/ethereum/go-ethereum/common"

type TransactionError struct {
	Hash         common.Hash
	ErrorMessage string
}
