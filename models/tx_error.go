package models

import "github.com/ethereum/go-ethereum/common"

type TxError struct {
	Hash         common.Hash
	ErrorMessage string
}
