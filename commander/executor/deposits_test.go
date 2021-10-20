package executor

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DepositsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage        *st.TestStorage
	client         *eth.TestClient
	executionCtx   *ExecutionContext
	depositSubtree models.PendingDepositSubTree
}

func (s *DepositsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.depositSubtree = models.PendingDepositSubTree{
		ID:   models.MakeUint256(1),
		Root: utils.RandomHash(),
		Deposits: []models.PendingDeposit{
			{
				ID:         models.DepositID{BlockNumber: 1, LogIndex: 0},
				ToPubKeyID: 1,
				TokenID:    models.MakeUint256(0),
				L2Amount:   models.MakeUint256(50),
			},
			{
				ID:         models.DepositID{BlockNumber: 1, LogIndex: 1},
				ToPubKeyID: 1,
				TokenID:    models.MakeUint256(0),
				L2Amount:   models.MakeUint256(50),
			},
		},
	}
}

func (s *DepositsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.executionCtx = NewTestExecutionContext(s.storage.Storage, s.client.Client, nil)
}

func (s *DepositsTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DepositsTestSuite) TestGetVacancyProof_EmptyTree() {
	stateID, err := s.executionCtx.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.executionCtx.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 0)
	s.Len(vacancyProof.Witness, 30)
}

func (s *DepositsTestSuite) TestGetVacancyProof_SingleLeafSet() {
	_, err := s.executionCtx.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)

	stateID, err := s.executionCtx.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.executionCtx.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 1)
	s.Len(vacancyProof.Witness, 30)
}

func (s *DepositsTestSuite) TestGetVacancyProof_TwoLeavesSet() {
	_, err := s.executionCtx.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)
	_, err = s.executionCtx.storage.StateTree.Set(4, &models.UserState{})
	s.NoError(err)

	stateID, err := s.executionCtx.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.executionCtx.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 2)
	s.Len(vacancyProof.Witness, 30)
}

func (s *DepositsTestSuite) TestGetVacancyProof_ProducesCorrectWitness() {
	userState := &models.UserState{}
	leafWitness, err := s.executionCtx.storage.StateTree.Set(0, userState)
	s.NoError(err)

	leaf, err := st.NewStateLeaf(0, userState)
	s.NoError(err)

	currentHash := leaf.DataHash
	for i := range leafWitness[:len(leafWitness)-2] {
		currentHash = utils.HashTwo(currentHash, leafWitness[i])
	}
	firstWitness := currentHash
	secondWitness := merkletree.GetZeroHash(31)

	stateID, err := s.executionCtx.storage.StateTree.NextVacantSubtree(30)
	s.NoError(err)

	vacancyProof, err := s.executionCtx.GetVacancyProof(*stateID, 30)
	s.NoError(err)

	s.Len(vacancyProof.Witness, 2)
	s.Equal(vacancyProof.Witness[0], firstWitness)
	s.Equal(vacancyProof.Witness[1], secondWitness)
}

//TestSubmitBatch_SubmitsCommitmentsOnChain
//TestSubmitBatch_StoresPendingBatchRecord
//TestSubmitBatch_AddsCommitments

func (s *DepositsTestSuite) TestExecuteDeposits_SetsUserStates() {
	err := s.storage.AddPendingDepositSubTree(&s.depositSubtree)
	s.NoError(err)

	_, err = s.executionCtx.ExecuteDeposits()
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

	_, err = s.executionCtx.ExecuteDeposits()
	s.NoError(err)

	subtree, err := s.storage.GetPendingDepositSubTree(s.depositSubtree.ID)
	s.True(st.IsNotFoundError(err))
	s.Nil(subtree)
}

func (s *DepositsTestSuite) TestExecuteDeposits_ReturnsCorrectVacancyProof() {
	_, err := s.executionCtx.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)

	err = s.storage.AddPendingDepositSubTree(&s.depositSubtree)
	s.NoError(err)

	vacancyProof, err := s.executionCtx.ExecuteDeposits()
	s.NoError(err)
	s.EqualValues(1, vacancyProof.PathAtDepth)
}

func (s *DepositsTestSuite) TestSubmitDepositBatch_SubmitsBatchOnChain() {
	s.T().SkipNow()
	s.addGenesisBatch()

	pendingBatch, err := s.executionCtx.NewPendingBatch(batchtype.Deposit)
	s.NoError(err)

	_, vacancyProof, err := s.executionCtx.getDepositSubtreeVacancyProof()
	s.NoError(err)

	err = s.executionCtx.SubmitDepositBatch(pendingBatch, vacancyProof)
	s.NoError(err)

	s.client.Commit()

	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *DepositsTestSuite) addGenesisBatch() {
	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	batch, err := s.client.GetBatch(models.NewUint256(0))
	s.NoError(err)

	batch.PrevStateRoot = root
	err = s.storage.AddBatch(batch)
	s.NoError(err)
}

func TestDepositsTestSuite(t *testing.T) {
	suite.Run(t, new(DepositsTestSuite))
}
