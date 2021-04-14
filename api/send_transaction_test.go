package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/encoder"
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
		AccountIndex: 123,
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
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
	s.api = &API{nil, storage, nil}

	wallet, err := bls.NewRandomWallet(mockDomain)
	s.NoError(err)

	err = storage.AddAccountIfNotExists(&models.Account{
		AccountIndex: 123,
		PublicKey:    *wallet.PublicKey(),
	})
	s.NoError(err)

	err = s.tree.Set(1, &userState)
	s.NoError(err)

	sanitizedTransfer, err := sanitizeTransfer(transferWithoutSignature)
	s.NoError(err)

	encodedTransfer, err := encoder.EncodeTransferForSigning(sanitizedTransfer)
	s.NoError(err)

	signature, err := wallet.Sign(encodedTransfer)
	s.NoError(err)

	s.transfer = transferWithoutSignature
	s.transfer.Signature = signature.Bytes()
}

func (s *SendTransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ReturnsNonNilHash() {
	hash, err := s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ValidatesNonce_TooLow() {
	userStateWithIncreasedNonce := userState
	userStateWithIncreasedNonce.Nonce = *models.NewUint256(1)

	err := s.tree.Set(1, &userStateWithIncreasedNonce)
	s.NoError(err)

	_, err = s.api.SendTransaction(dto.MakeTransaction(s.transfer))
	s.Equal(ErrNonceTooLow, err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ValidatesFee() {
	transferWithZeroFee := s.transfer
	transferWithZeroFee.Fee = models.NewUint256(0)

	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithZeroFee))
	s.Equal(ErrFeeTooLow, err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ValidatesBalance() {
	transferWithHugeAmount := s.transfer
	transferWithHugeAmount.Amount = models.NewUint256(500)
	_, err := s.api.SendTransaction(dto.MakeTransaction(transferWithHugeAmount))
	s.Equal(ErrNotEnoughBalance, err)
}

func (s *SendTransactionTestSuite) TestApi_SendTransaction_ValidatesSignature() {
	wallet, err := bls.NewRandomWallet(mockDomain)
	s.NoError(err)
	fakeSignature, err := wallet.Sign(utils.RandomBytes(2))
	s.NoError(err)

	transfer := transferWithoutSignature
	transfer.Signature = fakeSignature.Bytes()

	_, err = s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.Equal(ErrInvalidSignature, err)
}

func TestSendTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(SendTransactionTestSuite))
}
