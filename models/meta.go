package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type Meta struct {
	BatchType  batchtype.BatchType
	Size       uint8
	Committer  common.Address
	FinaliseOn uint32
}
