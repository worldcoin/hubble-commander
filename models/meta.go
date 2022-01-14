package models

import (
	"encoding/binary"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
)

const mmMetaLength = 72 // 4 + 32 + 32 + 4

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

func (m *MassMigrationMeta) Bytes() []byte {
	b := make([]byte, mmMetaLength)

	binary.BigEndian.PutUint32(b[0:4], m.SpokeID)
	copy(b[4:36], m.TokenID.Bytes())
	copy(b[36:68], m.Amount.Bytes())
	binary.BigEndian.PutUint32(b[68:72], m.FeeReceiver)

	return b
}

func (m *MassMigrationMeta) SetBytes(data []byte) error {
	if len(data) != mmMetaLength {
		return ErrInvalidLength
	}

	m.SpokeID = binary.BigEndian.Uint32(data[0:4])
	m.TokenID.SetBytes(data[4:36])
	m.Amount.SetBytes(data[36:68])
	m.FeeReceiver = binary.BigEndian.Uint32(data[68:72])

	return nil
}
