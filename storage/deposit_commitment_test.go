package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DepositCommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage    *TestStorage
	commitment models.DepositCommitment
}

func (s *DepositCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.commitment = models.DepositCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
			Type:          batchtype.Deposit,
			PostStateRoot: utils.RandomHash(),
		},
		SubtreeID:   models.MakeUint256(1),
		SubtreeRoot: utils.RandomHash(),
		Deposits: []models.PendingDeposit{
			{
				ID:         models.DepositID{},
				ToPubKeyID: 1,
				TokenID:    models.MakeUint256(1),
				L2Amount:   models.MakeUint256(10),
			},
		},
	}
}

func (s *DepositCommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *DepositCommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DepositCommitmentTestSuite) TestAddDepositCommitment_AddAndRetrieve() {
	err := s.storage.addDepositCommitment(&s.commitment)
	s.NoError(err)

	actual, err := s.storage.GetCommitment(&s.commitment.ID)
	s.NoError(err)
	s.Equal(s.commitment, *actual.ToDepositCommitment())
}

func TestDepositCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(DepositCommitmentTestSuite))
}
