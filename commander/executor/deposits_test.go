package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DepositsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *st.TestStorage
	transactionExecutor *TransactionExecutor
}

func (s *DepositsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DepositsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, &eth.Client{}, nil, context.Background())
}

func (s *DepositsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DepositsTestSuite) TestGetVacancyProof_EmptyTree() {
	stateID, err := s.transactionExecutor.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.transactionExecutor.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 0)
	s.Len(vacancyProof.Witness, 30)
}

func (s *DepositsTestSuite) TestGetVacancyProof_SingleLeafSet() {
	_, err := s.transactionExecutor.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)

	stateID, err := s.transactionExecutor.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.transactionExecutor.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 1)
	s.Len(vacancyProof.Witness, 30)
}

func (s *DepositsTestSuite) TestGetVacancyProof_TwoLeavesSet() {
	_, err := s.transactionExecutor.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)
	_, err = s.transactionExecutor.storage.StateTree.Set(4, &models.UserState{})
	s.NoError(err)

	stateID, err := s.transactionExecutor.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.transactionExecutor.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 2)
	s.Len(vacancyProof.Witness, 30)
}

func (s *DepositsTestSuite) TestGetVacancyProof_ProducesCorrectWitness() {
	userState := &models.UserState{}
	leafWitness, err := s.transactionExecutor.storage.StateTree.Set(0, userState)
	s.NoError(err)

	leaf, err := st.NewStateLeaf(0, userState)
	s.NoError(err)

	currentHash := leaf.DataHash
	for i := range leafWitness[:len(leafWitness)-2] {
		currentHash = utils.HashTwo(currentHash, leafWitness[i])
	}
	firstWitness := currentHash
	secondWitness := merkletree.GetZeroHash(31)

	stateID, err := s.transactionExecutor.storage.StateTree.NextVacantSubtree(30)
	s.NoError(err)

	vacancyProof, err := s.transactionExecutor.GetVacancyProof(*stateID, 30)
	s.NoError(err)

	s.Len(vacancyProof.Witness, 2)
	s.Equal(vacancyProof.Witness[0], firstWitness)
	s.Equal(vacancyProof.Witness[1], secondWitness)
}

func TestDepositsTestSuite(t *testing.T) {
	suite.Run(t, new(DepositsTestSuite))
}
