package prover

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VacancyProofTestSuite struct {
	*require.Assertions
	suite.Suite
	storage        *st.TestStorage
	proverCtx      *Context
	depositSubtree models.PendingDepositSubTree
}

func (s *VacancyProofTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.depositSubtree = models.PendingDepositSubTree{
		ID:   models.MakeUint256(1),
		Root: utils.RandomHash(),
		Deposits: []models.PendingDeposit{
			{
				ID: models.DepositID{
					SubtreeID:    models.MakeUint256(1),
					DepositIndex: models.MakeUint256(0),
				},
				ToPubKeyID: 1,
				TokenID:    models.MakeUint256(0),
				L2Amount:   models.MakeUint256(50),
			},
			{
				ID: models.DepositID{
					SubtreeID:    models.MakeUint256(1),
					DepositIndex: models.MakeUint256(1),
				},
				ToPubKeyID: 1,
				TokenID:    models.MakeUint256(0),
				L2Amount:   models.MakeUint256(50),
			},
		},
	}
}

func (s *VacancyProofTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.proverCtx = NewContext(s.storage.Storage)
}

func (s *VacancyProofTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *VacancyProofTestSuite) TestGetVacancyProof_EmptyTree() {
	stateID, err := s.proverCtx.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.proverCtx.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 0)
	s.Len(vacancyProof.Witness, 30)
}

func (s *VacancyProofTestSuite) TestGetVacancyProof_SingleLeafSet() {
	_, err := s.proverCtx.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)

	stateID, err := s.proverCtx.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.proverCtx.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 1)
	s.Len(vacancyProof.Witness, 30)
}

func (s *VacancyProofTestSuite) TestGetVacancyProof_TwoLeavesSet() {
	_, err := s.proverCtx.storage.StateTree.Set(0, &models.UserState{})
	s.NoError(err)
	_, err = s.proverCtx.storage.StateTree.Set(4, &models.UserState{})
	s.NoError(err)

	stateID, err := s.proverCtx.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.proverCtx.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 2)
	s.Len(vacancyProof.Witness, 30)
}

func (s *VacancyProofTestSuite) TestGetVacancyProof_ProducesCorrectWitness() {
	userState := &models.UserState{}
	leafWitness, err := s.proverCtx.storage.StateTree.Set(0, userState)
	s.NoError(err)

	leaf, err := st.NewStateLeaf(0, userState)
	s.NoError(err)

	currentHash := leaf.DataHash
	for i := range leafWitness[:len(leafWitness)-2] {
		currentHash = utils.HashTwo(currentHash, leafWitness[i])
	}
	firstWitness := currentHash
	secondWitness := merkletree.GetZeroHash(31)

	stateID, err := s.proverCtx.storage.StateTree.NextVacantSubtree(30)
	s.NoError(err)

	vacancyProof, err := s.proverCtx.GetVacancyProof(*stateID, 30)
	s.NoError(err)

	s.Len(vacancyProof.Witness, 2)
	s.Equal(vacancyProof.Witness[0], firstWitness)
	s.Equal(vacancyProof.Witness[1], secondWitness)
}

func TestVacancyProofTestSuite(t *testing.T) {
	suite.Run(t, new(VacancyProofTestSuite))
}
