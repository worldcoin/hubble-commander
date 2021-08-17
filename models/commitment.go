package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	commitmentDataLength    = 137
	commitmentKeyDataLength = 36
)

type Commitment struct {
	BatchID           Uint256
	IndexInBatch      int32
	Type              txtype.TransactionType
	FeeReceiver       uint32
	CombinedSignature Signature
	PostStateRoot     common.Hash
	IncludedInBatch   *Uint256
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
	//TODO: don't serialize BatchID and IndexInBatch, that can be decoded from key
	copy(encoded[0:32], utils.PadLeft(c.BatchID.Bytes(), 32))
	//TODO: replace to uint32 in struct
	binary.BigEndian.PutUint32(encoded[32:36], uint32(c.IndexInBatch))
	encoded[36] = byte(c.Type)
	binary.BigEndian.PutUint32(encoded[37:41], c.FeeReceiver)
	copy(encoded[41:105], c.CombinedSignature.Bytes())
	copy(encoded[105:137], c.PostStateRoot.Bytes())
	copy(encoded[137:], c.Transactions)

	return encoded
}

func (c *Commitment) SetBytes(data []byte) error {
	if len(data) < commitmentDataLength {
		return ErrInvalidLength
	}

	c.BatchID.SetBytes(data[0:32])
	c.IndexInBatch = int32(binary.BigEndian.Uint32(data[32:36]))
	c.Type = txtype.TransactionType(data[36])
	c.FeeReceiver = binary.BigEndian.Uint32(data[37:41])

	err := c.CombinedSignature.SetBytes(data[41:105])
	if err != nil {
		return err
	}

	c.PostStateRoot.SetBytes(data[105:137])
	c.Transactions = data[137:]
	return nil
}

type CommitmentKey struct {
	BatchID      Uint256
	IndexInBatch uint32
}

func (c *CommitmentKey) Bytes() []byte {
	encoded := make([]byte, commitmentKeyDataLength)
	copy(encoded[0:32], utils.PadLeft(c.BatchID.Bytes(), 32))
	binary.BigEndian.PutUint32(encoded[32:36], c.IndexInBatch)

	return encoded
}

func (c *CommitmentKey) SetBytes(data []byte) error {
	if len(data) != commitmentKeyDataLength {
		return ErrInvalidLength
	}

	c.BatchID.SetBytes(data[0:32])
	c.IndexInBatch = binary.BigEndian.Uint32(data[32:36])
	return nil
}

type CommitmentWithTokenID struct {
	ID                 int32 `db:"commitment_id"`
	LeafHash           common.Hash
	Transactions       []byte `json:"-"`
	TokenID            Uint256
	FeeReceiverStateID uint32      `db:"fee_receiver"`
	CombinedSignature  Signature   `db:"combined_signature"`
	PostStateRoot      common.Hash `db:"post_state_root"`
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
