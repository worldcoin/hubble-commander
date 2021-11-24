package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
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
		SpokeID: models.MakeUint256(5),
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

func (s *TransferTestSuite) TestAddMassMigration_AddAndRetrieve() {
	err := s.storage.AddMassMigration(&massMigration)
	s.NoError(err)

	expected := massMigration

	res, err := s.storage.GetMassMigration(massMigration.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *TransferTestSuite) TestAddMassMigration_AddAndRetrieveIncludedMassMigration() {
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

func (s *TransferTestSuite) TestGetMassMigration_NonexistentMassMigration() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetMassMigration(hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(res)
}

func TestMassMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MassMigrationTestSuite))
}
