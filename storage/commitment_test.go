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

func (s *CommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
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

	s.Panicsf(
		func() {
			err := s.storage.AddCommitment(genesisCommitment)
			s.NoError(err)
		},
		"invalid commitment type",
	)
}

func (s *CommitmentTestSuite) TestGetCommitment_NonexistentCommitment() {
	res, err := s.storage.GetCommitment(&s.txCommitment.ID)
	s.ErrorIs(err, NewNotFoundError("commitment"))
	s.Nil(res)
}

func TestCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(CommitmentTestSuite))
}
