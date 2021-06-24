package encoder

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type DecodedCommitment struct {
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

// DecodeBatchCalldata
//   bytes32[] stateRoots,
//   uint256[2][] signatures,
//   uint256[] feeReceivers,
//   bytes[] txss
func DecodeBatchCalldata(calldata []byte) ([]DecodedCommitment, error) {
	rollupAbi, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	unpacked, err := rollupAbi.Methods["submitTransfer"].Inputs.Unpack(calldata[4:])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	stateRoots := unpacked[0].([][32]uint8)
	signatures := unpacked[1].([][2]*big.Int)
	feeReceivers := unpacked[2].([]*big.Int)
	txss := unpacked[3].([][]uint8)

	size := len(stateRoots)

	commitments := make([]DecodedCommitment, size)
	for i := 0; i < size; i++ {
		commitments[i] = DecodedCommitment{
			StateRoot:         common.BytesToHash(stateRoots[i][:]),
			CombinedSignature: models.MakeSignatureFromBigInts(signatures[i]),
			FeeReceiver:       uint32(feeReceivers[i].Uint64()),
			Transactions:      txss[i],
		}
	}

	return commitments, nil
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
