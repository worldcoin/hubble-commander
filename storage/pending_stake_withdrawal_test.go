package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PendingStakeWithdrawalTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *PendingStakeWithdrawalTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *PendingStakeWithdrawalTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *PendingStakeWithdrawalTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *PendingStakeWithdrawalTestSuite) TestAddPendingStakeWithdrawal_AddAndDelete() {
	stake := models.PendingStakeWithdrawal{
		BatchID:           models.MakeUint256(1),
		FinalisationBlock: 100,
	}
	err := s.storage.AddPendingStakeWithdrawal(&stake)
	s.NoError(err)

	err = s.storage.RemovePendingStakeWithdrawal(stake.BatchID)
	s.NoError(err)
}

func (s *PendingStakeWithdrawalTestSuite) TestRemovePendingStakeWithdrawal_NonexistentStake() {
	err := s.storage.RemovePendingStakeWithdrawal(models.MakeUint256(42))
	s.ErrorIs(err, NewNotFoundError("pending stake withdrawal"))
}

func (s *PendingStakeWithdrawalTestSuite) TestGetReadyStateWithdrawals_AddAndGet() {
	stakes := []models.PendingStakeWithdrawal{
		{
			BatchID:           models.MakeUint256(uint64(2)),
			FinalisationBlock: uint32(12),
		},
		{
			BatchID:           models.MakeUint256(uint64(1)),
			FinalisationBlock: uint32(10),
		},
	}

	for _, stake := range stakes {
		err := s.storage.AddPendingStakeWithdrawal(&stake)
		s.NoError(err)
	}

	expectedStake := stakes[1]
	stStakes, err := s.storage.GetReadyStateWithdrawals(expectedStake.FinalisationBlock)
	s.NoError(err)
	s.Len(stStakes, 1)
	s.Equal(expectedStake, stStakes[0])
}

func (s *PendingStakeWithdrawalTestSuite) TestGetReadyStateWithdrawals_NonexistentStake() {
	stakes, err := s.storage.GetReadyStateWithdrawals(129)
	s.NoError(err)
	s.Len(stakes, 0)
}

func TestPendingStakeWithdrawalTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(PendingStakeWithdrawalTestSuite))
}
