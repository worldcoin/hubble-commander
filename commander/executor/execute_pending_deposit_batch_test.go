package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ExecutePendingDepositBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	storage        *st.TestStorage
	client         *eth.TestClient
	depositsCtx    *DepositsContext
	depositSubtree models.PendingDepositSubtree
	pendingBatch   models.PendingBatch
}

func (s *ExecutePendingDepositBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.depositSubtree = models.PendingDepositSubtree{
		ID:       models.MakeUint256(1),
		Root:     utils.RandomHash(),
		Deposits: testutils.GetFourDeposits(),
	}
}

func (s *ExecutePendingDepositBatchTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	executionCtx := NewTestExecutionContext(s.storage.Storage, s.client.Client, nil)
	s.depositsCtx = NewTestDepositsContext(executionCtx)

	s.pendingBatch = models.PendingBatch{
		ID:              models.MakeUint256(1),
		Type:            batchtype.Deposit,
		TransactionHash: utils.RandomHash(),
		PrevStateRoot:   utils.RandomHash(),
		Commitments: []models.PendingCommitment{
			{
				Commitment: &models.DepositCommitment{
					CommitmentBase: models.CommitmentBase{
						ID: models.CommitmentID{
							BatchID:      models.MakeUint256(1),
							IndexInBatch: 0,
						},
						Type:          batchtype.Deposit,
						PostStateRoot: utils.RandomHash(),
					},
					SubtreeID:   s.depositSubtree.ID,
					SubtreeRoot: s.depositSubtree.Root,
					Deposits:    s.depositSubtree.Deposits,
				},
			},
		},
	}
}

func (s *ExecutePendingDepositBatchTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ExecutePendingDepositBatchTestSuite) TestExecutePendingBatch_AddsBatch() {
	err := s.storage.AddPendingDepositSubtree(&s.depositSubtree)
	s.NoError(err)

	err = s.depositsCtx.ExecutePendingBatch(&s.pendingBatch)
	s.NoError(err)

	expectedBatch := models.Batch{
		ID:              s.pendingBatch.ID,
		Type:            s.pendingBatch.Type,
		TransactionHash: s.pendingBatch.TransactionHash,
		PrevStateRoot:   s.pendingBatch.PrevStateRoot,
	}

	batch, err := s.storage.GetBatch(s.pendingBatch.ID)
	s.NoError(err)
	s.Equal(expectedBatch, *batch)

	commitments, err := s.storage.GetCommitmentsByBatchID(s.pendingBatch.ID)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Equal(s.pendingBatch.Commitments[0].Commitment, commitments[0])

	_, err = s.storage.GetPendingDepositSubtree(s.depositSubtree.ID)
	s.True(st.IsNotFoundError(err))
}

func TestExecutePendingDepositBatchTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutePendingDepositBatchTestSuite))
}
