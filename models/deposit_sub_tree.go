package models

import (
	"github.com/ethereum/go-ethereum/common"
)

var PendingDepositSubTreePrefix = getBadgerHoldPrefix(PendingDepositSubTree{})

type PendingDepositSubTree struct {
	ID       Uint256
	Root     common.Hash
	Deposits []PendingDeposit
}

func (d *PendingDepositSubTree) Bytes() []byte {
	b := make([]byte, common.HashLength+depositDataLength*len(d.Deposits))

	copy(b[0:common.HashLength], d.Root.Bytes())

	for i := range d.Deposits {
		start := common.HashLength + i*depositDataLength
		end := start + depositDataLength
		copy(b[start:end], d.Deposits[i].Bytes())
	}

	return b
}

func (d *PendingDepositSubTree) SetBytes(data []byte) error {
	dataLength := len(data)

	if dataLength < common.HashLength || (dataLength-common.HashLength)%depositDataLength != 0 {
		return ErrInvalidLength
	}

	d.Root.SetBytes(data[0:common.HashLength])

	leafCount := (dataLength - common.HashLength) / depositDataLength

	if leafCount > 0 {
		d.Deposits = make([]PendingDeposit, 0, leafCount)
	} else {
		d.Deposits = nil
	}

	for i := 0; i < leafCount; i++ {
		start := common.HashLength + i*depositDataLength
		end := start + depositDataLength
		leaf := PendingDeposit{}
		err := leaf.SetBytes(data[start:end])
		if err != nil {
			return err
		}
		d.Deposits = append(d.Deposits, leaf)
	}

	return nil
}
