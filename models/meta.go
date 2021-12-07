package models

import (
	"math/big"

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

func NewMassMigrationMetaFromBigInts(meta [4]*big.Int) *MassMigrationMeta {
	return &MassMigrationMeta{
		SpokeID:     uint32(meta[0].Uint64()),
		TokenID:     MakeUint256FromBig(*meta[1]),
		Amount:      MakeUint256FromBig(*meta[2]),
		FeeReceiver: uint32(meta[3].Uint64()),
	}
}

func (m *MassMigrationMeta) BigInts() [4]*big.Int {
	return [4]*big.Int{
		big.NewInt(int64(m.SpokeID)),
		m.TokenID.ToBig(),
		m.Amount.ToBig(),
		big.NewInt(int64(m.FeeReceiver)),
	}
}
