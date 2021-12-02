package models

import (
	"encoding/binary"
)

const (
	depositDataLength   = 132
	depositIDDataLength = 64
)

var PendingDepositPrefix = getBadgerHoldPrefix(PendingDeposit{})

type DepositID struct {
	SubtreeID    Uint256
	DepositIndex Uint256
}

type PendingDeposit struct {
	ID         DepositID
	ToPubKeyID uint32
	TokenID    Uint256
	L2Amount   Uint256
}

func (d *DepositID) Bytes() []byte {
	b := make([]byte, depositIDDataLength)
	copy(b[0:32], d.SubtreeID.Bytes())
	copy(b[32:64], d.DepositIndex.Bytes())
	return b
}

func (d *DepositID) SetBytes(data []byte) error {
	if len(data) != depositIDDataLength {
		return ErrInvalidLength
	}

	d.SubtreeID.SetBytes(data[0:32])
	d.DepositIndex.SetBytes(data[32:64])

	return nil
}

func (d *PendingDeposit) Bytes() []byte {
	b := make([]byte, depositDataLength)

	copy(b[0:64], d.ID.Bytes())
	binary.BigEndian.PutUint32(b[64:68], d.ToPubKeyID)
	copy(b[68:100], d.TokenID.Bytes())
	copy(b[100:132], d.L2Amount.Bytes())

	return b
}

func (d *PendingDeposit) SetBytes(data []byte) error {
	if len(data) != depositDataLength {
		return ErrInvalidLength
	}

	err := d.ID.SetBytes(data[0:64])
	if err != nil {
		return err
	}

	d.ToPubKeyID = binary.BigEndian.Uint32(data[64:68])
	d.TokenID.SetBytes(data[68:100])
	d.L2Amount.SetBytes(data[100:132])

	return nil
}
