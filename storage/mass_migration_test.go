package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	massMigration = models.MassMigration{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.MassMigration,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
		},
		SpokeID: 5,
	}
)

type MassMigrationTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *MassMigrationTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MassMigrationTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *MassMigrationTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MassMigrationTestSuite) TestAddMassMigration_AddAndRetrieve() {
	err := s.storage.AddMassMigration(&massMigration)
	s.NoError(err)

	expected := massMigration

	res, err := s.storage.GetMassMigration(massMigration.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *MassMigrationTestSuite) TestAddMassMigration_AddAndRetrieveIncludedMassMigration() {
	includedMassMigration := massMigration
	includedMassMigration.CommitmentID = &models.CommitmentID{
		BatchID:      models.MakeUint256(3),
		IndexInBatch: 1,
	}
	err := s.storage.AddMassMigration(&includedMassMigration)
	s.NoError(err)

	res, err := s.storage.GetMassMigration(massMigration.Hash)
	s.NoError(err)
	s.Equal(includedMassMigration, *res)
}

func (s *MassMigrationTestSuite) TestBatchAddMassMigration() {
	txs := make([]models.MassMigration, 2)
	txs[0] = massMigration
	txs[0].Hash = utils.RandomHash()
	txs[1] = massMigration
	txs[1].Hash = utils.RandomHash()

	err := s.storage.BatchAddMassMigration(txs)
	s.NoError(err)

	massMigration, err := s.storage.GetMassMigration(txs[0].Hash)
	s.NoError(err)
	s.Equal(txs[0], *massMigration)
	massMigration, err = s.storage.GetMassMigration(txs[1].Hash)
	s.NoError(err)
	s.Equal(txs[1], *massMigration)
}

func (s *MassMigrationTestSuite) TestBatchAddMassMigration_NoMassMigrations() {
	err := s.storage.BatchAddMassMigration([]models.MassMigration{})
	s.ErrorIs(err, ErrNoRowsAffected)
}

func (s *MassMigrationTestSuite) TestGetMassMigration_NonexistentMassMigration() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetMassMigration(hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(res)
}

func (s *MassMigrationTestSuite) TestGetPendingMassMigrations() {
	massMigrations := make([]models.MassMigration, 4)
	for i := range massMigrations {
		massMigrations[i] = massMigration
		massMigrations[i].Hash = utils.RandomHash()
	}
	massMigrations[2].CommitmentID = &models.CommitmentID{BatchID: models.MakeUint256(3)}
	massMigrations[3].ErrorMessage = ref.String("A very boring error message")

	err := s.storage.BatchAddMassMigration(massMigrations)
	s.NoError(err)

	res, err := s.storage.GetPendingMassMigrations()
	s.NoError(err)

	s.Len(res, 2)
	s.Contains(res, massMigrations[0])
	s.Contains(res, massMigrations[1])
}

func (s *MassMigrationTestSuite) TestGetPendingMassMigration_OrdersMassMigrationsByNonceAndTxHashAscending() {
	massMigration.Nonce = models.MakeUint256(1)
	massMigration.Hash = utils.RandomHash()
	massMigration2 := massMigration
	massMigration2.Nonce = models.MakeUint256(4)
	massMigration2.Hash = utils.RandomHash()
	massMigration3 := massMigration
	massMigration3.Nonce = models.MakeUint256(7)
	massMigration3.Hash = utils.RandomHash()
	massMigration4 := massMigration
	massMigration4.Nonce = models.MakeUint256(5)
	massMigration4.Hash = common.Hash{66, 66, 66, 66}
	massMigration5 := massMigration
	massMigration5.Nonce = models.MakeUint256(5)
	massMigration5.Hash = common.Hash{65, 65, 65, 65}

	massMigrations := []models.MassMigration{
		massMigration,
		massMigration2,
		massMigration3,
		massMigration4,
		massMigration5,
	}

	err := s.storage.BatchAddMassMigration(massMigrations)
	s.NoError(err)

	res, err := s.storage.GetPendingMassMigrations()
	s.NoError(err)

	s.Equal(models.MassMigrationArray{
		massMigration,
		massMigration2,
		massMigration5,
		massMigration4,
		massMigration3,
	}, res)
}

func (s *MassMigrationTestSuite) TestMarkMassMigrationsAsIncluded() {
	txs := make([]models.MassMigration, 2)
	for i := 0; i < len(txs); i++ {
		txs[i] = massMigration
		txs[i].Hash = utils.RandomHash()
		err := s.storage.AddMassMigration(&txs[i])
		s.NoError(err)
	}

	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 1,
	}
	err := s.storage.MarkMassMigrationsAsIncluded(txs, &commitmentID)
	s.NoError(err)

	for i := range txs {
		tx, err := s.storage.GetMassMigration(txs[i].Hash)
		s.NoError(err)
		s.Equal(commitmentID, *tx.CommitmentID)
	}
}

func TestMassMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MassMigrationTestSuite))
}
