package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	userState = models.UserState{
		PubKeyID:   123,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}
	transferWithoutSignature = dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   []byte{},
	}
)

type SendTransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	api      *API
	db       *db.TestDB
	tree     *st.StateTree
	transfer dto.Transfer
	wallet   *bls.Wallet
}

func (s *SendTransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SendTransactionTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.db = testDB

	storage := st.NewTestStorage(testDB.DB)
	s.tree = st.NewStateTree(storage)
	s.api = &API{
		cfg:     &config.APIConfig{},
		storage: storage,
		client:  nil,
	}

	s.wallet, err = bls.NewRandomWallet(mockDomain)
	s.NoError(err)

	err = storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  123,
		PublicKey: *s.wallet.PublicKey(),
	})
	s.NoError(err)

	err = s.tree.Set(1, &userState)
	s.NoError(err)

	s.transfer = s.signTransfer(transferWithoutSignature)
}

func (s *SendTransactionTestSuite) signTransfer(transfer dto.Transfer) dto.Transfer {
	signedTransfer, err := SignTransfer(s.wallet, transfer)
	s.NoError(err)
	return *signedTransfer
}

func (s *SendTransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ReturnsNonNilHash() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ValidatesNonceTooLow_NoTransactions() {
	userStateWithIncreasedNonce := userState
	userStateWithIncreasedNonce.Nonce = *models.NewUint256(1)

	err := s.tree.Set(1, &userStateWithIncreasedNonce)
	s.NoError(err)

	_, err = s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.Equal(ErrNonceTooLow, err)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ValidatesNonceTooHigh_NoTransactions() {
	transferWithIncreasedNonce := s.transfer
	transferWithIncreasedNonce.Nonce = models.NewUint256(1)
	transferWithIncreasedNonce = s.signTransfer(transferWithIncreasedNonce)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithIncreasedNonce))
	s.Equal(ErrNonceTooHigh, err)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ValidatesNonceTooHigh_ExistingTransactions() {
	_, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)

	transferWithIncreasedNonce := s.transfer
	transferWithIncreasedNonce.Nonce = models.NewUint256(2)
	transferWithIncreasedNonce = s.signTransfer(transferWithIncreasedNonce)

	_, err = s.api.SendTransaction(dto.MakeTransaction(transferWithIncreasedNonce))
	s.Equal(ErrNonceTooHigh, err)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ValidatesNonceTooLow_ExistingTransactions() {
	_, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)

	secondTransfer := s.transfer
	secondTransfer.Nonce = models.NewUint256(1)
	secondTransfer = s.signTransfer(secondTransfer)

	_, err = s.api.SendTransaction(dto.MakeTransaction(secondTransfer))
	s.NoError(err)

	thirdTransfer := s.transfer
	thirdTransfer = s.signTransfer(thirdTransfer)

	_, err = s.api.SendTransaction(dto.MakeTransaction(thirdTransfer))
	s.Equal(ErrNonceTooLow, err)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ValidatesFeeValue() {
	transferWithZeroFee := s.transfer
	transferWithZeroFee.Fee = models.NewUint256(0)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithZeroFee))
	s.Equal(ErrFeeTooLow, err)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ValidatesFeeEncodability() {
	transferWithBadFee := s.transfer
	transferWithBadFee.Fee = models.NewUint256(66666666)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithBadFee))
	s.Equal(NewNotDecimalEncodableError("fee"), err)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ValidatesAmountEncodability() {
	transferWithBadAmount := s.transfer
	transferWithBadAmount.Amount = models.NewUint256(66666666)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithBadAmount))
	s.Equal(NewNotDecimalEncodableError("amount"), err)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ValidatesBalance() {
	transferWithHugeAmount := s.transfer
	transferWithHugeAmount.Amount = models.NewUint256(500)
	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithHugeAmount))
	s.Equal(ErrNotEnoughBalance, err)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ValidatesSignature() {
	wallet, err := bls.NewRandomWallet(mockDomain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	transfer := transferWithoutSignature
	transfer.Signature = fakeSignature.Bytes()

	_, err = s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.Equal(ErrInvalidSignature, err)
}

func (s *SendTransactionTestSuite) TestSendTransaction_ValidatesSignature_DevMode() {
	s.api.cfg = &config.APIConfig{DevMode: true}

	wallet, err := bls.NewRandomWallet(mockDomain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	transfer := transferWithoutSignature
	transfer.Signature = fakeSignature.Bytes()

	_, err = s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.NoError(err)
}

func TestSendTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(SendTransactionTestSuite))
}
