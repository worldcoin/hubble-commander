package stored

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestBatchedTx_GenericTransactionRoundTrip(t *testing.T) {
	batchedTx := BatchedTx{
		PendingTx: PendingTx{
			Hash: utils.RandomHash(),

			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
			ReceiveTime: nil,

			Body: &TxTransferBody{
				ToStateID: 2,
			},
		},
		ID: models.CommitmentSlot{
			BatchID:           models.MakeUint256(1),
			IndexInBatch:      2,
			IndexInCommitment: 3,
		},
	}

	origBytes := batchedTx.Bytes()
	roundTripBytes := NewBatchedTx(batchedTx.ToGenericTransaction()).Bytes()
	require.Equal(t, origBytes, roundTripBytes)
}
