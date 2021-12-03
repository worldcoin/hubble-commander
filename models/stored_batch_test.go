package models

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestStoredBatch_Bytes(t *testing.T) {
	batch := &StoredBatch{
		ID:                MakeUint256(10),
		BType:             batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              &common.Hash{8, 6, 4},
		FinalisationBlock: ref.Uint32(25),
		AccountTreeRoot:   &common.Hash{4, 5, 6},
		PrevStateRoot:     &common.Hash{7, 8, 9},
		SubmissionTime:    NewTimestamp(time.Unix(10, 11).UTC()),
	}

	bytes := batch.Bytes()

	decodedStoredBatch := StoredBatch{}
	err := decodedStoredBatch.SetBytes(bytes)
	require.NoError(t, err)
	require.EqualValues(t, *batch, decodedStoredBatch)
}
