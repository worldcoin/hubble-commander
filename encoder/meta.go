package encoder

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

func DecodeMeta(meta [32]byte) models.Meta {
	return models.Meta{
		BatchType:  txtype.TransactionType(meta[0]),
		Size:       meta[1],
		Committer:  common.BytesToAddress(meta[2:22]),
		FinaliseOn: binary.BigEndian.Uint32(meta[22:26]),
	}
}
