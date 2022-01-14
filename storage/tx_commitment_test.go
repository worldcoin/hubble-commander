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
	txCommitment = models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeRandomSignature(),
		BodyHash:          utils.NewRandomHash(),
	}
)

type TxCommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *TxCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxCommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *TxCommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TxCommitmentTestSuite) TestAddTxCommitment_AddAndRetrieve() {
	err := s.storage.addTxCommitment(&txCommitment)
	s.NoError(err)

	actual, err := s.storage.GetCommitment(&txCommitment.ID)
	s.NoError(err)
	s.Equal(txCommitment, *actual.ToTxCommitment())
}

func TestTxCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(TxCommitmentTestSuite))
}
