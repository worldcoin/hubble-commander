package encoder

import (
	"encoding/binary"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
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

func (c *DecodedCommitment) BodyHash(accountRoot common.Hash) common.Hash {
	arr := make([]byte, 32+64+32+len(c.Transactions))

	copy(arr[0:32], accountRoot.Bytes())
	copy(arr[32:96], c.CombinedSignature.Bytes())
	binary.BigEndian.PutUint32(arr[124:128], c.FeeReceiver)
	copy(arr[128:], c.Transactions)

	return crypto.Keccak256Hash(arr)
}

func (c *DecodedCommitment) LeafHash(accountRoot common.Hash) common.Hash {
	return utils.HashTwo(c.StateRoot, c.BodyHash(accountRoot))
}

func CommitmentToCalldataFields(commitments []models.Commitment) (
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	feeReceivers []*big.Int,
	transactions [][]byte,
) {
	count := len(commitments)

	stateRoots = make([][32]byte, 0, count)
	signatures = make([][2]*big.Int, 0, count)
	feeReceivers = make([]*big.Int, 0, count)
	transactions = make([][]byte, 0, count)

	for i := range commitments {
		stateRoots = append(stateRoots, commitments[i].PostStateRoot)
		signatures = append(signatures, commitments[i].CombinedSignature.BigInts())
		feeReceivers = append(feeReceivers, new(big.Int).SetUint64(uint64(commitments[i].FeeReceiver)))
		transactions = append(transactions, commitments[i].Transactions)
	}
	return
}
