package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
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

	deposits, err := s.storage.GetFirstPendingDeposits(1)
	s.NoError(err)
	s.Equal(deposit, deposits[0])
}

func (s *DepositTestSuite) TestRemovePendingDeposits() {
	deposits := []models.PendingDeposit{
		{
			ID: models.DepositID{
				BlockNumber: 123,
				LogIndex:    1,
			},
		},
		{
			ID: models.DepositID{
				BlockNumber: 582,
				LogIndex:    17,
			},
		},
	}

	err := s.storage.AddPendingDeposit(&deposits[0])
	s.NoError(err)
	err = s.storage.AddPendingDeposit(&deposits[1])
	s.NoError(err)

	err = s.storage.RemovePendingDeposits(deposits)
	s.NoError(err)

	_, err = s.storage.GetFirstPendingDeposits(2)
	s.ErrorIs(err, db.ErrIteratorFinished)
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
