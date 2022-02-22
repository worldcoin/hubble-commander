package models

import "github.com/ethereum/go-ethereum/common"

type TxError struct {
	TxHash        common.Hash
	SenderStateID uint32
	ErrorMessage  string
}
