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
	rollupABI, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	require.NoError(t, err)

	batchID := models.NewUint256(1)
	commitment := models.CommitmentWithTxs{
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
	arg1, arg2, arg3, arg4, arg5 := CommitmentsToTransferAndC2TSubmitBatchFields(batchID, []models.CommitmentWithTxs{commitment})
	calldata, err := rollupABI.Pack("submitTransfer", arg1, arg2, arg3, arg4, arg5)
	require.NoError(t, err)

	decodedCommitments, err := DecodeBatchCalldata(&rollupABI, calldata)
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
