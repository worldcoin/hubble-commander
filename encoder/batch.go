package encoder

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// DecodeBatchCalldata
//   uint256 batchID
//   bytes32[] stateRoots,
//   uint256[2][] signatures,
//   uint256[] feeReceivers,
//   bytes[] txss
func DecodeBatchCalldata(rollupABI *abi.ABI, calldata []byte) ([]DecodedCommitment, error) {
	unpacked, err := rollupABI.Methods["submitTransfer"].Inputs.Unpack(calldata[4:])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	batchID := unpacked[0].(*big.Int)
	stateRoots := unpacked[1].([][32]uint8)
	signatures := unpacked[2].([][2]*big.Int)
	feeReceivers := unpacked[3].([]*big.Int)
	txss := unpacked[4].([][]uint8)

	size := len(stateRoots)

	commitments := make([]DecodedCommitment, size)
	for i := 0; i < size; i++ {
		commitments[i] = DecodedCommitment{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256FromBig(*batchID),
				IndexInBatch: uint8(i),
			},
			StateRoot:         common.BytesToHash(stateRoots[i][:]),
			CombinedSignature: models.MakeSignatureFromBigInts(signatures[i]),
			FeeReceiver:       uint32(feeReceivers[i].Uint64()),
			Transactions:      txss[i],
		}
	}

	return commitments, nil
}

func CommitmentsToSubmitBatchFields(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (
	bigBatchID *big.Int,
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	feeReceivers []*big.Int,
	transactions [][]byte,
) {
	bigBatchID = batchID.ToBig()

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
