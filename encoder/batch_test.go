package encoder

import (
	"strings"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestDecodeBatchCalldata(t *testing.T) {
	//goland:noinspection GoDeprecation
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

	decodedCommitments, err := DecodeTransferBatchCalldata(&rollupABI, calldata)
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

func TestDecodeMMBatchCalldata(t *testing.T) {
	//goland:noinspection GoDeprecation
	rollupABI, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	require.NoError(t, err)

	batchID := models.NewUint256(1)
	commitments := []models.CommitmentWithTxs{{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				Type:          batchtype.MassMigration,
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       uint32(1234),
			CombinedSignature: models.MakeRandomSignature(),
		},
		Transactions: utils.RandomBytes(8),
	}}
	metas := []models.MassMigrationMeta{{
		SpokeID:     1,
		TokenID:     models.MakeUint256(1),
		Amount:      models.MakeUint256(100),
		FeeReceiver: 1,
	}}
	withdrawRoots := []common.Hash{utils.RandomHash()}

	arg1, arg2, arg3, arg4, arg5, arg6 := CommitmentsToSubmitMassMigrationBatchFields(batchID, commitments, metas, withdrawRoots)
	calldata, err := rollupABI.Pack("submitMassMigration", arg1, arg2, arg3, arg4, arg5, arg6)
	require.NoError(t, err)

	decodedCommitments, err := DecodeMMBatchCalldata(&rollupABI, calldata)
	require.NoError(t, err)
	require.Equal(t, 1, len(decodedCommitments))

	decoded := &decodedCommitments[0]
	require.Equal(t, commitments[0].PostStateRoot, decoded.StateRoot)
	require.Equal(t, commitments[0].CombinedSignature, decoded.CombinedSignature)
	require.Equal(t, metas[0], *decoded.Meta)
	require.Equal(t, metas[0].FeeReceiver, decoded.FeeReceiver)
	require.Equal(t, withdrawRoots[0], decoded.WithdrawRoot)
	require.Equal(t, commitments[0].Transactions, decoded.Transactions)
	require.Equal(t, *batchID, decoded.ID.BatchID)
	require.EqualValues(t, 0, decoded.ID.IndexInBatch)
}
