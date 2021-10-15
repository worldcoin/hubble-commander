package models

import (
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

type CommitmentBase struct {
	ID            CommitmentID
	Type          batchtype.BatchType
	PostStateRoot common.Hash
}
