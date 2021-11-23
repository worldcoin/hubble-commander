package encoder

import (
	"math/big"
	"strings"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// TODO move to eth.Client and reuse rollupAbi

// DecodeBatchCalldata
//   uint256 batchID
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
