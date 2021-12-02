package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	depositSubTree = models.PendingDepositSubTree{
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

type DepositSubTreeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *DepositSubTreeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DepositSubTreeTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *DepositSubTreeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DepositSubTreeTestSuite) TestAddPendingDepositSubTree_AddAndRetrieve() {
	err := s.storage.AddPendingDepositSubTree(&depositSubTree)
	s.NoError(err)

	actual, err := s.storage.GetPendingDepositSubTree(depositSubTree.ID)
	s.NoError(err)
	s.Equal(depositSubTree, *actual)
}

func (s *DepositSubTreeTestSuite) TestGetPendingDepositSubTree_NonexistentTree() {
	_, err := s.storage.GetPendingDepositSubTree(depositSubTree.ID)
	s.ErrorIs(err, NewNotFoundError("deposit sub tree"))
	s.True(IsNotFoundError(err))
}

func (s *DepositSubTreeTestSuite) TestDeletePendingDepositSubTrees() {
	subTrees := []models.PendingDepositSubTree{
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
	for i := range subTrees {
		err := s.storage.AddPendingDepositSubTree(&subTrees[i])
		s.NoError(err)
	}

	err := s.storage.DeletePendingDepositSubTrees(subTrees[0].ID, subTrees[1].ID)
	s.NoError(err)

	for i := range subTrees {
		_, err = s.storage.GetPendingDepositSubTree(subTrees[i].ID)
		s.ErrorIs(err, NewNotFoundError("deposit sub tree"))
	}
}

func (s *DepositSubTreeTestSuite) TestDeletePendingDepositSubTrees_NonexistentTree() {
	err := s.storage.DeletePendingDepositSubTrees(models.MakeUint256(1))
	s.ErrorIs(err, NewNotFoundError("deposit sub tree"))
}

func (s *DepositSubTreeTestSuite) TestGetFirstPendingDepositSubTree() {
	err := s.storage.AddPendingDepositSubTree(&depositSubTree)
	s.NoError(err)

	secondSubTree := depositSubTree
	secondSubTree.ID = models.MakeUint256(1)
	err = s.storage.AddPendingDepositSubTree(&secondSubTree)
	s.NoError(err)

	subTree, err := s.storage.GetFirstPendingDepositSubTree()
	s.NoError(err)
	s.Equal(secondSubTree, *subTree)
}

func (s *DepositSubTreeTestSuite) TestGetFirstPendingDepositSubTree_NoPendingDepositSubTrees() {
	subTree, err := s.storage.GetFirstPendingDepositSubTree()
	s.ErrorIs(err, NewNotFoundError("deposit sub tree"))
	s.Nil(subTree)
}

func TestDepositSubTreeTestSuite(t *testing.T) {
	suite.Run(t, new(DepositSubTreeTestSuite))
}
