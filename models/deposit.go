package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/utils"
)

const (
	depositDataLength             = 68
	depositInCommitmentDataLength = 101
	depositIDDataLength           = 8
)

type DepositID struct {
	BlockNumber uint32
	LogIndex    uint32
}

type Deposit struct {
	ID                   DepositID
	ToPubKeyID           uint32
	TokenID              Uint256
	L2Amount             Uint256
	IncludedInCommitment *CommitmentID
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

func (d *Deposit) Bytes() []byte {
	var b []byte

	if d.IncludedInCommitment != nil {
		b = make([]byte, depositInCommitmentDataLength)
		copy(b[68:101], d.IncludedInCommitment.Bytes())
	} else {
		b = make([]byte, depositDataLength)
	}

	binary.BigEndian.PutUint32(b[0:4], d.ToPubKeyID)
	copy(b[4:36], utils.PadLeft(d.TokenID.Bytes(), 32))
	copy(b[36:68], utils.PadLeft(d.L2Amount.Bytes(), 32))

	return b
}

func (d *Deposit) SetBytes(data []byte) error {
	if len(data) != depositDataLength && len(data) != depositInCommitmentDataLength {
		return ErrInvalidLength
	}

	if len(data) == depositInCommitmentDataLength {
		d.IncludedInCommitment = &CommitmentID{}
		err := d.IncludedInCommitment.SetBytes(data[68:101])
		if err != nil {
			return err
		}
	}

	d.ToPubKeyID = binary.BigEndian.Uint32(data[0:4])
	d.TokenID.SetBytes(data[4:36])
	d.L2Amount.SetBytes(data[36:68])

	return nil
}
