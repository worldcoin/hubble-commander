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
	arr := make([]byte, 32+64+32+len(c.Transactions))

	copy(arr[0:32], c.AccountTreeRoot.Bytes())
	copy(arr[32:64], utils.PadLeft(c.CombinedSignature[0].Bytes(), 32))
	copy(arr[64:96], utils.PadLeft(c.CombinedSignature[1].Bytes(), 32))
	binary.BigEndian.PutUint32(arr[124:128], c.FeeReceiver)
	copy(arr[128:], c.Transactions)

	return crypto.Keccak256Hash(arr)
}

func (c *Commitment) LeafHash() common.Hash {
	return utils.HashTwo(c.PostStateRoot, c.BodyHash())
}

type CommitmentWithTokenID struct {
	ID       int32       `db:"commitment_id"`
	LeafHash common.Hash //commitment.LeafHash()
	//from feeReceiverStateID -> getStateLeafByPath
	TokenID            Uint256     `db:"token_index"`
	FeeReceiverStateID uint32      `db:"fee_receiver"`
	CombinedSignature  Signature   `db:"combined_signature"`
	PostStateRoot      common.Hash `db:"post_state_root"`
}
