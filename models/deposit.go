package models

import (
	"encoding/binary"
)

const (
	depositDataLength   = 76
	depositIDDataLength = 8
)

type DepositID struct {
	BlockNumber uint32
	LogIndex    uint32
}

type PendingDeposit struct {
	ID         DepositID
	ToPubKeyID uint32
	TokenID    Uint256
	L2Amount   Uint256
}

func (d *DepositID) Bytes() []byte {
	b := make([]byte, depositIDDataLength)
	binary.BigEndian.PutUint32(b[0:4], d.BlockNumber)
	binary.BigEndian.PutUint32(b[4:8], d.LogIndex)
	return b
}

func (d *DepositID) SetBytes(data []byte) error {
	if len(data) != depositIDDataLength {
		return ErrInvalidLength
	}

	d.BlockNumber = binary.BigEndian.Uint32(data[0:4])
	d.LogIndex = binary.BigEndian.Uint32(data[4:8])

	return nil
}

func (d *PendingDeposit) Bytes() []byte {
	b := make([]byte, depositDataLength)

	copy(b[0:8], d.ID.Bytes())
	binary.BigEndian.PutUint32(b[8:12], d.ToPubKeyID)
	copy(b[12:44], d.TokenID.Bytes())
	copy(b[44:76], d.L2Amount.Bytes())

	return b
}

func (d *PendingDeposit) SetBytes(data []byte) error {
	if len(data) != depositDataLength {
		return ErrInvalidLength
	}

	err := d.ID.SetBytes(data[0:8])
	if err != nil {
		return err
	}

	d.ToPubKeyID = binary.BigEndian.Uint32(data[8:12])
	d.TokenID.SetBytes(data[12:44])
	d.L2Amount.SetBytes(data[44:76])

	return nil
}
