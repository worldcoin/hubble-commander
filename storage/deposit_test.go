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
		IncludedInCommitment: &models.CommitmentID{
			BatchID:      models.MakeUint256(9),
			IndexInBatch: 17,
		},
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

func TestDepositTestSuite(t *testing.T) {
	suite.Run(t, new(DepositTestSuite))
}
