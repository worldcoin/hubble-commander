package encoder

import (
	"math/big"
	"strings"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// DecodeBatchCalldata
//   bytes32[] stateRoots,
//   uint256[2][] signatures,
//   uint256[] feeReceivers,
//   bytes[] txss
func DecodeBatchCalldata(calldata []byte, batchID *models.Uint256) ([]DecodedCommitment, error) {
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
			ID: models.CommitmentID{
				BatchID:      *batchID,
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

func DecodeDepositBatchCalldata(calldata []byte) (*uint32, error) {
	rollupAbi, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	unpacked, err := rollupAbi.Methods["submitDeposits"].Inputs.Unpack(calldata[4:])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	vacancyProof := unpacked[1].(struct {
		PathAtDepth *big.Int   `json:"pathAtDepth"`
		Witness     [][32]byte `json:"witness"`
	})

	return ref.Uint32(uint32(vacancyProof.PathAtDepth.Uint64())), nil
}
