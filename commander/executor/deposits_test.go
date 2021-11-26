package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DepositsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage        *st.TestStorage
	client         *eth.TestClient
	depositsCtx    *DepositsContext
	depositSubtree models.PendingDepositSubTree
}

func (s *DepositsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.depositSubtree = models.PendingDepositSubTree{
		ID:       models.MakeUint256(1),
		Root:     utils.RandomHash(),
		Deposits: testutils.GetFourDeposits(),
	}
}

func (s *DepositsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	executionCtx := NewTestExecutionContext(s.storage.Storage, s.client.Client, nil)
	s.depositsCtx = NewTestDepositsContext(executionCtx)
}

func (s *DepositsTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DepositsTestSuite) TestCreateCommitment_AddsCommitment() {
	err := s.storage.AddPendingDepositSubTree(&s.depositSubtree)
	s.NoError(err)

	batchID := models.MakeUint256(1)
	_, err = s.depositsCtx.createCommitment(batchID)
	s.NoError(err)

	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	commitment, err := s.storage.GetDepositCommitment(&models.CommitmentID{
		BatchID:      batchID,
		IndexInBatch: 0,
	})
	s.NoError(err)
	s.Equal(*root, commitment.PostStateRoot)
	s.Equal(s.depositSubtree.ID, commitment.SubTreeID)
	s.Equal(s.depositSubtree.Root, commitment.SubTreeRoot)
	s.Equal(s.depositSubtree.Deposits, commitment.Deposits)
}

func (s *DepositsTestSuite) TestCreateCommitment_NotEnoughDeposits() {
	vacancyProof, err := s.depositsCtx.createCommitment(models.MakeUint256(1))
	s.ErrorIs(err, ErrNotEnoughDeposits)
	s.Nil(vacancyProof)
}

func (s *DepositsTestSuite) TestExecuteDeposits_SetsUserStates() {
	err := s.storage.AddPendingDepositSubTree(&s.depositSubtree)
	s.NoError(err)

	_, err = s.depositsCtx.executeDeposits(&s.depositSubtree)
	s.NoError(err)

	for i := range s.depositSubtree.Deposits {
		stateLeaf, err := s.storage.StateTree.Leaf(uint32(i))
		s.NoError(err)
		s.Equal(s.depositSubtree.Deposits[i].L2Amount, stateLeaf.Balance)
	}
}

func (s *DepositsTestSuite) TestExecuteDeposits_RemovesDepositSubtree() {
	err := s.storage.AddPendingDepositSubTree(&s.depositSubtree)
	s.NoError(err)

	_, err = s.depositsCtx.executeDeposits(&s.depositSubtree)
	s.NoError(err)

	subtree, err := s.storage.GetPendingDepositSubTree(s.depositSubtree.ID)
	s.True(st.IsNotFoundError(err))
	s.Nil(subtree)
}

func (s *DepositsTestSuite) TestExecuteDeposits_ReturnsCorrectVacancyProof() {
	_, err := s.depositsCtx.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)

	err = s.storage.AddPendingDepositSubTree(&s.depositSubtree)
	s.NoError(err)

	vacancyProof, err := s.depositsCtx.executeDeposits(&s.depositSubtree)
	s.NoError(err)
	s.EqualValues(1, vacancyProof.PathAtDepth)
}

func TestDepositsTestSuite(t *testing.T) {
	suite.Run(t, new(DepositsTestSuite))
}
