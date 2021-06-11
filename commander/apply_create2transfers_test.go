package commander

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
	transactionExecutor *transactionExecutor
	events              chan *accountregistry.AccountRegistryPubkeyRegistered
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

	err = s.tree.Set(1, &senderState)
	s.NoError(err)
	err = s.tree.Set(2, &receiverState)
	s.NoError(err)
	err = s.tree.Set(3, &feeReceiverState)
	s.NoError(err)

	s.events, s.unsubscribe, err = s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)

	s.transactionExecutor = newTestTransactionExecutor(s.storage, s.client.Client, s.cfg, transactionExecutorOpts{})
}

func (s *ApplyCreate2TransfersTestSuite) TearDownTest() {
	s.unsubscribe()
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_AllValid() {
	generatedTransfers := generateValidCreate2Transfers(3, &s.publicKey)

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 3)
	s.Len(transfers.invalidTransfers, 0)
	s.Len(transfers.addedPubKeyIDs, 3)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_SomeValid() {
	generatedTransfers := generateValidCreate2Transfers(2, &s.publicKey)
	generatedTransfers = append(generatedTransfers, generateInvalidCreate2Transfers(3, &s.publicKey)...)

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers)
	s.NoError(err)

	s.Len(transfers.appliedTransfers, 2)
	s.Len(transfers.invalidTransfers, 3)
	s.Len(transfers.addedPubKeyIDs, 2)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfers_MoreThanSpecifiedInConfigTxsPerCommitment() {
	generatedTransfers := generateValidCreate2Transfers(13, &s.publicKey)

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers)
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
		err := s.storage.AddCreate2Transfer(&generatedTransfers[i])
		s.NoError(err)
	}

	transfers, err := s.transactionExecutor.ApplyCreate2Transfers(generatedTransfers)
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

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_SomeValid() {
	generatedTransfers := generateValidCreate2Transfers(2, &s.publicKey)
	generatedTransfers = append(generatedTransfers, generateInvalidCreate2Transfers(3, &s.publicKey)...)

	transfers, err := s.transactionExecutor.ApplyCreate2TransfersForSync(generatedTransfers, []uint32{1, 2, 3, 4, 5})
	s.NoError(err)
	s.Len(transfers.appliedTransfers, 2)
	s.Len(transfers.invalidTransfers, 3)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_MoreThanSpecifiedInConfigTxsPerCommitment() {
	generatedTransfers := generateValidCreate2Transfers(7, &s.publicKey)
	pubKeyIDs := make([]uint32, 0, len(generatedTransfers))
	for i := range generatedTransfers {
		pubKeyIDs = append(pubKeyIDs, uint32(i+1))
	}

	transfers, err := s.transactionExecutor.ApplyCreate2TransfersForSync(generatedTransfers, pubKeyIDs)
	s.NoError(err)
	s.Len(transfers.appliedTransfers, 6)
	s.Len(transfers.invalidTransfers, 0)

	state, err := s.storage.GetStateLeaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(6), state.Nonce)
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2TransfersForSync_InvalidSlicesLength() {
	generatedTransfers := generateValidCreate2Transfers(3, &s.publicKey)
	_, err := s.transactionExecutor.ApplyCreate2TransfersForSync(generatedTransfers, []uint32{1, 2})
	s.Equal(ErrInvalidSliceLength, err)
}

func (s *ApplyCreate2TransfersTestSuite) TestGetOrRegisterPubKeyID_AccountNotExists() {
	transfer := create2Transfer
	transfer.ToPublicKey = models.PublicKey{10, 11, 12}

	pubKeyID, err := s.transactionExecutor.getOrRegisterPubKeyID(s.events, &transfer, models.MakeUint256(1))
	s.NoError(err)
	s.Equal(uint32(0), *pubKeyID)
}

func (s *ApplyCreate2TransfersTestSuite) TestGetOrRegisterPubKeyID_AccountForTokenIndexNotExists() {
	transfer := create2Transfer
	transfer.ToPublicKey = s.publicKey

	pubKeyID, err := s.transactionExecutor.getOrRegisterPubKeyID(s.events, &transfer, models.MakeUint256(1))
	s.NoError(err)
	s.Equal(uint32(4), *pubKeyID)
}

func (s *ApplyCreate2TransfersTestSuite) TestHandleApplyC2T_ValidTransfer() {
	transfers := generateValidCreate2Transfers(1, &s.publicKey)
	combinedFee := models.NewUint256(100)

	appliedStruct := AppliedC2Transfers{
		appliedTransfers: make([]models.Create2Transfer, 0),
		invalidTransfers: make([]models.Create2Transfer, 0),
	}

	ok, err := s.transactionExecutor.handleApplyC2T(&transfers[0], 1, &appliedStruct, combinedFee, models.NewUint256(1))
	s.NoError(err)
	s.True(ok)
	s.Len(appliedStruct.appliedTransfers, 1)
	s.Len(appliedStruct.invalidTransfers, 0)
	s.Equal(*transfers[0].Amount.AddN(100), *combinedFee)
}

func (s *ApplyCreate2TransfersTestSuite) TestHandleApplyC2T_InvalidTransfer() {
	transfers := generateInvalidCreate2Transfers(1, &s.publicKey)
	transfers[0].Amount = models.MakeUint256(500)
	combinedFee := models.NewUint256(100)

	appliedStruct := AppliedC2Transfers{
		appliedTransfers: make([]models.Create2Transfer, 0),
		invalidTransfers: make([]models.Create2Transfer, 0),
	}

	ok, err := s.transactionExecutor.handleApplyC2T(&transfers[0], 1, &appliedStruct, combinedFee, models.NewUint256(1))
	s.NoError(err)
	s.False(ok)
	s.Len(appliedStruct.appliedTransfers, 0)
	s.Len(appliedStruct.invalidTransfers, 1)
	s.Equal(uint64(100), combinedFee.Uint64())
}

func TestApplyCreate2TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyCreate2TransfersTestSuite))
}

func generateValidCreate2Transfers(transfersAmount int, publicKey *models.PublicKey) []models.Create2Transfer {
	transfers := make([]models.Create2Transfer, 0, transfersAmount)
	for i := 0; i < transfersAmount; i++ {
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

func generateInvalidCreate2Transfers(transfersAmount int, publicKey *models.PublicKey) []models.Create2Transfer {
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
			ToStateID:   nil,
			ToPublicKey: *publicKey,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}
