package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	depositSubtree = models.PendingDepositSubtree{
		ID:   models.MakeUint256(932),
		Root: utils.RandomHash(),
		Deposits: []models.PendingDeposit{
			{
				ID: models.DepositID{
					SubtreeID:    models.MakeUint256(932),
					DepositIndex: models.MakeUint256(0),
				},
				ToPubKeyID: 3,
				TokenID:    models.MakeUint256(4),
				L2Amount:   models.MakeUint256(500),
			},
			{
				ID: models.DepositID{
					SubtreeID:    models.MakeUint256(932),
					DepositIndex: models.MakeUint256(1),
				},
				ToPubKeyID: 8,
				TokenID:    models.MakeUint256(9),
				L2Amount:   models.MakeUint256(1000),
			},
		},
	}
)

type DepositSubtreeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *DepositSubtreeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DepositSubtreeTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *DepositSubtreeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DepositSubtreeTestSuite) TestAddPendingDepositSubtree_AddAndRetrieve() {
	err := s.storage.AddPendingDepositSubtree(&depositSubtree)
	s.NoError(err)

	actual, err := s.storage.GetPendingDepositSubtree(depositSubtree.ID)
	s.NoError(err)
	s.Equal(depositSubtree, *actual)
}

func (s *DepositSubtreeTestSuite) TestGetPendingDepositSubtree_NonexistentTree() {
	_, err := s.storage.GetPendingDepositSubtree(depositSubtree.ID)
	s.ErrorIs(err, NewNotFoundError("deposit sub tree"))
	s.True(IsNotFoundError(err))
}

func (s *DepositSubtreeTestSuite) TestDeletePendingDepositSubtrees() {
	subtrees := []models.PendingDepositSubtree{
		{
			ID:   models.MakeUint256(1),
			Root: utils.RandomHash(),
			Deposits: []models.PendingDeposit{
				{
					ID: models.DepositID{
						SubtreeID:    models.MakeUint256(1),
						DepositIndex: models.MakeUint256(0),
					},
					ToPubKeyID: 3,
					TokenID:    models.MakeUint256(4),
					L2Amount:   models.MakeUint256(500),
				},
			},
		},
		{
			ID:   models.MakeUint256(4),
			Root: utils.RandomHash(),
			Deposits: []models.PendingDeposit{
				{
					ID: models.DepositID{
						SubtreeID:    models.MakeUint256(4),
						DepositIndex: models.MakeUint256(0),
					},
					ToPubKeyID: 8,
					TokenID:    models.MakeUint256(9),
					L2Amount:   models.MakeUint256(1000),
				},
			},
		},
	}
	for i := range subtrees {
		err := s.storage.AddPendingDepositSubtree(&subtrees[i])
		s.NoError(err)
	}

	err := s.storage.DeletePendingDepositSubtrees(subtrees[0].ID, subtrees[1].ID)
	s.NoError(err)

	for i := range subtrees {
		_, err = s.storage.GetPendingDepositSubtree(subtrees[i].ID)
		s.ErrorIs(err, NewNotFoundError("deposit sub tree"))
	}
}

func (s *DepositSubtreeTestSuite) TestDeletePendingDepositSubtrees_NonexistentTree() {
	err := s.storage.DeletePendingDepositSubtrees(models.MakeUint256(1))
	s.ErrorIs(err, NewNotFoundError("deposit sub tree"))
}

func (s *DepositSubtreeTestSuite) TestGetFirstPendingDepositSubtree() {
	err := s.storage.AddPendingDepositSubtree(&depositSubtree)
	s.NoError(err)

	secondSubtree := depositSubtree
	secondSubtree.ID = models.MakeUint256(1)
	err = s.storage.AddPendingDepositSubtree(&secondSubtree)
	s.NoError(err)

	subtree, err := s.storage.GetFirstPendingDepositSubtree()
	s.NoError(err)
	s.Equal(secondSubtree, *subtree)
}

func (s *DepositSubtreeTestSuite) TestGetFirstPendingDepositSubtree_NoPendingDepositSubtrees() {
	subtree, err := s.storage.GetFirstPendingDepositSubtree()
	s.ErrorIs(err, NewNotFoundError("deposit sub tree"))
	s.Nil(subtree)
}

func TestDepositSubtreeTestSuite(t *testing.T) {
	suite.Run(t, new(DepositSubtreeTestSuite))
}
