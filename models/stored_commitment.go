package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	storedCommitmentBytesLength       = 66
	storedCommitmentTxBodyLength      = 68
	storedCommitmentDepositBodyLength = 64
)

var (
	StoredCommitmentName                = getTypeName(StoredCommitment{})
	StoredCommitmentPrefix              = getBadgerHoldPrefix(StoredCommitment{})
	errInvalidStoredCommitmentIndexType = errors.New("invalid StoredCommitment index type")
)

type StoredCommitment struct {
	ID            CommitmentID
	Type          batchtype.BatchType
	PostStateRoot common.Hash

	Body StoredCommitmentBody
}

func MakeStoredCommitmentFromTransaction(t *GenericTransaction) StoredCommitment {
	return StoredCommitment{
		ID: CommitmentID{
			BatchID:      Uint256{},
			IndexInBatch: 0,
		},
		Type:          c.Type,
		PostStateRoot: c.PostStateRoot,
		Body:          nil,
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
	if len(data) < storedCommitmentTxBodyLength {
		return ErrInvalidLength
	}

	err := c.ID.SetBytes(data[0:33])
	if err != nil {
		return err
	}

	c.Type = batchtype.BatchType(data[33])
	c.PostStateRoot.SetBytes(data[34:66])

	body, err := commitmentBody(data[66:], c.Type)
	if err != nil {
		return err
	}

	c.Body = body

	return nil
}

func (c *StoredCommitment) BytesLen() int {
	return storedCommitmentBytesLength + c.Body.BytesLen()
}

func commitmentBody(data []byte, commitmentType batchtype.BatchType) (StoredCommitmentBody, error) {
	switch commitmentType {
	case batchtype.Deposit:
		body := new(StoredCommitmentDepositBody)
		err := body.SetBytes(data)
		return body, err
	case batchtype.Transfer:
	case batchtype.Create2Transfer:
		body := new(StoredCommitmentTxBody)
		err := body.SetBytes(data)
		return body, err
	case batchtype.MassMigration:
		return nil, errors.Errorf("unsupported commitment type: %s", commitmentType)
	}

	return nil, nil
}

type StoredCommitmentBody interface {
	ByteEncoder
	BytesLen() int
}

type StoredCommitmentTxBody struct {
	FeeReceiver       uint32
	CombinedSignature Signature
	Transactions      []byte
}

func (c *StoredCommitmentTxBody) Bytes() []byte {
	b := make([]byte, storedCommitmentTxBodyLength)
	binary.BigEndian.PutUint32(b[0:4], c.FeeReceiver)
	copy(b[4:68], c.CombinedSignature.Bytes())
	copy(b[68:], c.Transactions)
	return b
}

func (c *StoredCommitmentTxBody) SetBytes(data []byte) error {
	c.FeeReceiver = binary.BigEndian.Uint32(data[0:4])
	err := c.CombinedSignature.SetBytes(data[4:68])
	if err != nil {
		return err
	}
	c.Transactions = data[68:]
	return nil
}

func (c *StoredCommitmentTxBody) BytesLen() int {
	return storedCommitmentTxBodyLength + len(c.Transactions)
}

type StoredCommitmentDepositBody struct {
	SubTreeID   Uint256
	SubTreeRoot common.Hash
	Deposits    []PendingDeposit
}

func (c *StoredCommitmentDepositBody) Bytes() []byte {
	b := make([]byte, storedCommitmentTxBodyLength)
	copy(b[0:32], c.SubTreeID.Bytes())
	copy(b[32:64], c.SubTreeRoot.Bytes())

	for i := range c.Deposits {
		start := storedCommitmentDepositBodyLength + i*depositDataLength
		end := start + depositDataLength
		copy(b[start:end], c.Deposits[i].Bytes())
	}

	return b
}

func (c *StoredCommitmentDepositBody) SetBytes(data []byte) error {
	dataLength := len(data)

	// TODO-SC check if commitment can have 0 deposits?
	if dataLength < storedCommitmentDepositBodyLength || (dataLength-storedCommitmentDepositBodyLength)%depositDataLength != 0 {
		return ErrInvalidLength
	}

	c.SubTreeID.SetBytes(data[0:32])
	c.SubTreeRoot.SetBytes(data[32:64])

	depositCount := (dataLength - storedCommitmentDepositBodyLength) / depositDataLength

	// TODO-SC check the TODO above
	if depositCount > 0 {
		c.Deposits = make([]PendingDeposit, 0, depositCount)
	} else {
		c.Deposits = nil
	}

	for i := 0; i < depositCount; i++ {
		start := storedCommitmentDepositBodyLength + i*depositDataLength
		end := start + depositDataLength
		deposit := PendingDeposit{}
		err := deposit.SetBytes(data[start:end])
		if err != nil {
			return err
		}
		c.Deposits = append(c.Deposits, deposit)
	}

	return nil
}

func (c *StoredCommitmentDepositBody) BytesLen() int {
	return storedCommitmentDepositBodyLength + len(c.Deposits)
}
