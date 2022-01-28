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
		WithdrawManager:                utils.RandomAddress(),
		Rollup:                         utils.RandomAddress(),
		SyncedBlock:                    11293,
		GenesisAccounts: []models.GenesisAccount{
			{
				PublicKey: models.PublicKey{4, 4, 1, 9},
				StateID:   32,
				State: models.UserState{
					PubKeyID: 77,
					TokenID:  models.MakeUint256(0),
					Balance:  models.MakeUint256(29384),
					Nonce:    models.Uint256{},
				},
			},
			{
				PublicKey: models.PublicKey{7, 3, 1, 1},
				StateID:   293,
				State: models.UserState{
					PubKeyID: 443,
					TokenID:  models.MakeUint256(0),
					Balance:  models.MakeUint256(3004),
					Nonce:    models.Uint256{},
				},
			},
		},
	}
)

type SendTransferTestSuite struct {
	*require.Assertions
	suite.Suite
	api       *API
	storage   *st.TestStorage
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

	s.transfer = s.signTransfer(transferWithoutSignature)
}

func (s *SendTransferTestSuite) signTransfer(transfer dto.Transfer) dto.Transfer {
	signedTransfer, err := SignTransfer(s.wallet, transfer)
	s.NoError(err)
	return *signedTransfer
}

func (s *SendTransferTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SendTransferTestSuite) TestSendTransaction_ReturnsNonNilHash() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *SendTransferTestSuite) TestSendTransaction_ValidatesNonceTooLow_NoTransactions() {
	userStateWithIncreasedNonce := s.userState
	userStateWithIncreasedNonce.Nonce = models.MakeUint256(1)

	_, err := s.storage.StateTree.Set(1, userStateWithIncreasedNonce)
	s.NoError(err)

	_, err = s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.Equal(APIErrNonceTooLow, err)
}

func (s *SendTransferTestSuite) TestSendTransaction_ValidatesFeeValue() {
	transferWithZeroFee := s.transfer
	transferWithZeroFee.Fee = models.NewUint256(0)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithZeroFee))
	s.Equal(APIErrFeeTooLow, err)
}

func (s *SendTransferTestSuite) TestSendTransaction_ValidatesFeeEncodability() {
	transferWithBadFee := s.transfer
	transferWithBadFee.Fee = models.NewUint256(66666666)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithBadFee))
	s.Equal(APINotDecimalEncodableFeeError, err)
}

func (s *SendTransferTestSuite) TestSendTransaction_ValidatesAmountEncodability() {
	transferWithBadAmount := s.transfer
	transferWithBadAmount.Amount = models.NewUint256(66666666)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithBadAmount))
	s.Equal(APINotDecimalEncodableAmountError, err)
}

func (s *SendTransferTestSuite) TestSendTransaction_ValidatesAmountValue() {
	transferWithZeroAmount := s.transfer
	transferWithZeroAmount.Amount = models.NewUint256(0)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithZeroAmount))
	s.Equal(APIErrInvalidAmount, err)
}

func (s *SendTransferTestSuite) TestSendTransaction_ValidatesBalance() {
	transferWithHugeAmount := s.transfer
	transferWithHugeAmount.Amount = models.NewUint256(500)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithHugeAmount))
	s.Equal(APIErrNotEnoughBalance, err)
}

func (s *SendTransferTestSuite) TestSendTransaction_ValidatesSignature() {
	wallet, err := bls.NewRandomWallet(*s.domain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	transfer := transferWithoutSignature
	transfer.Signature = fakeSignature.ModelsSignature()

	_, err = s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.Equal(APIErrInvalidSignature, err)
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

func (s *SendTransferTestSuite) TestSendTransaction_AddsTransferToStorage() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)
	s.NotNil(hash)

	transfer, err := s.api.storage.GetTransfer(*hash)
	s.NoError(err)
	s.NotNil(transfer)
}

func (s *SendTransferTestSuite) TestSendTransaction_UpdatesFailedTransaction() {
	originalHash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)

	err = s.storage.SetTransactionErrors(models.TxError{
		TxHash:       *originalHash,
		ErrorMessage: "some error",
	})
	s.NoError(err)

	originalTx, err := s.storage.GetTransfer(*originalHash)
	s.NoError(err)

	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)
	s.Equal(*originalHash, *hash)

	tx, err := s.storage.GetTransfer(*originalHash)
	s.NoError(err)
	s.Nil(tx.ErrorMessage)
	s.NotEqual(*originalTx.ReceiveTime, tx.ReceiveTime)
}

func (s *SendTransferTestSuite) TestSendTransaction_DoesNotUpdatePendingTransfer() {
	_, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)

	_, err = s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.Equal(APIErrPendingTransaction, err)
}

func (s *SendTransferTestSuite) TestSendTransaction_DoesNotUpdateMinedTransfer() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)

	tx, err := s.storage.GetTransfer(*hash)
	s.NoError(err)

	err = s.storage.MarkTransfersAsIncluded([]models.Transfer{*tx}, &models.CommitmentID{BatchID: models.MakeUint256(1)})
	s.NoError(err)

	_, err = s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.Equal(APIErrMinedTransaction, err)
}

func (s *SendTransferTestSuite) TestSendTransaction_DoesNotAcceptTransactions() {
	s.api.isAcceptingTransactions = false
	_, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.Equal(APIErrTxSendingDisabled, err)
}

func TestSendTransferTestSuite(t *testing.T) {
	suite.Run(t, new(SendTransferTestSuite))
}
