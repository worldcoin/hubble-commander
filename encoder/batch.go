package encoder

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// DecodeTransferBatchCalldata
//
//	uint256 batchID
//	bytes32[] stateRoots,
//	uint256[2][] signatures,
//	uint256[] feeReceivers,
//	bytes[] txss
func DecodeTransferBatchCalldata(rollupABI *abi.ABI, calldata []byte) ([]DecodedCommitment, error) {
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

// DecodeMMBatchCalldata
//
//	uint256 batchID,
//	bytes32[] stateRoots,
//	uint256[2][] signatures,
//	uint256[4][] meta,
//	bytes32[] withdrawRoots,
//	bytes[] txss
func DecodeMMBatchCalldata(rollupABI *abi.ABI, calldata []byte) ([]DecodedMMCommitment, error) {
	unpacked, err := rollupABI.Methods["submitMassMigration"].Inputs.Unpack(calldata[4:])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	batchID := unpacked[0].(*big.Int)
	stateRoots := unpacked[1].([][32]uint8)
	signatures := unpacked[2].([][2]*big.Int)
	meta := unpacked[3].([][4]*big.Int)
	withdrawRoots := unpacked[4].([][32]uint8)
	txss := unpacked[5].([][]uint8)

	size := len(stateRoots)

	commitments := make([]DecodedMMCommitment, size)
	for i := 0; i < size; i++ {
		mmMeta := models.NewMassMigrationMetaFromBigInts(meta[i])
		commitments[i] = DecodedMMCommitment{
			DecodedCommitment: DecodedCommitment{
				ID: models.CommitmentID{
					BatchID:      models.MakeUint256FromBig(*batchID),
					IndexInBatch: uint8(i),
				},
				StateRoot:         common.BytesToHash(stateRoots[i][:]),
				CombinedSignature: models.MakeSignatureFromBigInts(signatures[i]),
				FeeReceiver:       mmMeta.FeeReceiver,
				Transactions:      txss[i],
			},
			Meta:         mmMeta,
			WithdrawRoot: common.BytesToHash(withdrawRoots[i][:]),
		}
	}

	return commitments, nil
}

func CommitmentsToTransferAndC2TSubmitBatchFields(batchID *models.Uint256, commitments []models.CommitmentWithTxs) (
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
		commitment := commitments[i].ToTxCommitmentWithTxs()

		stateRoots = append(stateRoots, commitment.PostStateRoot)
		signatures = append(signatures, commitment.CombinedSignature.BigInts())
		feeReceivers = append(feeReceivers, new(big.Int).SetUint64(uint64(commitment.FeeReceiver)))
		transactions = append(transactions, commitment.Transactions)
	}
	return
}

//nolint:gocritic
func CommitmentsToSubmitMMBatchFields(
	batchID *models.Uint256,
	commitments []models.CommitmentWithTxs,
) (
	bigBatchID *big.Int,
	stateRoots [][32]byte,
	signatures [][2]*big.Int,
	meta [][4]*big.Int,
	bytesWithdrawRoots [][32]byte,
	transactions [][]byte,
) {
	count := len(commitments)

	bigBatchID = batchID.ToBig()
	stateRoots = make([][32]byte, 0, count)
	signatures = make([][2]*big.Int, 0, count)
	meta = make([][4]*big.Int, 0, count)
	bytesWithdrawRoots = make([][32]byte, 0, count)
	transactions = make([][]byte, 0, count)

	for i := range commitments {
		commitment := commitments[i].ToMMCommitmentWithTxs()

		stateRoots = append(stateRoots, commitment.PostStateRoot)
		signatures = append(signatures, commitment.CombinedSignature.BigInts())
		meta = append(meta, commitment.Meta.BigInts())
		bytesWithdrawRoots = append(bytesWithdrawRoots, commitment.WithdrawRoot)
		transactions = append(transactions, commitment.Transactions)
	}
	return
}
