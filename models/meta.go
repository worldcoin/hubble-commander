package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type Meta struct {
	BatchType  txtype.TransactionType
	Size       uint8
	Committer  common.Address
	FinaliseOn uint32
}
