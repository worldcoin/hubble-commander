package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	storedTxCommitmentBodyLength          = 4 + 64 + 33
	storedDepositCommitmentBodyBaseLength = 32 + 32
)

type StoredCommitmentBody interface {
	ByteEncoder
	BytesLen() int
}

func NewStoredCommitmentBody(commitmentType batchtype.BatchType) (StoredCommitmentBody, error) {
	// nolint:exhaustive
	switch commitmentType {
	case batchtype.Deposit:
		return new(StoredDepositCommitmentBody), nil
	case batchtype.Transfer, batchtype.Create2Transfer, batchtype.MassMigration:
		return new(StoredTxCommitmentBody), nil
	default:
		return nil, errors.Errorf("unsupported commitment type: %s", commitmentType)
	}
}

type StoredTxCommitmentBody struct {
	FeeReceiver       uint32
	CombinedSignature Signature
	BodyHash          *common.Hash
}

func (c *StoredTxCommitmentBody) Bytes() []byte {
	b := make([]byte, c.BytesLen())
	binary.BigEndian.PutUint32(b[0:4], c.FeeReceiver)
	copy(b[4:68], c.CombinedSignature.Bytes())
	copy(b[68:101], EncodeHashPointer(c.BodyHash))
	return b
}

func (c *StoredTxCommitmentBody) SetBytes(data []byte) error {
	if len(data) != storedTxCommitmentBodyLength {
		return ErrInvalidLength
	}
	err := c.CombinedSignature.SetBytes(data[4:68])
	if err != nil {
		return err
	}

	c.FeeReceiver = binary.BigEndian.Uint32(data[0:4])
	c.BodyHash = decodeHashPointer(data[68:101])
	return nil
}

func (c *StoredTxCommitmentBody) BytesLen() int {
	return storedTxCommitmentBodyLength
}

type StoredDepositCommitmentBody struct {
	SubTreeID   Uint256
	SubTreeRoot common.Hash
	Deposits    []PendingDeposit
}

func (c *StoredDepositCommitmentBody) Bytes() []byte {
	b := make([]byte, c.BytesLen())
	copy(b[0:32], c.SubTreeID.Bytes())
	copy(b[32:64], c.SubTreeRoot.Bytes())

	startIndex := storedDepositCommitmentBodyBaseLength
	for i := range c.Deposits {
		startIndex += copy(b[startIndex:startIndex+depositDataLength], c.Deposits[i].Bytes())
	}

	return b
}

func (c *StoredDepositCommitmentBody) SetBytes(data []byte) error {
	overallDepositsLength := len(data) - storedDepositCommitmentBodyBaseLength
	if len(data) <= storedDepositCommitmentBodyBaseLength || overallDepositsLength%depositDataLength != 0 {
		return ErrInvalidLength
	}

	depositCount := overallDepositsLength / depositDataLength
	c.Deposits = make([]PendingDeposit, 0, depositCount)

	startIndex := storedDepositCommitmentBodyBaseLength
	for i := 0; i < depositCount; i++ {
		endIndex := startIndex + depositDataLength
		deposit := PendingDeposit{}
		err := deposit.SetBytes(data[startIndex:endIndex])
		if err != nil {
			return err
		}
		c.Deposits = append(c.Deposits, deposit)
		startIndex = endIndex
	}

	c.SubTreeID.SetBytes(data[0:32])
	c.SubTreeRoot.SetBytes(data[32:64])
	return nil
}

func (c *StoredDepositCommitmentBody) BytesLen() int {
	return storedDepositCommitmentBodyBaseLength + len(c.Deposits)*depositDataLength
}
