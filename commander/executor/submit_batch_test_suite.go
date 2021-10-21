package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
)

var (
	baseCommitment = models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			Type:          batchtype.Transfer,
			PostStateRoot: utils.RandomHash(),
		},
		Transactions:      utils.RandomBytes(24),
		FeeReceiver:       1,
		CombinedSignature: models.MakeRandomSignature(),
	}
)

type submitBatchTestSuite struct {
	testSuiteWithRollupContext
}

func (s *submitBatchTestSuite) setupUser() {
	userState := models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	}

	_, err := s.storage.StateTree.Set(1, &userState)
	s.NoError(err)
}

func getCommitments(count int, batchID models.Uint256) []models.TxCommitment {
	commitments := make([]models.TxCommitment, 0, count)
	for i := 0; i < count; i++ {
		commitment := baseCommitment
		commitment.ID.BatchID = batchID
		commitment.ID.IndexInBatch = uint8(i)

		commitments = append(commitments, commitment)
	}
	return commitments
}
