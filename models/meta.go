package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type Meta struct {
	BatchType  uint8
	Size       uint8
	Committer  common.Address
	FinaliseOn uint32
}
