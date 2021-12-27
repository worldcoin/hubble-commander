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

func (s *PendingStakeWithdrawalTestSuite) TestRemovePendingStakeWithdrawal_NonexistentStoke() {
	err := s.storage.RemovePendingStakeWithdrawal(models.MakeUint256(42))
	s.ErrorIs(err, NewNotFoundError("pending stake withdrawal"))
}

func TestPendingStakeWithdrawalTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(PendingStakeWithdrawalTestSuite))
}
