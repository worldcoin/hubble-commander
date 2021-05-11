package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyTransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown func() error
	storage  *storage.Storage
	tree     *storage.StateTree
	cfg      *config.RollupConfig
}

func (s *ApplyTransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyTransfersTestSuite) SetupTest() {
	testStorage, err := storage.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.NoError(err)
	s.tree = storage.NewStateTree(s.storage)
	s.cfg = &config.RollupConfig{
		FeeReceiverPubKeyID: 3,
		TxsPerCommitment:    6,
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

func (s *ApplyTransfersTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfers_AllValid() {
	transfers := generateValidTransfers(3)

	validTransfers, invalidTransfers, _, err := ApplyTransfers(s.storage, transfers, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 3)
	s.Len(invalidTransfers, 0)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfers_SomeValid() {
	transfers := generateValidTransfers(2)
	transfers = append(transfers, generateInvalidTransfers(3)...)

	validTransfers, invalidTransfers, _, err := ApplyTransfers(s.storage, transfers, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 2)
	s.Len(invalidTransfers, 3)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfers_MoreThan32() {
	transfers := generateValidTransfers(13)

	validTransfers, invalidTransfers, _, err := ApplyTransfers(s.storage, transfers, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 6)
	s.Len(invalidTransfers, 0)

	state, _ := s.tree.Leaf(1)
	s.Equal(models.MakeUint256(6), state.Nonce)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfersTestSuite_SavesTransferErrors() {
	transfers := generateValidTransfers(3)
	transfers = append(transfers, generateInvalidTransfers(2)...)

	for i := range transfers {
		err := s.storage.AddTransfer(&transfers[i])
		s.NoError(err)
	}

	validTransfers, invalidTransfers, _, err := ApplyTransfers(s.storage, transfers, s.cfg)
	s.NoError(err)

	s.Len(validTransfers, 3)
	s.Len(invalidTransfers, 2)

	for i := range transfers {
		transfer, err := s.storage.GetTransfer(transfers[i].Hash)
		s.NoError(err)
		if i < 3 {
			s.Nil(transfer.ErrorMessage)
		} else {
			s.Equal(*transfer.ErrorMessage, ErrNonceTooLow.Error())
		}
	}
}

func TestApplyTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransfersTestSuite))
}

func generateValidTransfers(transfersAmount int) []models.Transfer {
	transfers := make([]models.Transfer, 0, transfersAmount)
	for i := 0; i < transfersAmount; i++ {
		transfer := models.Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(int64(i)),
			},
			ToStateID: 2,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}

func generateInvalidTransfers(transfersAmount int) []models.Transfer {
	transfers := make([]models.Transfer, 0, transfersAmount)
	for i := 0; i < transfersAmount; i++ {
		transfer := models.Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 2,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}
