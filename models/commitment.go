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
	FeeReceiver       uint32       `db:"fee_receiver"`
	CombinedSignature Signature    `db:"combined_signature"`
	PostStateRoot     common.Hash  `db:"post_state_root"`
	AccountTreeRoot   *common.Hash `db:"account_tree_root"`
	IncludedInBatch   *common.Hash `db:"included_in_batch"`
}

func (c *Commitment) BodyHash() common.Hash {
	return calcBodyHash(c.FeeReceiver, c.CombinedSignature, c.Transactions, c.AccountTreeRoot.Bytes())
}

func (c *Commitment) LeafHash() common.Hash {
	return utils.HashTwo(c.PostStateRoot, c.BodyHash())
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

func (c *CommitmentWithTokenID) CalcLeafHash(accountTreeRoot *common.Hash) common.Hash {
	bodyHash := calcBodyHash(c.FeeReceiverStateID, c.CombinedSignature, c.Transactions, accountTreeRoot.Bytes())
	return utils.HashTwo(c.PostStateRoot, bodyHash)
}

func calcBodyHash(feeReceiver uint32, combinedSignature Signature, transactions, accountTreeRoot []byte) common.Hash {
	arr := make([]byte, 32+64+32+len(transactions))

	copy(arr[0:32], accountTreeRoot)
	copy(arr[32:64], utils.PadLeft(combinedSignature[0].Bytes(), 32))
	copy(arr[64:96], utils.PadLeft(combinedSignature[1].Bytes(), 32))
	binary.BigEndian.PutUint32(arr[124:128], feeReceiver)
	copy(arr[128:], transactions)

	return crypto.Keccak256Hash(arr)
}
