package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
)

func TestStoredTxReceipt_Bytes(t *testing.T) {
	txReceipt := StoredTxReceipt{
		Hash: utils.RandomHash(),
		CommitmentID: &CommitmentID{
			BatchID:      MakeUint256(10),
			IndexInBatch: 2,
		},
		ToStateID:    ref.Uint32(12),
		ErrorMessage: ref.String("some error message"),
	}

	bytes := txReceipt.Bytes()

	decodedStoredTxReceipt := StoredTxReceipt{}
	err := decodedStoredTxReceipt.SetBytes(bytes)
	require.NoError(t, err)
	require.EqualValues(t, txReceipt, decodedStoredTxReceipt)
}
