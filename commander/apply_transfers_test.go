package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
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
	teardown            func() error
	storage             *storage.Storage
	tree                *storage.StateTree
	cfg                 *config.RollupConfig
	transactionExecutor *TransactionExecutor
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

	s.transactionExecutor = NewTestTransactionExecutor(s.storage, &eth.Client{}, s.cfg, TransactionExecutorOpts{})
}

func (s *ApplyTransfersTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfers_AllValid() {
	generatedTransfers := generateValidTransfers(3)

	transfers, err := s.transactionExecutor.ApplyTransfers(generatedTransfers)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 0)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfers_SomeValid() {
	generatedTransfers := generateValidTransfers(2)
	generatedTransfers = append(generatedTransfers, generateInvalidTransfers(3)...)

	transfers, err := s.transactionExecutor.ApplyTransfers(generatedTransfers)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 2)
	s.Len(transfers.invalidTransfers, 3)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfers_MoreThan32() {
	generatedTransfers := generateValidTransfers(13)

	transfers, err := s.transactionExecutor.ApplyTransfers(generatedTransfers)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 6)
	s.Len(transfers.invalidTransfers, 0)

	state, err := s.storage.GetStateLeaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(6), state.Nonce)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfersTestSuite_SavesTransferErrors() {
	generatedTransfers := generateValidTransfers(3)
	generatedTransfers = append(generatedTransfers, generateInvalidTransfers(2)...)

	for i := range generatedTransfers {
		err := s.storage.AddTransfer(&generatedTransfers[i])
		s.NoError(err)
	}

	transfers, err := s.transactionExecutor.ApplyTransfers(generatedTransfers)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 2)

	for i := range generatedTransfers {
		transfer, err := s.storage.GetTransfer(generatedTransfers[i].Hash)
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
				Nonce:       models.MakeUint256(uint64(i)),
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
