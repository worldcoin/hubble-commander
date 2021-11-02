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

func (c *TxCommitment) CalcBodyHash(accountRoot common.Hash) common.Hash {
	return calcBodyHash(c.FeeReceiver, c.CombinedSignature, c.Transactions, accountRoot.Bytes())
}

func (c *TxCommitment) LeafHash() common.Hash {
	return utils.HashTwo(c.PostStateRoot, *c.BodyHash)
}

func calcBodyHash(feeReceiver uint32, combinedSignature Signature, transactions, accountTreeRoot []byte) common.Hash {
	arr := make([]byte, 32+64+32+len(transactions))

	copy(arr[0:32], accountTreeRoot)
	copy(arr[32:96], combinedSignature.Bytes())
	binary.BigEndian.PutUint32(arr[124:128], feeReceiver)
	copy(arr[128:], transactions)

	return crypto.Keccak256Hash(arr)
}

type TxCommitmentWithTxs struct {
	CommitmentBase
	FeeReceiver       uint32
	CombinedSignature Signature
	Transactions      []byte
}

func (c *TxCommitmentWithTxs) CalcBodyHash(accountRoot common.Hash) common.Hash {
	return calcBodyHash(c.FeeReceiver, c.CombinedSignature, c.Transactions, accountRoot.Bytes())
}

func (c *TxCommitmentWithTxs) LeafHash() common.Hash {
	return utils.HashTwo(c.PostStateRoot, *c.BodyHash)
}

func (c *TxCommitmentWithTxs) ToTxCommitment() *TxCommitment {
	return &TxCommitment{
		CommitmentBase:    c.CommitmentBase,
		FeeReceiver:       c.FeeReceiver,
		CombinedSignature: c.CombinedSignature,
		Transactions:      c.Transactions,
	}
}
