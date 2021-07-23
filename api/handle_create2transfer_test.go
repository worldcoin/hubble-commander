package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	create2TransferWithoutSignature = dto.Create2Transfer{
		FromStateID: ref.Uint32(1),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
	}
)

type SendCreate2TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	api             *API
	teardown        func() error
	tree            *st.StateTree
	userState       *models.UserState
	create2Transfer dto.Create2Transfer
	wallet          *bls.Wallet
	domain          *bls.Domain
}

func (s *SendCreate2TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SendCreate2TransferTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.tree = st.NewStateTree(testStorage.InternalStorage)
	s.api = &API{
		cfg:     &config.APIConfig{},
		storage: testStorage.InternalStorage,
		client: &eth.Client{
			ChainState: chainState,
		},
	}

	err = testStorage.SetChainState(&chainState)
	s.NoError(err)
	s.domain, err = testStorage.GetDomain(chainState.ChainID)
	s.NoError(err)

	s.wallet, err = bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	receiverWallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)

	err = testStorage.AddAccountLeafIfNotExists(&models.AccountLeaf{
		PubKeyID:  123,
		PublicKey: *s.wallet.PublicKey(),
	})
	s.NoError(err)

	err = testStorage.AddAccountLeafIfNotExists(&models.AccountLeaf{
		PubKeyID:  10,
		PublicKey: *receiverWallet.PublicKey(),
	})
	s.NoError(err)

	s.userState = &models.UserState{
		PubKeyID: 123,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}

	_, err = s.tree.Set(1, s.userState)
	s.NoError(err)

	create2TransferWithoutSignature.ToPublicKey = receiverWallet.PublicKey()
	s.create2Transfer = s.signCreate2Transfer(create2TransferWithoutSignature)
}

func (s *SendCreate2TransferTestSuite) signCreate2Transfer(create2Transfer dto.Create2Transfer) dto.Create2Transfer {
	signedTransfer, err := SignCreate2Transfer(s.wallet, create2Transfer)
	s.NoError(err)
	return *signedTransfer
}

func (s *SendCreate2TransferTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ReturnsNonNilHash() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.create2Transfer))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesNonceTooLow_NoTransactions() {
	userStateWithIncreasedNonce := s.userState
	userStateWithIncreasedNonce.Nonce = models.MakeUint256(1)

	_, err := s.tree.Set(1, userStateWithIncreasedNonce)
	s.NoError(err)

	_, err = s.api.SendTransaction(dto.MakeTransaction(s.create2Transfer))
	s.Equal(ErrNonceTooLow, err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesNonceTooHigh_NoTransactions() {
	transferWithIncreasedNonce := s.create2Transfer
	transferWithIncreasedNonce.Nonce = models.NewUint256(1)
	transferWithIncreasedNonce = s.signCreate2Transfer(transferWithIncreasedNonce)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithIncreasedNonce))
	s.Equal(ErrNonceTooHigh, err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesNonceTooHigh_ExistingTransactions() {
	_, err := s.api.SendTransaction(dto.MakeTransaction(s.create2Transfer))
	s.NoError(err)

	transferWithIncreasedNonce := s.create2Transfer
	transferWithIncreasedNonce.Nonce = models.NewUint256(2)
	transferWithIncreasedNonce = s.signCreate2Transfer(transferWithIncreasedNonce)

	_, err = s.api.SendTransaction(dto.MakeTransaction(transferWithIncreasedNonce))
	s.Equal(ErrNonceTooHigh, err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesNonceTooLow_ExistingTransactions() {
	_, err := s.api.SendTransaction(dto.MakeTransaction(s.create2Transfer))
	s.NoError(err)

	secondTransfer := s.create2Transfer
	secondTransfer.Nonce = models.NewUint256(1)
	secondTransfer = s.signCreate2Transfer(secondTransfer)

	_, err = s.api.SendTransaction(dto.MakeTransaction(secondTransfer))
	s.NoError(err)

	thirdTransfer := s.create2Transfer
	thirdTransfer = s.signCreate2Transfer(thirdTransfer)

	_, err = s.api.SendTransaction(dto.MakeTransaction(thirdTransfer))
	s.Equal(ErrNonceTooLow, err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesFeeValue() {
	transferWithZeroFee := s.create2Transfer
	transferWithZeroFee.Fee = models.NewUint256(0)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithZeroFee))
	s.Equal(ErrFeeTooLow, err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesFeeEncodability() {
	transferWithBadFee := s.create2Transfer
	transferWithBadFee.Fee = models.NewUint256(66666666)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithBadFee))
	s.Equal(NewNotDecimalEncodableError("fee"), err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesAmountEncodability() {
	transferWithBadAmount := s.create2Transfer
	transferWithBadAmount.Amount = models.NewUint256(66666666)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithBadAmount))
	s.Equal(NewNotDecimalEncodableError("amount"), err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesAmountValue() {
	transferWithZeroAmount := s.create2Transfer
	transferWithZeroAmount.Amount = models.NewUint256(0)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithZeroAmount))
	s.Equal(ErrInvalidAmount, err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesBalance() {
	transferWithHugeAmount := s.create2Transfer
	transferWithHugeAmount.Amount = models.NewUint256(500)
	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithHugeAmount))
	s.Equal(ErrNotEnoughBalance, err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesSignature() {
	wallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	transfer := create2TransferWithoutSignature
	transfer.Signature = fakeSignature.ModelsSignature()

	_, err = s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.Equal(ErrInvalidSignature, err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_ValidatesSignature_DevMode() {
	s.api.devMode = true

	wallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	transfer := create2TransferWithoutSignature
	transfer.Signature = fakeSignature.ModelsSignature()

	_, err = s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.NoError(err)
}

func (s *SendCreate2TransferTestSuite) TestSendCreate2Transfer_AddsTransferToStorage() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.create2Transfer))
	s.NoError(err)
	s.NotNil(hash)

	transfer, err := s.api.storage.GetCreate2Transfer(*hash)
	s.NoError(err)
	s.NotNil(transfer)
}

func TestSendCreate2TransferTestSuite(t *testing.T) {
	suite.Run(t, new(SendCreate2TransferTestSuite))
}
