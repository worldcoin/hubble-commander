package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type BatchMeta struct {
	BatchType  batchtype.BatchType
	Size       uint8
	Committer  common.Address
	FinaliseOn uint32
}

type MassMigrationMeta struct {
	SpokeID     uint32
	TokenID     Uint256
	Amount      Uint256
	FeeReceiver uint32
}
