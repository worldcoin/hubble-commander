package models

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
)

func TestTransactionBase_Bytes(t *testing.T) {
	transactionBase := &TransactionBase{
		Hash:         utils.RandomHash(),
		TxType:       txtype.Transfer,
		FromStateID:  11,
		Amount:       MakeUint256(10),
		Fee:          MakeUint256(111),
		Nonce:        MakeUint256(1),
		Signature:    Signature{1, 2, 3, 4, 5},
		ReceiveTime:  NewTimestamp(time.Unix(10, 0).UTC()),
		BatchID:      nil,
		IndexInBatch: nil,
		CommitmentID: &CommitmentID{
			BatchID:      MakeUint256(10),
			IndexInBatch: 2,
		},
		ErrorMessage: ref.String("some error message"),
	}

	bytes := transactionBase.Bytes()

	decodedTransactionBase := TransactionBase{}
	err := decodedTransactionBase.SetBytes(bytes)
	require.NoError(t, err)
	require.EqualValues(t, *transactionBase, decodedTransactionBase)
}
