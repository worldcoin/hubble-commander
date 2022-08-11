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

	account := models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{3, 4, 5},
	}
	err = s.storage.AccountTree.SetSingle(&account)
	s.NoError(err)

	leaf := models.StateLeaf{
		StateID: 1,
		UserState: models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(2000),
		},
	}
	_, err = s.storage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)
}

func (s *MassMigrationTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MassMigrationTestSuite) TestAddMassMigration_AddAndRetrieve() {
	err := s.storage.AddMempoolTx(&massMigration)
	s.NoError(err)

	expected := massMigration

	res, err := s.storage.GetMassMigration(massMigration.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *MassMigrationTestSuite) TestAddMassMigration_AddAndRetrieveIncludedMassMigration() {
	includedMassMigration := massMigration
	includedMassMigration.CommitmentSlot = &models.CommitmentSlot{
		BatchID:      models.MakeUint256(3),
		IndexInBatch: 1,
	}
	err := s.storage.AddTransaction(&includedMassMigration)
	s.NoError(err)

	res, err := s.storage.GetMassMigration(massMigration.Hash)
	s.NoError(err)
	s.Equal(includedMassMigration, *res)
}

func (s *MassMigrationTestSuite) TestGetMassMigration_NonexistentMassMigration() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetMassMigration(hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(res)
}

func TestMassMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(MassMigrationTestSuite))
}
