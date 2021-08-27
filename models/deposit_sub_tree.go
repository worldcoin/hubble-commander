package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type PendingDepositSubTree struct {
	ID       Uint256
	Root     common.Hash
	Deposits []DepositID
}

func (d *PendingDepositSubTree) Bytes() []byte {
	b := make([]byte, common.HashLength+depositIDDataLength*len(d.Deposits))

	copy(b[0:common.HashLength], d.Root.Bytes())

	for i := range d.Deposits {
		start := common.HashLength + i*depositIDDataLength
		end := start + depositIDDataLength
		copy(b[start:end], d.Deposits[i].Bytes())
	}

	return b
}

func (d *PendingDepositSubTree) SetBytes(data []byte) error {
	dataLength := len(data)

	if dataLength < common.HashLength || (dataLength-common.HashLength)%depositIDDataLength != 0 {
		return ErrInvalidLength
	}

	d.Root.SetBytes(data[0:common.HashLength])

	leafCount := (dataLength - common.HashLength) / depositIDDataLength

	if leafCount > 0 {
		d.Deposits = make([]DepositID, 0, leafCount)
	}

	for i := 0; i < leafCount; i++ {
		start := common.HashLength + i*depositIDDataLength
		end := start + depositIDDataLength
		leaf := DepositID{}
		err := leaf.SetBytes(data[start:end])
		if err != nil {
			return err
		}
		d.Deposits = append(d.Deposits, leaf)
	}

	return nil
}
