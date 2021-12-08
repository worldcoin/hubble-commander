package stored

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
)

func TestTxReceipt_Bytes(t *testing.T) {
	txReceipt := TxReceipt{
		Hash: utils.RandomHash(),
		CommitmentID: &models.CommitmentID{
			BatchID:      models.MakeUint256(10),
			IndexInBatch: 2,
		},
		ToStateID:    ref.Uint32(12),
		ErrorMessage: ref.String("some error message"),
	}

	bytes := txReceipt.Bytes()

	decodedTxReceipt := TxReceipt{}
	err := decodedTxReceipt.SetBytes(bytes)
	require.NoError(t, err)
	require.EqualValues(t, txReceipt, decodedTxReceipt)
}
