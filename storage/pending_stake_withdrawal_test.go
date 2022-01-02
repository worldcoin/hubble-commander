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

func (s *PendingStakeWithdrawalTestSuite) TestGetPendingStakeWithdrawalsByFinalisationBlock_AddAndGet() {
	stake := models.PendingStakeWithdrawal{
		BatchID:           models.MakeUint256(uint64(100)),
		FinalisationBlock: uint32(125),
	}

	err := s.storage.AddPendingStakeWithdrawal(&stake)
	s.NoError(err)

	stStakes, err := s.storage.GetPendingStakeWithdrawalsByFinalisationBlock(stake.FinalisationBlock, stake.FinalisationBlock)
	s.NoError(err)
	s.Len(stStakes, 1)
	s.Equal(stake, *stStakes[0])
}

func (s *PendingStakeWithdrawalTestSuite) TestGetPendingStakeWithdrawalsByFinalisationBlock_NonexistentStake() {
	stakes, err := s.storage.GetPendingStakeWithdrawalsByFinalisationBlock(123, 129)
	s.NoError(err)
	s.Len(stakes, 0)
}

func TestPendingStakeWithdrawalTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(PendingStakeWithdrawalTestSuite))
}
