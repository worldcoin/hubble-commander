package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyCreate2TransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	db      *db.TestDB
	storage *storage.Storage
	tree    *storage.StateTree
	cfg *config.RollupConfig
}

func (s *ApplyCreate2TransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyCreate2TransfersTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.db = testDB
	s.storage = storage.NewTestStorage(testDB.DB)
	s.tree = storage.NewStateTree(s.storage)
	s.cfg = &config.RollupConfig{
		FeeReceiverIndex: 3,
		TxsPerCommitment: 32,
	}

	senderState := models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}
	receiverState := models.UserState{
		PubKeyID:   2,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	}
	feeReceiverState := models.UserState{
		PubKeyID:   3,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(0),
	}

	accounts := []models.Account{
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  3,
			PublicKey: models.PublicKey{1, 2, 3},
		},
	}
	for i := range accounts {
		err = s.storage.AddAccountIfNotExists(&accounts[i])
		s.NoError(err)
	}

	err = s.tree.Set(1, &senderState)
	s.NoError(err)
	err = s.tree.Set(2, &receiverState)
	s.NoError(err)
	err = s.tree.Set(3, &feeReceiverState)
	s.NoError(err)
}

func (s *ApplyCreate2TransfersTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_AllValid() {
	transfers := generateValidCreate2Transfers(10)

	validTransfers, err := ApplyCreate2Transfers(s.storage, transfers, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 10)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_SomeValid() {
	transfers := generateValidCreate2Transfers(10)
	transfers = append(transfers, generateInvalidCreate2Transfers(10)...)

	validTransfers, err := ApplyCreate2Transfers(s.storage, transfers, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 10)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_MoreThan32() {
	transfers := generateValidCreate2Transfers(60)

	validTransfers, err := ApplyCreate2Transfers(s.storage, transfers, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 32)

	state, _ := s.tree.Leaf(1)
	s.Equal(models.MakeUint256(32), state.Nonce)
}

func TestApplyCreate2TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyCreate2TransfersTestSuite))
}

func generateValidCreate2Transfers(transfersAmount int) []models.Create2Transfer {
	transfers := make([]models.Create2Transfer, 0, transfersAmount)
	for i := 0; i < transfersAmount; i++ {
		transfer := models.Create2Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(int64(i)),
			},
			ToStateID: 2,
			ToPubKeyID: 2,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}

func generateInvalidCreate2Transfers(transfersAmount int) []models.Create2Transfer {
	transfers := make([]models.Create2Transfer, 0, transfersAmount)
	for i := 0; i < transfersAmount; i++ {
		transfer := models.Create2Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 2,
			ToPubKeyID: 2,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}
