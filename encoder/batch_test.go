package encoder

import (
	"strings"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"
)

func TestDecodeBatchCalldata(t *testing.T) {
	rollupAbi, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	require.NoError(t, err)

	batchID := models.NewUint256(1)
	commitment := models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				Type:          batchtype.Transfer,
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       uint32(1234),
			CombinedSignature: models.MakeRandomSignature(),
		},
		Transactions: utils.RandomBytes(12),
	}
	arg1, arg2, arg3, arg4 := CommitmentToCalldataFields([]models.TxCommitmentWithTxs{commitment})
	calldata, err := rollupAbi.Pack("submitTransfer", arg1, arg2, arg3, arg4)
	require.NoError(t, err)

	decodedCommitments, err := DecodeBatchCalldata(calldata, batchID)
	require.NoError(t, err)
	require.Equal(t, 1, len(decodedCommitments))

	decoded := &decodedCommitments[0]
	require.Equal(t, commitment.PostStateRoot, decoded.StateRoot)
	require.Equal(t, commitment.CombinedSignature, decoded.CombinedSignature)
	require.Equal(t, commitment.FeeReceiver, decoded.FeeReceiver)
	require.Equal(t, commitment.Transactions, decoded.Transactions)
	require.Equal(t, *batchID, decoded.ID.BatchID)
	require.EqualValues(t, 0, decoded.ID.IndexInBatch)
}
