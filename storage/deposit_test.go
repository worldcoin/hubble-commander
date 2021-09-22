package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	deposit = models.PendingDeposit{
		ID: models.DepositID{
			BlockNumber: 16,
			LogIndex:    32,
		},
		ToPubKeyID: 4,
		TokenID:    models.MakeUint256(4),
		L2Amount:   models.MakeUint256(1024),
	}
)

type DepositTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *DepositTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DepositTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *DepositTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DepositTestSuite) TestAddPendingDeposit_AddAndRetrieve() {
	err := s.storage.AddPendingDeposit(&deposit)
	s.NoError(err)

	actual, err := s.storage.GetPendingDeposit(&deposit.ID)
	s.NoError(err)
	s.Equal(deposit, *actual)
}

func (s *DepositTestSuite) TestGetPendingDeposit_NotFound() {
	_, err := s.storage.GetPendingDeposit(&deposit.ID)
	s.ErrorIs(err, NewNotFoundError("pending deposit"))
	s.True(IsNotFoundError(err))
}

func (s *DepositTestSuite) TestGetFirstPendingDeposits() {
	allDeposits := []models.PendingDeposit{
		{
			ID: models.DepositID{
				BlockNumber: 1,
				LogIndex:    0,
			},
			ToPubKeyID: 4,
			TokenID:    models.MakeUint256(4),
			L2Amount:   models.MakeUint256(1024),
		},
		{
			ID: models.DepositID{
				BlockNumber: 1,
				LogIndex:    2,
			},
			ToPubKeyID: 4,
			TokenID:    models.MakeUint256(4),
			L2Amount:   models.MakeUint256(1024),
		},
		{
			ID: models.DepositID{
				BlockNumber: 3,
				LogIndex:    7,
			},
			ToPubKeyID: 4,
			TokenID:    models.MakeUint256(4),
			L2Amount:   models.MakeUint256(1024),
		},
		{
			ID: models.DepositID{
				BlockNumber: 3,
				LogIndex:    12,
			},
			ToPubKeyID: 4,
			TokenID:    models.MakeUint256(4),
			L2Amount:   models.MakeUint256(1024),
		},
	}

	for i := range allDeposits {
		err := s.storage.AddPendingDeposit(&allDeposits[i])
		s.NoError(err)
	}

	amount := 3
	pendingDeposits, err := s.storage.GetFirstPendingDeposits(amount)
	s.NoError(err)
	s.Len(pendingDeposits, amount)
	s.Equal(allDeposits[0], pendingDeposits[0])
	s.Equal(allDeposits[1], pendingDeposits[1])
	s.Equal(allDeposits[2], pendingDeposits[2])
}

func TestDepositTestSuite(t *testing.T) {
	suite.Run(t, new(DepositTestSuite))
}
