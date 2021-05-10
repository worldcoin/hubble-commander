package encoder

import (
	"strings"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"
)

func TestDecodeTransferBatch(t *testing.T) {
	rollupAbi, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	require.NoError(t, err)

	commitments := []models.Commitment{
		{
			Type:              txtype.Transfer,
			Transactions:      utils.RandomBytes(12),
			FeeReceiver:       uint32(1234),
			CombinedSignature: models.MakeRandomSignature(),
			PostStateRoot:     utils.RandomHash(),
		},
	}
	arg1, arg2, arg3, arg4 := CommitmentToCalldataFields(commitments)
	calldata, err := rollupAbi.Pack("submitTransfer", arg1, arg2, arg3, arg4)
	require.NoError(t, err)

	decoded, err := DecodeTransferBatch(calldata)
	require.NoError(t, err)

	require.Equal(t, len(commitments), len(decoded))
	require.Equal(t, commitments[0].PostStateRoot, decoded[0].StateRoot)
	require.Equal(t, commitments[0].Transactions, decoded[0].Transactions)
	require.Equal(t, commitments[0].FeeReceiver, decoded[0].FeeReceiver)

}
