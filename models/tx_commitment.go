package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type TxCommitment struct {
	CommitmentBase
	FeeReceiver       uint32
	CombinedSignature Signature
	BodyHash          *common.Hash
}

func (c *TxCommitment) GetCommitmentBase() CommitmentBase {
	return c.CommitmentBase
}

func (c *TxCommitment) GetBodyHash() common.Hash {
	return *c.BodyHash
}

func (c *TxCommitment) LeafHash() common.Hash {
	return utils.HashTwo(c.PostStateRoot, *c.BodyHash)
}

func (c *TxCommitment) ToTxCommitment() *TxCommitment {
	return c
}

func (c *TxCommitment) ToMMCommitment() *MMCommitment {
	panic("cannot cast TxCommitment to MMCommitment")
}

func (c *TxCommitment) ToDepositCommitment() *DepositCommitment {
	panic("cannot cast TxCommitment to DepositCommitment")
}

type TxCommitmentWithTxs struct {
	TxCommitment
	Transactions []byte
}

func (c *TxCommitmentWithTxs) SetBodyHash(accountRoot common.Hash) {
	c.BodyHash = calcBodyHash(c.FeeReceiver, c.CombinedSignature, c.Transactions, accountRoot.Bytes())
}

func (c *TxCommitmentWithTxs) CalcBodyHash(accountRoot common.Hash) *common.Hash {
	return calcBodyHash(c.FeeReceiver, c.CombinedSignature, c.Transactions, accountRoot.Bytes())
}

func (c *TxCommitmentWithTxs) ToTxCommitmentWithTxs() *TxCommitmentWithTxs {
	return c
}

func (c *TxCommitmentWithTxs) ToMMCommitmentWithTxs() *MMCommitmentWithTxs {
	panic("cannot cast TxCommitmentWithTxs to MMCommitmentWithTxs")
}

func calcBodyHash(feeReceiver uint32, combinedSignature Signature, transactions, accountTreeRoot []byte) *common.Hash {
	arr := make([]byte, 32+64+32+len(transactions))

	copy(arr[0:32], accountTreeRoot)
	copy(arr[32:96], combinedSignature.Bytes())
	binary.BigEndian.PutUint32(arr[124:128], feeReceiver)
	copy(arr[128:], transactions)

	return ref.Hash(crypto.Keccak256Hash(arr))
}
