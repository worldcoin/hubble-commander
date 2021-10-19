package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	storedCommitmentBytesLength       = 66
	storedTxCommitmentBodyLength      = 68
	storedDepositCommitmentBodyLength = 64
)

var StoredCommitmentPrefix = getBadgerHoldPrefix(StoredCommitment{})

type StoredCommitment struct {
	CommitmentBase
	Body StoredCommitmentBody
}

func MakeStoredCommitmentFromTxCommitment(c *TxCommitment) StoredCommitment {
	return StoredCommitment{
		CommitmentBase: c.CommitmentBase,
		Body: &StoredTxCommitmentBody{
			FeeReceiver:       c.FeeReceiver,
			CombinedSignature: c.CombinedSignature,
			Transactions:      c.Transactions,
		},
	}
}

func MakeStoredCommitmentFromDepositCommitment(c *DepositCommitment) StoredCommitment {
	return StoredCommitment{
		CommitmentBase: c.CommitmentBase,
		Body: &StoredDepositCommitmentBody{
			SubTreeID:   c.SubTreeID,
			SubTreeRoot: c.SubTreeRoot,
			Deposits:    c.Deposits,
		},
	}
}

func (c *StoredCommitment) Bytes() []byte {
	b := make([]byte, c.BytesLen())
	copy(b[0:33], c.ID.Bytes())
	b[33] = byte(c.Type)
	copy(b[34:66], c.PostStateRoot.Bytes())
	copy(b[66:], c.Body.Bytes())

	return b
}

func (c *StoredCommitment) SetBytes(data []byte) error {
	if len(data) < storedCommitmentBytesLength {
		return ErrInvalidLength
	}
	err := c.ID.SetBytes(data[0:33])
	if err != nil {
		return err
	}

	commitmentType := batchtype.BatchType(data[33])
	body, err := commitmentBody(data[66:], commitmentType)
	if err != nil {
		return err
	}

	c.Type = commitmentType
	c.PostStateRoot.SetBytes(data[34:66])
	c.Body = body
	return nil
}

func (c *StoredCommitment) BytesLen() int {
	return storedCommitmentBytesLength + c.Body.BytesLen()
}

func (c *StoredCommitment) ToTxCommitment() *TxCommitment {
	txCommitmentBody, ok := c.Body.(*StoredTxCommitmentBody)
	if !ok {
		panic("invalid TxCommitment body type")
	}

	return &TxCommitment{
		CommitmentBase:    c.CommitmentBase,
		FeeReceiver:       txCommitmentBody.FeeReceiver,
		CombinedSignature: txCommitmentBody.CombinedSignature,
		Transactions:      txCommitmentBody.Transactions,
	}
}

func (c *StoredCommitment) ToDepositCommitment() *DepositCommitment {
	depositCommitmentBody, ok := c.Body.(*StoredDepositCommitmentBody)
	if !ok {
		panic("invalid DepositCommitment body type")
	}

	return &DepositCommitment{
		CommitmentBase: c.CommitmentBase,
		SubTreeID:      depositCommitmentBody.SubTreeID,
		SubTreeRoot:    depositCommitmentBody.SubTreeRoot,
		Deposits:       depositCommitmentBody.Deposits,
	}
}

func commitmentBody(data []byte, commitmentType batchtype.BatchType) (StoredCommitmentBody, error) {
	switch commitmentType {
	case batchtype.Deposit:
		body := new(StoredDepositCommitmentBody)
		err := body.SetBytes(data)
		return body, err
	case batchtype.Transfer, batchtype.Create2Transfer:
		body := new(StoredTxCommitmentBody)
		err := body.SetBytes(data)
		return body, err
	case batchtype.Genesis, batchtype.MassMigration:
		return nil, errors.Errorf("unsupported commitment type: %s", commitmentType)
	}

	return nil, nil
}

type StoredCommitmentBody interface {
	ByteEncoder
	BytesLen() int
}

type StoredTxCommitmentBody struct {
	FeeReceiver       uint32
	CombinedSignature Signature
	Transactions      []byte
}

func (c *StoredTxCommitmentBody) Bytes() []byte {
	b := make([]byte, c.BytesLen())
	binary.BigEndian.PutUint32(b[0:4], c.FeeReceiver)
	copy(b[4:68], c.CombinedSignature.Bytes())
	copy(b[68:], c.Transactions)
	return b
}

func (c *StoredTxCommitmentBody) SetBytes(data []byte) error {
	err := c.CombinedSignature.SetBytes(data[4:68])
	if err != nil {
		return err
	}

	c.FeeReceiver = binary.BigEndian.Uint32(data[0:4])
	c.Transactions = data[68:]
	return nil
}

func (c *StoredTxCommitmentBody) BytesLen() int {
	return storedTxCommitmentBodyLength + len(c.Transactions)
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

	startIndex := storedDepositCommitmentBodyLength
	for i := range c.Deposits {
		startIndex += copy(b[startIndex:startIndex+depositDataLength], c.Deposits[i].Bytes())
	}

	return b
}

func (c *StoredDepositCommitmentBody) SetBytes(data []byte) error {
	if len(data) <= storedDepositCommitmentBodyLength || (len(data)-storedDepositCommitmentBodyLength)%depositDataLength != 0 {
		return ErrInvalidLength
	}

	depositCount := (len(data) - storedDepositCommitmentBodyLength) / depositDataLength
	c.Deposits = make([]PendingDeposit, 0, depositCount)

	startIndex := storedDepositCommitmentBodyLength
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
	return storedDepositCommitmentBodyLength + len(c.Deposits)*depositDataLength
}
