package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage           *TestStorage
	txCommitment      *models.TxCommitment
	mmCommitment      *models.MMCommitment
	depositCommitment *models.DepositCommitment
}

func (s *CommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)

	s.txCommitment = &models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       1,
		CombinedSignature: models.MakeRandomSignature(),
	}

	s.mmCommitment = &models.MMCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
			Type:          batchtype.MassMigration,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       1,
		CombinedSignature: models.MakeRandomSignature(),
		Meta: &models.MassMigrationMeta{
			SpokeID:     1,
			TokenID:     models.MakeUint256(2),
			Amount:      models.MakeUint256(200),
			FeeReceiver: 1,
		},
		WithdrawRoot: utils.RandomHash(),
	}

	s.depositCommitment = &models.DepositCommitment{
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
				ID: models.DepositID{
					SubtreeID:    models.MakeUint256(1),
					DepositIndex: models.MakeUint256(0),
				},
				ToPubKeyID: 1,
				TokenID:    models.MakeUint256(1),
				L2Amount:   models.MakeUint256(1000),
			},
		},
	}
}

func (s *CommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *CommitmentTestSuite) TestAddCommitment_AddTxCommitmentAndRetrieve() {
	err := s.storage.AddCommitment(s.txCommitment)
	s.NoError(err)

	actual, err := s.storage.GetCommitment(&s.txCommitment.ID)
	s.NoError(err)
	s.Equal(s.txCommitment, actual.ToTxCommitment())
}

func (s *CommitmentTestSuite) TestAddCommitment_AddMMCommitmentAndRetrieve() {
	err := s.storage.AddCommitment(s.mmCommitment)
	s.NoError(err)

	actual, err := s.storage.GetCommitment(&s.mmCommitment.ID)
	s.NoError(err)
	s.Equal(s.mmCommitment, actual.ToMMCommitment())
}

func (s *CommitmentTestSuite) TestAddCommitment_AddDepositCommitmentAndRetrieve() {
	err := s.storage.AddCommitment(s.depositCommitment)
	s.NoError(err)

	actual, err := s.storage.GetCommitment(&s.depositCommitment.ID)
	s.NoError(err)
	s.Equal(s.depositCommitment, actual.ToDepositCommitment())
}

func (s *CommitmentTestSuite) TestAddCommitment_GenesisCommitment() {
	genesisCommitment := &models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
			Type: batchtype.Genesis,
		},
	}

	s.Panicsf(func() {
		_ = s.storage.AddCommitment(genesisCommitment)
	}, "invalid commitment type")
}

func (s *CommitmentTestSuite) TestGetCommitment_NonexistentCommitment() {
	res, err := s.storage.GetCommitment(&s.txCommitment.ID)
	s.ErrorIs(err, NewNotFoundError("commitment"))
	s.Nil(res)
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchID_TransfersAndC2TsBatch() {
	err := s.storage.AddCommitment(s.txCommitment)
	s.NoError(err)

	commitments, err := s.storage.GetCommitmentsByBatchID(s.txCommitment.ID.BatchID)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Equal(s.txCommitment, commitments[0].ToTxCommitment())
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchID_MMsBatch() {
	err := s.storage.AddCommitment(s.mmCommitment)
	s.NoError(err)

	commitments, err := s.storage.GetCommitmentsByBatchID(s.mmCommitment.ID.BatchID)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Equal(s.mmCommitment, commitments[0].ToMMCommitment())
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchID_DepositsBatch() {
	err := s.storage.AddCommitment(s.depositCommitment)
	s.NoError(err)

	commitments, err := s.storage.GetCommitmentsByBatchID(s.depositCommitment.ID.BatchID)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Equal(s.depositCommitment, commitments[0].ToDepositCommitment())
}

func (s *CommitmentTestSuite) TestUpdateCommitments_TxCommitment() {
	s.testUpdateCommitments(s.txCommitment)
}

func (s *CommitmentTestSuite) TestUpdateCommitments_MMCommitment() {
	s.testUpdateCommitments(s.mmCommitment)
}

func (s *CommitmentTestSuite) TestUpdateCommitments_DepositCommitment() {
	s.testUpdateCommitments(s.depositCommitment)
}

func (s *CommitmentTestSuite) TestUpdateCommitments_InvalidCommitmentType() {
	commitment := models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(10),
				IndexInBatch: 0,
			},
			Type: batchtype.Genesis,
		},
	}

	s.Panicsf(func() {
		_ = s.storage.UpdateCommitments([]models.Commitment{&commitment})
	}, "invalid commitment type")
}

func (s *CommitmentTestSuite) TestUpdateCommitments_NonexistentCommitment() {
	commitment := *s.txCommitment
	commitment.BodyHash = utils.NewRandomHash()
	err := s.storage.UpdateCommitments([]models.Commitment{&commitment})
	s.ErrorIs(err, NewNotFoundError("commitment"))
}

func (s *CommitmentTestSuite) testUpdateCommitments(commitment models.Commitment) {
	commitment.GetCommitmentBase().ID.IndexInBatch = uint8(0)

	err := s.storage.AddCommitment(commitment)
	s.NoError(err)

	commitment.GetCommitmentBase().PostStateRoot = utils.RandomHash()

	err = s.storage.UpdateCommitments([]models.Commitment{commitment})
	s.NoError(err)

	commitments, err := s.storage.GetCommitmentsByBatchID(commitment.GetCommitmentBase().ID.BatchID)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Equal(commitment, commitments[0])
}

func TestCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(CommitmentTestSuite))
}
