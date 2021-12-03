package api

import (
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
		cfg:              &config.APIConfig{},
		storage:          s.storage.Storage,
		client:           eth.DomainOnlyTestClient,
		commanderMetrics: metrics.NewCommanderMetrics(),
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

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ReturnsNonNilHash() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.massMigration))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesNonceTooLow_NoTransactions() {
	userStateWithIncreasedNonce := s.userState
	userStateWithIncreasedNonce.Nonce = models.MakeUint256(1)

	_, err := s.storage.StateTree.Set(1, userStateWithIncreasedNonce)
	s.NoError(err)

	_, err = s.api.SendTransaction(dto.MakeTransaction(s.massMigration))
	s.Equal(APIErrNonceTooLow, err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesNonceTooHigh_NoTransactions() {
	massMigrationWithIncreasedNonce := s.massMigration
	massMigrationWithIncreasedNonce.Nonce = models.NewUint256(1)
	massMigrationWithIncreasedNonce = s.signMassMigration(massMigrationWithIncreasedNonce)

	_, err := s.api.SendTransaction(dto.MakeTransaction(massMigrationWithIncreasedNonce))
	s.Equal(APIErrNonceTooHigh, err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesNonceTooHigh_ExistingTransactions() {
	_, err := s.api.SendTransaction(dto.MakeTransaction(s.massMigration))
	s.NoError(err)

	massMigrationWithIncreasedNonce := s.massMigration
	massMigrationWithIncreasedNonce.Nonce = models.NewUint256(2)
	massMigrationWithIncreasedNonce = s.signMassMigration(massMigrationWithIncreasedNonce)

	_, err = s.api.SendTransaction(dto.MakeTransaction(massMigrationWithIncreasedNonce))
	s.Equal(APIErrNonceTooHigh, err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesNonceTooLow_ExistingTransactions() {
	_, err := s.api.SendTransaction(dto.MakeTransaction(s.massMigration))
	s.NoError(err)

	secondMassMigration := s.massMigration
	secondMassMigration.Nonce = models.NewUint256(1)
	secondMassMigration = s.signMassMigration(secondMassMigration)

	_, err = s.api.SendTransaction(dto.MakeTransaction(secondMassMigration))
	s.NoError(err)

	thirdMassMigration := s.massMigration
	thirdMassMigration = s.signMassMigration(thirdMassMigration)

	_, err = s.api.SendTransaction(dto.MakeTransaction(thirdMassMigration))
	s.Equal(APIErrNonceTooLow, err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesFeeValue() {
	massMigrationWithZeroFee := s.massMigration
	massMigrationWithZeroFee.Fee = models.NewUint256(0)

	_, err := s.api.SendTransaction(dto.MakeTransaction(massMigrationWithZeroFee))
	s.Equal(APIErrFeeTooLow, err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesFeeEncodability() {
	massMigrationWithBadFee := s.massMigration
	massMigrationWithBadFee.Fee = models.NewUint256(66666666)

	_, err := s.api.SendTransaction(dto.MakeTransaction(massMigrationWithBadFee))
	s.Equal(APINotDecimalEncodableFeeError, err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesAmountEncodability() {
	massMigrationWithBadAmount := s.massMigration
	massMigrationWithBadAmount.Amount = models.NewUint256(66666666)

	_, err := s.api.SendTransaction(dto.MakeTransaction(massMigrationWithBadAmount))
	s.Equal(APINotDecimalEncodableAmountError, err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesSpokeID() {
	massMigrationWithBadSpokeID := s.massMigration
	massMigrationWithBadSpokeID.SpokeID = ref.Uint32(0)

	_, err := s.api.SendTransaction(dto.MakeTransaction(massMigrationWithBadSpokeID))
	s.Equal(APIErrInvalidSpokeID, err)

	massMigrationWithCorrectSpokeID := s.massMigration
	massMigrationWithCorrectSpokeID.SpokeID = ref.Uint32(1)

	hash, err := s.api.SendTransaction(dto.MakeTransaction(massMigrationWithCorrectSpokeID))
	s.Nil(err)
	s.NotNil(hash)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesAmountValue() {
	massMigrationWithZeroAmount := s.massMigration
	massMigrationWithZeroAmount.Amount = models.NewUint256(0)

	_, err := s.api.SendTransaction(dto.MakeTransaction(massMigrationWithZeroAmount))
	s.Equal(APIErrInvalidAmount, err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesBalance() {
	massMigrationWithHugeAmount := s.massMigration
	massMigrationWithHugeAmount.Amount = models.NewUint256(500)

	_, err := s.api.SendTransaction(dto.MakeTransaction(massMigrationWithHugeAmount))
	s.Equal(APIErrNotEnoughBalance, err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesSignature() {
	wallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	massMigration := massMigrationWithoutSignature
	massMigration.Signature = fakeSignature.ModelsSignature()

	_, err = s.api.SendTransaction(dto.MakeTransaction(massMigration))
	s.Equal(APIErrInvalidSignature, err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_ValidatesSignature_DisabledSignatures() {
	s.api.disableSignatures = true

	wallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	massMigration := massMigrationWithoutSignature
	massMigration.Signature = fakeSignature.ModelsSignature()

	_, err = s.api.SendTransaction(dto.MakeTransaction(massMigration))
	s.NoError(err)
}

func (s *SendMassMigrationTestSuite) TestSendMassMigration_AddsMassMigrationToStorage() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.massMigration))
	s.NoError(err)
	s.NotNil(hash)

	transfer, err := s.api.storage.GetMassMigration(*hash)
	s.NoError(err)
	s.NotNil(transfer)
}

func TestSendMassMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(SendMassMigrationTestSuite))
}
