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
	transferWithoutSignature = dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
	}
	chainState = models.ChainState{
		ChainID:                        models.MakeUint256(1337),
		AccountRegistry:                utils.RandomAddress(),
		AccountRegistryDeploymentBlock: 9483,
		TokenRegistry:                  utils.RandomAddress(),
		DepositManager:                 utils.RandomAddress(),
		Rollup:                         utils.RandomAddress(),
		SyncedBlock:                    11293,
		GenesisAccounts: []models.PopulatedGenesisAccount{
			{
				PublicKey: models.PublicKey{4, 4, 1, 9},
				PubKeyID:  77,
				StateID:   32,
				Balance:   models.MakeUint256(29384),
			},
			{
				PublicKey: models.PublicKey{7, 3, 1, 1},
				PubKeyID:  443,
				StateID:   293,
				Balance:   models.MakeUint256(3004),
			},
		},
	}
)

type SendTransferTestSuite struct {
	*require.Assertions
	suite.Suite
	api       *API
	storage   *st.TestStorage
	client    *eth.TestClient
	userState *models.UserState
	transfer  dto.Transfer
	wallet    *bls.Wallet
	domain    *bls.Domain
}

func (s *SendTransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SendTransferTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.api = &API{
		cfg:     &config.APIConfig{},
		storage: s.storage.Storage,
		client:  s.client.Client,
	}

	s.domain, err = s.client.GetDomain()
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

	s.transfer = s.signTransfer(transferWithoutSignature)
}

func (s *SendTransferTestSuite) signTransfer(transfer dto.Transfer) dto.Transfer {
	signedTransfer, err := SignTransfer(s.wallet, transfer)
	s.NoError(err)
	return *signedTransfer
}

func (s *SendTransferTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SendTransferTestSuite) TestSendTransfer_ReturnsNonNilHash() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *SendTransferTestSuite) TestSendTransfer_ValidatesNonceTooLow_NoTransactions() {
	userStateWithIncreasedNonce := s.userState
	userStateWithIncreasedNonce.Nonce = models.MakeUint256(1)

	_, err := s.storage.StateTree.Set(1, userStateWithIncreasedNonce)
	s.NoError(err)

	apiErr := &APIError{
		Code:    10005,
		Message: "nonce too low",
	}

	_, err = s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.Equal(apiErr, err)
}

func (s *SendTransferTestSuite) TestSendTransfer_ValidatesNonceTooHigh_NoTransactions() {
	transferWithIncreasedNonce := s.transfer
	transferWithIncreasedNonce.Nonce = models.NewUint256(1)
	transferWithIncreasedNonce = s.signTransfer(transferWithIncreasedNonce)

	apiErr := &APIError{
		Code:    10006,
		Message: "nonce too high",
	}

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithIncreasedNonce))
	s.Equal(apiErr, err)
}

func (s *SendTransferTestSuite) TestSendTransfer_ValidatesNonceTooHigh_ExistingTransactions() {
	_, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)

	transferWithIncreasedNonce := s.transfer
	transferWithIncreasedNonce.Nonce = models.NewUint256(2)
	transferWithIncreasedNonce = s.signTransfer(transferWithIncreasedNonce)

	apiErr := &APIError{
		Code:    10006,
		Message: "nonce too high",
	}

	_, err = s.api.SendTransaction(dto.MakeTransaction(transferWithIncreasedNonce))
	s.Equal(apiErr, err)
}

func (s *SendTransferTestSuite) TestSendTransfer_ValidatesNonceTooLow_ExistingTransactions() {
	_, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)

	secondTransfer := s.transfer
	secondTransfer.Nonce = models.NewUint256(1)
	secondTransfer = s.signTransfer(secondTransfer)

	_, err = s.api.SendTransaction(dto.MakeTransaction(secondTransfer))
	s.NoError(err)

	thirdTransfer := s.transfer
	thirdTransfer = s.signTransfer(thirdTransfer)

	apiErr := &APIError{
		Code:    10005,
		Message: "nonce too low",
	}

	_, err = s.api.SendTransaction(dto.MakeTransaction(thirdTransfer))
	s.Equal(apiErr, err)
}

func (s *SendTransferTestSuite) TestSendTransfer_ValidatesFeeValue() {
	transferWithZeroFee := s.transfer
	transferWithZeroFee.Fee = models.NewUint256(0)

	apiErr := &APIError{
		Code:    10009,
		Message: "fee too low",
	}

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithZeroFee))
	s.Equal(apiErr, err)
}

func (s *SendTransferTestSuite) TestSendTransfer_ValidatesFeeEncodability() {
	transferWithBadFee := s.transfer
	transferWithBadFee.Fee = models.NewUint256(66666666)

	apiErr := &APIError{
		Code:    10012,
		Message: "fee is not encodable as multi-precission decimal",
	}

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithBadFee))
	s.Equal(apiErr, err)
}

func (s *SendTransferTestSuite) TestSendTransfer_ValidatesAmountEncodability() {
	transferWithBadAmount := s.transfer
	transferWithBadAmount.Amount = models.NewUint256(66666666)

	apiErr := &APIError{
		Code:    10011,
		Message: "amount is not encodable as multi-precission decimal",
	}

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithBadAmount))
	s.Equal(apiErr, err)
}

func (s *SendTransferTestSuite) TestSendTransfer_ValidatesAmountValue() {
	transferWithZeroAmount := s.transfer
	transferWithZeroAmount.Amount = models.NewUint256(0)

	apiErr := &APIError{
		Code:    10008,
		Message: "amount must be greater than 0",
	}

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithZeroAmount))
	s.Equal(apiErr, err)
}

func (s *SendTransferTestSuite) TestSendTransfer_ValidatesBalance() {
	transferWithHugeAmount := s.transfer
	transferWithHugeAmount.Amount = models.NewUint256(500)

	apiErr := &APIError{
		Code:    10007,
		Message: "not enough balance",
	}

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithHugeAmount))
	s.Equal(apiErr, err)
}

func (s *SendTransferTestSuite) TestSendTransfer_ValidatesSignature() {
	wallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	transfer := transferWithoutSignature
	transfer.Signature = fakeSignature.ModelsSignature()

	apiErr := &APIError{
		Code:    10010,
		Message: "invalid signature",
	}

	_, err = s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.Equal(apiErr, err)
}

func (s *SendTransferTestSuite) TestSendTransaction_ValidatesSignature_DisabledSignatures() {
	s.api.disableSignatures = true

	wallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	transfer := transferWithoutSignature
	transfer.Signature = fakeSignature.ModelsSignature()

	_, err = s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.NoError(err)
}

func (s *SendTransferTestSuite) TestSendTransfer_AddsTransferToStorage() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)
	s.NotNil(hash)

	transfer, err := s.api.storage.GetTransfer(*hash)
	s.NoError(err)
	s.NotNil(transfer)
}

func TestSendTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(SendTransferTestSuite))
}
