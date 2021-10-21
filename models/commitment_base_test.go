package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestCommitmentBase_Bytes(t *testing.T) {
	base := CommitmentBase{
		ID: CommitmentID{
			BatchID:      MakeUint256(1),
			IndexInBatch: 2,
		},
		Type:          batchtype.Create2Transfer,
		PostStateRoot: utils.RandomHash(),
	}

	bytes := base.Bytes()

	var decodedBase CommitmentBase
	err := decodedBase.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, base, decodedBase)
}
