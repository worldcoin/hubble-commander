package encoder

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type DecodedCommitment struct {
	ID                models.CommitmentID
	StateRoot         common.Hash
	CombinedSignature models.Signature
	FeeReceiver       uint32
	Transactions      []byte
}

func (c *DecodedCommitment) BodyHash(accountRoot common.Hash) *common.Hash {
	arr := make([]byte, 32+64+32+len(c.Transactions))

	copy(arr[0:32], accountRoot.Bytes())
	copy(arr[32:96], c.CombinedSignature.Bytes())
	binary.BigEndian.PutUint32(arr[124:128], c.FeeReceiver)
	copy(arr[128:], c.Transactions)

	return ref.Hash(crypto.Keccak256Hash(arr))
}

func (c *DecodedCommitment) LeafHash(accountRoot common.Hash) common.Hash {
	return utils.HashTwo(c.StateRoot, *c.BodyHash(accountRoot))
}

type DecodedMassMigrationCommitment struct {
	ID                models.CommitmentID
	StateRoot         common.Hash
	CombinedSignature models.Signature
	Meta              models.MassMigrationMeta
	WithdrawRoot      common.Hash
	Transactions      []byte
}

func (c *DecodedMassMigrationCommitment) BodyHash(accountRoot common.Hash) *common.Hash {
	arr := make([]byte, 32+64+32+32+32+32+32+len(c.Transactions))

	copy(arr[0:32], accountRoot.Bytes())
	copy(arr[32:96], c.CombinedSignature.Bytes())
	binary.BigEndian.PutUint32(arr[124:128], c.Meta.SpokeID)
	copy(arr[128:160], c.WithdrawRoot.Bytes())
	copy(arr[160:192], c.Meta.TokenID.Bytes())
	copy(arr[192:224], c.Meta.Amount.Bytes())
	binary.BigEndian.PutUint32(arr[252:256], c.Meta.FeeReceiver)
	copy(arr[256:], c.Transactions)

	return ref.Hash(crypto.Keccak256Hash(arr))
}

func (c *DecodedMassMigrationCommitment) LeafHash(accountRoot common.Hash) common.Hash {
	return utils.HashTwo(c.StateRoot, *c.BodyHash(accountRoot))
}
