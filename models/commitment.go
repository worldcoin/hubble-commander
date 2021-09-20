package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	commitmentDataLength   = 101
	commitmentIDDataLength = 33
)

var CommitmentPrefix = getBadgerHoldPrefix(Commitment{})

type Commitment struct {
	ID                CommitmentID
	Type              batchtype.BatchType
	FeeReceiver       uint32
	CombinedSignature Signature
	PostStateRoot     common.Hash
	Transactions      []byte
}

func (c *Commitment) BodyHash(accountRoot common.Hash) common.Hash {
	return calcBodyHash(c.FeeReceiver, c.CombinedSignature, c.Transactions, accountRoot.Bytes())
}

func (c *Commitment) LeafHash(accountRoot common.Hash) common.Hash {
	return utils.HashTwo(c.PostStateRoot, c.BodyHash(accountRoot))
}

func (c *Commitment) Bytes() []byte {
	encoded := make([]byte, commitmentDataLength+len(c.Transactions))
	encoded[0] = byte(c.Type)
	binary.BigEndian.PutUint32(encoded[1:5], c.FeeReceiver)
	copy(encoded[5:69], c.CombinedSignature.Bytes())
	copy(encoded[69:101], c.PostStateRoot.Bytes())
	copy(encoded[101:], c.Transactions)

	return encoded
}

func (c *Commitment) SetBytes(data []byte) error {
	if len(data) < commitmentDataLength {
		return ErrInvalidLength
	}
	err := c.CombinedSignature.SetBytes(data[5:69])
	if err != nil {
		return err
	}

	c.Type = batchtype.BatchType(data[0])
	c.FeeReceiver = binary.BigEndian.Uint32(data[1:5])
	c.PostStateRoot.SetBytes(data[69:101])
	c.Transactions = data[101:]
	return nil
}

type CommitmentID struct {
	BatchID      Uint256
	IndexInBatch uint8
}

func (c *CommitmentID) Bytes() []byte {
	encoded := make([]byte, commitmentIDDataLength)
	copy(encoded[0:32], utils.PadLeft(c.BatchID.Bytes(), 32))
	encoded[32] = c.IndexInBatch

	return encoded
}

func (c *CommitmentID) SetBytes(data []byte) error {
	if len(data) != commitmentIDDataLength {
		return ErrInvalidLength
	}

	c.BatchID.SetBytes(data[0:32])
	c.IndexInBatch = data[32]
	return nil
}

type CommitmentWithTokenID struct {
	ID                 CommitmentID
	LeafHash           common.Hash
	Transactions       []byte `json:"-"`
	TokenID            Uint256
	FeeReceiverStateID uint32
	CombinedSignature  Signature
	PostStateRoot      common.Hash
}

func (c *CommitmentWithTokenID) BodyHash(accountRoot common.Hash) common.Hash {
	return calcBodyHash(c.FeeReceiverStateID, c.CombinedSignature, c.Transactions, accountRoot.Bytes())
}

func (c *CommitmentWithTokenID) CalcLeafHash(accountTreeRoot *common.Hash) common.Hash {
	bodyHash := calcBodyHash(c.FeeReceiverStateID, c.CombinedSignature, c.Transactions, accountTreeRoot.Bytes())
	return utils.HashTwo(c.PostStateRoot, bodyHash)
}

func calcBodyHash(feeReceiver uint32, combinedSignature Signature, transactions, accountTreeRoot []byte) common.Hash {
	arr := make([]byte, 32+64+32+len(transactions))

	copy(arr[0:32], accountTreeRoot)
	copy(arr[32:96], combinedSignature.Bytes())
	binary.BigEndian.PutUint32(arr[124:128], feeReceiver)
	copy(arr[128:], transactions)

	return crypto.Keccak256Hash(arr)
}
