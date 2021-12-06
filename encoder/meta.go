package encoder

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

func DecodeMeta(meta [32]byte) models.BatchMeta {
	return models.BatchMeta{
		BatchType:  batchtype.BatchType(meta[0]),
		Size:       meta[1],
		Committer:  common.BytesToAddress(meta[2:22]),
		FinaliseOn: binary.BigEndian.Uint32(meta[22:26]),
	}
}
