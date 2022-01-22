package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	mmCommitment = models.MMCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
			Type:          batchtype.MassMigration,
			PostStateRoot: utils.RandomHash(),
		},
		CombinedSignature: models.MakeRandomSignature(),
		BodyHash:          utils.NewRandomHash(),
		Meta: &models.MassMigrationMeta{
			SpokeID:     1,
			TokenID:     models.MakeUint256(2),
			Amount:      models.MakeUint256(3),
			FeeReceiver: 4,
		},
		WithdrawRoot: utils.RandomHash(),
	}
)

type MMCommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *MMCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MMCommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *MMCommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MMCommitmentTestSuite) TestAddMMCommitment_AddAndRetrieve() {
	err := s.storage.addMMCommitment(&mmCommitment)
	s.NoError(err)

	actual, err := s.storage.GetCommitment(&mmCommitment.ID)
	s.NoError(err)
	s.Equal(mmCommitment, *actual.ToMMCommitment())
}

func TestMMCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(MMCommitmentTestSuite))
}
