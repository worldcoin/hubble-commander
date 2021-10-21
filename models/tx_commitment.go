package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type TxCommitment struct {
	CommitmentBase
	FeeReceiver       uint32
	CombinedSignature Signature
	Transactions      []byte
}

func (c *TxCommitment) BodyHash(accountRoot common.Hash) common.Hash {
	return calcBodyHash(c.FeeReceiver, c.CombinedSignature, c.Transactions, accountRoot.Bytes())
}

func (c *TxCommitment) LeafHash(accountRoot common.Hash) common.Hash {
	return utils.HashTwo(c.PostStateRoot, c.BodyHash(accountRoot))
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