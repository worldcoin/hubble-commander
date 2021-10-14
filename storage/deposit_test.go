package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	exampleDeposit := s.addPendingDeposit(16, 32)

	deposits, err := s.storage.GetFirstPendingDeposits(1)
	s.NoError(err)
	s.Equal(exampleDeposit, deposits[0])
}

func (s *DepositTestSuite) TestRemovePendingDeposits() {
	deposits := []models.PendingDeposit{
		s.addPendingDeposit(123, 1),
		s.addPendingDeposit(582, 17),
	}

	err := s.storage.RemovePendingDeposits(deposits)
	s.NoError(err)

	_, err = s.storage.GetFirstPendingDeposits(2)
	s.True(IsNotFoundError(err))
}

func (s *DepositTestSuite) TestGetFirstPendingDeposits() {
	allDeposits := []models.PendingDeposit{
		s.addPendingDeposit(1, 0),
		s.addPendingDeposit(1, 2),
		s.addPendingDeposit(3, 7),
		s.addPendingDeposit(3, 12),
	}

	amount := 3
	pendingDeposits, err := s.storage.GetFirstPendingDeposits(amount)
	s.NoError(err)
	s.Len(pendingDeposits, amount)
	s.Equal(allDeposits[0], pendingDeposits[0])
	s.Equal(allDeposits[1], pendingDeposits[1])
	s.Equal(allDeposits[2], pendingDeposits[2])
}

func (s *DepositTestSuite) TestGetFirstPendingDeposits_NoDeposits() {
	deposits, err := s.storage.GetFirstPendingDeposits(1)
	s.True(IsNotFoundError(err))
	s.Nil(deposits)
}

func (s *DepositTestSuite) addPendingDeposit(blockNumber, logIndex uint32) models.PendingDeposit {
	deposit := models.PendingDeposit{
		ID: models.DepositID{
			BlockNumber: blockNumber,
			LogIndex:    logIndex,
		},
		ToPubKeyID: 4,
		TokenID:    models.MakeUint256(4),
		L2Amount:   models.MakeUint256(1024),
	}
	err := s.storage.AddPendingDeposit(&deposit)
	s.NoError(err)
	return deposit
}

func TestDepositTestSuite(t *testing.T) {
	suite.Run(t, new(DepositTestSuite))
}
