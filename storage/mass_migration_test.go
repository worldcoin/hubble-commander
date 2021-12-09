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

func (s *MassMigrationTestSuite) TestGetMassMigrationsByCommitmentID() {
	massMigration1 := massMigration
	massMigration1.CommitmentID = &txCommitment.ID

	err := s.storage.AddMassMigration(&massMigration1)
	s.NoError(err)

	massMigrations, err := s.storage.GetMassMigrationsByCommitmentID(txCommitment.ID)
	s.NoError(err)
	s.Len(massMigrations, 1)
}

func (s *MassMigrationTestSuite) TestGetMassMigrationsByCommitmentID_NoTransactions() {
	massMigrations, err := s.storage.GetMassMigrationsByCommitmentID(txCommitment.ID)
	s.NoError(err)
	s.Len(massMigrations, 0)
}

func (s *MassMigrationTestSuite) TestGetMassMigrationsByCommitmentID_NoMassMigrationsButSomeTransfers() {
	transfer1 := transfer
	transfer1.CommitmentID = &txCommitment.ID
	err := s.storage.AddTransfer(&transfer1)
	s.NoError(err)

	massMigrations, err := s.storage.GetMassMigrationsByCommitmentID(txCommitment.ID)
	s.NoError(err)
	s.Len(massMigrations, 0)
}

func TestMassMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MassMigrationTestSuite))
}
