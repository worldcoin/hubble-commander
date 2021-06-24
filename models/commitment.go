package models

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Commitment struct {
	ID                int32 `db:"commitment_id"`
	Type              txtype.TransactionType
	Transactions      []byte
	FeeReceiver       uint32      `db:"fee_receiver"`
	CombinedSignature Signature   `db:"combined_signature"`
	PostStateRoot     common.Hash `db:"post_state_root"`
	IncludedInBatch   *Uint256    `db:"included_in_batch"`
}

func (c *Commitment) BodyHash(accountRoot common.Hash) common.Hash {
	return calcBodyHash(c.FeeReceiver, c.CombinedSignature, c.Transactions, accountRoot.Bytes())
}

func (c *Commitment) LeafHash(accountRoot common.Hash) common.Hash {
	return utils.HashTwo(c.PostStateRoot, c.BodyHash(accountRoot))
}

type CommitmentWithTokenID struct {
	ID                 int32 `db:"commitment_id"`
	LeafHash           common.Hash
	Transactions       []byte      `json:"-"`
	TokenID            Uint256     `db:"token_index"`
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
