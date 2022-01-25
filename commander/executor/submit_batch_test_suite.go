package executor

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
)

var (
	baseCommitment = models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				Type:          batchtype.Transfer,
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       1,
			CombinedSignature: models.MakeRandomSignature(),
		},
		Transactions: utils.RandomBytes(24),
	}
)

type submitBatchTestSuite struct {
	testSuiteWithTxsContext
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

func getTxCommitments(count int, batchID models.Uint256, batchType batchtype.BatchType) []models.CommitmentWithTxs {
	commitments := make([]models.CommitmentWithTxs, 0, count)
	for i := 0; i < count; i++ {
		commitment := baseCommitment
		commitment.Type = batchType
		commitment.ID.BatchID = batchID
		commitment.ID.IndexInBatch = uint8(i)

		commitments = append(commitments, &commitment)
	}
	return commitments
}

func getMMCommitments(count int, batchID models.Uint256) []models.CommitmentWithTxs {
	commitments := make([]models.CommitmentWithTxs, 0, count)
	for i := 0; i < count; i++ {
		commitment := models.MMCommitmentWithTxs{
			MMCommitment: models.MMCommitment{
				CommitmentBase: models.CommitmentBase{
					Type:          batchtype.MassMigration,
					PostStateRoot: utils.RandomHash(),
				},
				Meta: &models.MassMigrationMeta{
					SpokeID:     1,
					TokenID:     models.MakeUint256(2),
					Amount:      models.MakeUint256(3),
					FeeReceiver: 1,
				},
				CombinedSignature: models.MakeRandomSignature(),
				WithdrawRoot:      utils.RandomHash(),
			},
			Transactions: utils.RandomBytes(8),
		}
		commitment.ID.BatchID = batchID
		commitment.ID.IndexInBatch = uint8(i)

		commitments = append(commitments, &commitment)
	}
	return commitments
}
