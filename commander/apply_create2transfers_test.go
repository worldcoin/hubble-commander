package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
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
	cfg     *config.RollupConfig
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
		TxsPerCommitment: 6,
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

	for i := 1; i <= 50; i++ {
		err = s.storage.AddAccountIfNotExists(&models.Account{
			PubKeyID:  uint32(i),
			PublicKey: models.PublicKey{1, 2, 3},
		})
		s.NoError(err)
	}

	for i := 1; i <= 10; i++ {
		err = s.storage.AddAccountIfNotExists(&models.Account{
			PubKeyID:  uint32(100 + i),
			PublicKey: models.PublicKey{1, 2, 3},
		})
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
	transfers := generateValidCreate2Transfers(3)

	addedAccounts := make(map[uint32]struct{})
	validTransfers, invalidTransfers, err := ApplyCreate2Transfers(s.storage, transfers, addedAccounts, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 3)
	s.Len(invalidTransfers, 0)
	s.Len(addedAccounts, 3)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_SomeValid() {
	transfers := generateValidCreate2Transfers(2)
	transfers = append(transfers, generateInvalidCreate2Transfers(3)...)

	addedAccounts := make(map[uint32]struct{})
	validTransfers, invalidTransfers, err := ApplyCreate2Transfers(s.storage, transfers, addedAccounts, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 2)
	s.Len(invalidTransfers, 3)
	s.Len(addedAccounts, 2)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_MoreThanSpecifiedInConfigTxsPerCommitment() {
	transfers := generateValidCreate2Transfers(13)

	addedAccounts := make(map[uint32]struct{})
	validTransfers, invalidTransfers, err := ApplyCreate2Transfers(s.storage, transfers, addedAccounts, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 6)
	s.Len(invalidTransfers, 0)
	s.Len(addedAccounts, 6)

	state, _ := s.tree.Leaf(1)
	s.Equal(models.MakeUint256(6), state.Nonce)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_FailsTransferIfAccountWasAlreadyAdded() {
	transfers := generateValidCreate2Transfers(3)
	transfers[2].ToPubKeyID = 10

	addedAccounts := map[uint32]struct{}{
		10: {},
	}

	validTransfers, invalidTransfers, err := ApplyCreate2Transfers(s.storage, transfers, addedAccounts, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 2)
	s.Len(invalidTransfers, 1)
	s.Len(addedAccounts, 3)
	s.Contains(addedAccounts, transfers[0].ToPubKeyID)
	s.Contains(addedAccounts, transfers[1].ToPubKeyID)
	s.Contains(addedAccounts, transfers[2].ToPubKeyID)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_SavesTransferErrors() {
	transfers := generateValidCreate2Transfers(3)
	transfers[2].ToPubKeyID = 10
	transfers = append(transfers, generateInvalidCreate2Transfers(2)...)

	for i := range transfers {
		err := s.storage.AddCreate2Transfer(&transfers[i])
		s.NoError(err)
	}

	addedAccounts := map[uint32]struct{}{
		10: {},
	}

	validTransfers, invalidTransfers, err := ApplyCreate2Transfers(s.storage, transfers, addedAccounts, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 2)
	s.Len(invalidTransfers, 3)
	s.Len(addedAccounts, 3)

	for i := range transfers {
		transfer, err := s.storage.GetCreate2Transfer(transfers[i].Hash)
		s.NoError(err)
		if i < 2 {
			s.Nil(transfer.ErrorMessage)
		} else if i == 2 {
			s.Equal(*transfer.ErrorMessage, ErrAccountAlreadyExists.Error())
		} else {
			s.Equal(*transfer.ErrorMessage, ErrNonceTooLow.Error())
		}
	}
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
				TxType:      txtype.Create2Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(int64(i)),
			},
			ToStateID:  2,
			ToPubKeyID: uint32(i + 1),
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
				TxType:      txtype.Create2Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID:  2,
			ToPubKeyID: uint32(i + 101),
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}
