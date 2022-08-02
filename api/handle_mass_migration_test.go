package api

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	massMigrationWithoutSignature = dto.MassMigration{
		FromStateID: ref.Uint32(1),
		SpokeID:     ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
	}
)

type SendMassMigrationTestSuite struct {
	*require.Assertions
	suite.Suite
	api           *API
	storage       *st.TestStorage
	userState     *models.UserState
	massMigration dto.MassMigration
	wallet        *bls.Wallet
	domain        *bls.Domain
}

func (s *SendMassMigrationTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SendMassMigrationTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{
		cfg:                     &config.APIConfig{},
		storage:                 s.storage.Storage,
		client:                  eth.DomainOnlyTestClient,
		commanderMetrics:        metrics.NewCommanderMetrics(),
		isAcceptingTransactions: true,
	}

	s.domain, err = s.api.client.GetDomain()
	s.NoError(err)
	s.wallet, err = bls.NewRandomWallet(*s.domain)
	s.NoError(err)

	err = s.storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  123,
		PublicKey: *s.wallet.PublicKey(),
	})
	s.NoError(err)

	s.userState = &models.UserState{
		PubKeyID: 123,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}

	_, err = s.storage.StateTree.Set(1, s.userState)
	s.NoError(err)

	s.massMigration = s.signMassMigration(massMigrationWithoutSignature)

	err = s.storage.AddRegisteredSpoke(&models.RegisteredSpoke{
		ID: models.MakeUint256(2),
	})
	s.NoError(err)
}

func (s *SendMassMigrationTestSuite) signMassMigration(massMigration dto.MassMigration) dto.MassMigration {
	signedMassMigration, err := SignMassMigration(s.wallet, massMigration)
	s.NoError(err)
	return *signedMassMigration
}

func (s *SendMassMigrationTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_ReturnsNonNilHash() {
	hash, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(s.massMigration))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_ValidatesNonceTooLow_NoTransactions() {
	userStateWithIncreasedNonce := s.userState
	userStateWithIncreasedNonce.Nonce = models.MakeUint256(1)

	_, err := s.storage.StateTree.Set(1, userStateWithIncreasedNonce)
	s.NoError(err)

	_, err = s.api.SendTransaction(context.Background(), dto.MakeTransaction(s.massMigration))
	s.Equal(APIErrNonceTooLow, err)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_ValidatesFeeValue() {
	massMigrationWithZeroFee := s.massMigration
	massMigrationWithZeroFee.Fee = models.NewUint256(0)

	_, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(massMigrationWithZeroFee))
	s.Equal(APIErrFeeTooLow, err)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_ValidatesFeeEncodability() {
	massMigrationWithBadFee := s.massMigration
	massMigrationWithBadFee.Fee = models.NewUint256(66666666)

	_, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(massMigrationWithBadFee))
	s.Equal(APINotDecimalEncodableFeeError, err)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_ValidatesAmountEncodability() {
	massMigrationWithBadAmount := s.massMigration
	massMigrationWithBadAmount.Amount = models.NewUint256(66666666)

	_, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(massMigrationWithBadAmount))
	s.Equal(APINotDecimalEncodableAmountError, err)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_ValidatesThatSpokeExists() {
	massMigrationWithNonexistentSpoke := s.massMigration
	massMigrationWithNonexistentSpoke.SpokeID = ref.Uint32(1)

	_, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(massMigrationWithNonexistentSpoke))
	s.Equal(APIErrSpokeDoesNotExist, err)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_ValidatesAmountValue() {
	massMigrationWithZeroAmount := s.massMigration
	massMigrationWithZeroAmount.Amount = models.NewUint256(0)

	_, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(massMigrationWithZeroAmount))
	s.Equal(APIErrInvalidAmount, err)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_ValidatesBalance() {
	massMigrationWithHugeAmount := s.massMigration
	massMigrationWithHugeAmount.Amount = models.NewUint256(500)

	_, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(massMigrationWithHugeAmount))
	s.Equal(APIErrNotEnoughBalance, err)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_ValidatesSignature() {
	wallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	massMigration := massMigrationWithoutSignature
	massMigration.Signature = fakeSignature.ModelsSignature()

	_, err = s.api.SendTransaction(context.Background(), dto.MakeTransaction(massMigration))
	s.Equal(APIErrInvalidSignature, err)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_ValidatesSignature_DisabledSignatures() {
	s.api.disableSignatures = true

	wallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	massMigration := massMigrationWithoutSignature
	massMigration.Signature = fakeSignature.ModelsSignature()

	_, err = s.api.SendTransaction(context.Background(), dto.MakeTransaction(massMigration))
	s.NoError(err)
}

func (s *SendMassMigrationTestSuite) TestSendTransaction_AddsMassMigrationToStorage() {
	hash, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(s.massMigration))
	s.NoError(err)
	s.NotNil(hash)

	transfer, err := s.api.storage.GetMassMigration(*hash)
	s.NoError(err)
	s.NotNil(transfer)
}

// TODO: figure out what to do about this test
/*
func (s *SendMassMigrationTestSuite) TestSendTransaction_UpdatesFailedTransaction() {
	originalHash, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(s.massMigration))
	s.NoError(err)

	err = s.storage.SetTransactionErrors(models.TxError{
		TxHash:       *originalHash,
		ErrorMessage: "some error",
	})
	s.NoError(err)

	originalTx, err := s.storage.GetMassMigration(*originalHash)
	s.NoError(err)

	hash, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(s.massMigration))
	s.NoError(err)
	s.Equal(*originalHash, *hash)

	tx, err := s.storage.GetMassMigration(*originalHash)
	s.NoError(err)
	s.Nil(tx.ErrorMessage)
	s.NotEqual(*originalTx.ReceiveTime, tx.ReceiveTime)
}
*/

func TestSendMassMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(SendMassMigrationTestSuite))
}
