package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyCreate2TransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown            func() error
	storage             *storage.Storage
	tree                *storage.StateTree
	cfg                 *config.RollupConfig
	client              *eth.TestClient
	publicKey           models.PublicKey
	transactionExecutor *TransactionExecutor
	feeReceiver         *FeeReceiver
	events              chan *accountregistry.AccountRegistrySinglePubkeyRegistered
	unsubscribe         func()
}

func (s *ApplyCreate2TransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyCreate2TransfersTestSuite) SetupTest() {
	testStorage, err := storage.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.tree = storage.NewStateTree(s.storage)
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		FeeReceiverPubKeyID: 3,
		TxsPerCommitment:    6,
	}

	senderState := models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	receiverState := models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	feeReceiverState := models.UserState{
		PubKeyID: 3,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	}
	s.publicKey = models.PublicKey{1, 2, 3}

	for i := 1; i <= 50; i++ {
		err = s.storage.AddAccountIfNotExists(&models.Account{
			PubKeyID:  uint32(i),
			PublicKey: s.publicKey,
		})
		s.NoError(err)
	}

	for i := 1; i <= 10; i++ {
		err = s.storage.AddAccountIfNotExists(&models.Account{
			PubKeyID:  uint32(100 + i),
			PublicKey: s.publicKey,
		})
		s.NoError(err)
	}

	_, err = s.tree.Set(1, &senderState)
	s.NoError(err)
	_, err = s.tree.Set(2, &receiverState)
	s.NoError(err)
	_, err = s.tree.Set(3, &feeReceiverState)
	s.NoError(err)

	s.events, s.unsubscribe, err = s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage, s.client.Client, s.cfg, TransactionExecutorOpts{})
	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func (s *ApplyCreate2TransfersTestSuite) TearDownTest() {
	s.unsubscribe()
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_AllValid() {
	generatedTransfers := generateValidCreate2Transfers(3, &s.publicKey)

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers, s.cfg.TxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 0)
	s.Len(transfers.addedPubKeyIDs, 3)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_SomeValid() {
	generatedTransfers := generateValidCreate2Transfers(2, &s.publicKey)
	generatedTransfers = append(generatedTransfers, generateInvalidCreate2Transfers(3, &s.publicKey)...)

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers, s.cfg.TxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 2)
	s.Len(transfers.invalidTransfers, 3)
	s.Len(transfers.addedPubKeyIDs, 2)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_MoreThanSpecifiedInConfigTxsPerCommitment() {
	generatedTransfers := generateValidCreate2Transfers(13, &s.publicKey)

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers, s.cfg.TxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 6)
	s.Len(transfers.invalidTransfers, 0)
	s.Len(transfers.addedPubKeyIDs, 6)

	state, err := s.storage.GetStateLeaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(6), state.Nonce)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_SavesTransferErrors() {
	generatedTransfers := generateValidCreate2Transfers(3, &s.publicKey)
	generatedTransfers = append(generatedTransfers, generateInvalidCreate2Transfers(2, &s.publicKey)...)

	for i := range generatedTransfers {
		_, err := s.storage.AddCreate2Transfer(&generatedTransfers[i])
		s.NoError(err)
	}

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers, s.cfg.TxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 2)
	s.Len(transfers.addedPubKeyIDs, 3)

	for i := range generatedTransfers {
		transfer, err := s.storage.GetCreate2Transfer(generatedTransfers[i].Hash)
		s.NoError(err)
		if i < 3 {
			s.Nil(transfer.ErrorMessage)
		} else {
			s.Equal(*transfer.ErrorMessage, ErrNonceTooLow.Error())
		}
	}
}

// TODO check the same test for normal Transfer
func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_InvalidSlicesLength() {
	s.T().SkipNow()
	generatedTransfers := generateValidCreate2Transfers(3, &s.publicKey)
	_, err := s.transactionExecutor.ApplyCreate2TransfersForSync(generatedTransfers, []uint32{1, 2}, s.feeReceiver)
	s.Equal(ErrInvalidSliceLength, err)
}

func (s *ApplyCreate2TransfersTestSuite) TestGetOrRegisterPubKeyID_AccountNotExists() {
	c2T := create2Transfer
	c2T.ToPublicKey = models.PublicKey{10, 11, 12}

	pubKeyID, err := s.transactionExecutor.getOrRegisterPubKeyID(s.events, &c2T, models.MakeUint256(1))
	s.NoError(err)
	s.Equal(uint32(0), *pubKeyID)
}

func (s *ApplyCreate2TransfersTestSuite) TestGetOrRegisterPubKeyID_AccountForTokenIDNotExists() {
	c2T := create2Transfer
	c2T.ToPublicKey = s.publicKey

	pubKeyID, err := s.transactionExecutor.getOrRegisterPubKeyID(s.events, &c2T, models.MakeUint256(1))
	s.NoError(err)
	s.Equal(uint32(4), *pubKeyID)
}

func TestApplyCreate2TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyCreate2TransfersTestSuite))
}

func generateValidCreate2Transfers(transfersAmount uint32, publicKey *models.PublicKey) []models.Create2Transfer {
	transfers := make([]models.Create2Transfer, 0, transfersAmount)
	for i := 0; i < int(transfersAmount); i++ {
		transfer := models.Create2Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Create2Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(uint64(i)),
			},
			ToStateID:   nil,
			ToPublicKey: *publicKey,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}

func generateInvalidCreate2Transfers(transfersAmount uint64, publicKey *models.PublicKey) []models.Create2Transfer {
	transfers := make([]models.Create2Transfer, 0, transfersAmount)
	for i := uint64(0); i < transfersAmount; i++ {
		transfer := models.Create2Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Create2Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID:   nil,
			ToPublicKey: *publicKey,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}
