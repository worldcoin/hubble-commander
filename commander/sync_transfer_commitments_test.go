package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncTransferCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	transactionExecutor *transactionExecutor
	storage             *st.Storage
	tree                *st.StateTree
	client              *eth.TestClient
	teardown            func() error
	cfg                 *config.RollupConfig
	wallets             []bls.Wallet
}

func (s *SyncTransferCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		TxsPerCommitment:       1,
	}
}

func (s *SyncTransferCommitmentsTestSuite) SetupTest() {
	var err error
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.tree = st.NewStateTree(s.storage)
	s.teardown = testStorage.Teardown
	s.transactionExecutor = newTestTransactionExecutor(s.storage, s.client.Client, s.cfg, transactionExecutorOpts{AssumeNonces: true})
	err = s.storage.SetChainState(&s.client.ChainState)
	s.NoError(err)
	s.addAccounts()
}

func (s *SyncTransferCommitmentsTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *SyncTransferCommitmentsTestSuite) TestVerifySignature_ValidSignature() {
	transfers := []models.Transfer{
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(200),
				Fee:         models.MakeUint256(100),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 1,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(150),
				Fee:         models.MakeUint256(10),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 0,
		},
	}
	for i := range transfers {
		signTransfer(s.T(), &s.wallets[i], &transfers[i])
	}

	combinedSignature, err := combineTransferSignatures(transfers, testDomain)
	s.NoError(err)
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: *combinedSignature,
	}

	valid, err := s.transactionExecutor.verifyTransferSignature(commitment, transfers)
	s.NoError(err)
	s.True(valid)
}

func (s *SyncTransferCommitmentsTestSuite) TestVerifySignature_InvalidSignature() {
	transfers := []models.Transfer{
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(200),
				Fee:         models.MakeUint256(100),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 1,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(150),
				Fee:         models.MakeUint256(10),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 0,
		},
	}
	for i := range transfers {
		invalidTransfer := transfers[i]
		invalidTransfer.Nonce = models.MakeUint256(4)
		signTransfer(s.T(), &s.wallets[i], &invalidTransfer)
	}

	combinedSignature, err := combineTransferSignatures(transfers, testDomain)
	s.NoError(err)
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: *combinedSignature,
	}

	valid, err := s.transactionExecutor.verifyTransferSignature(commitment, transfers)
	s.NoError(err)
	s.False(valid)
}

func (s *SyncTransferCommitmentsTestSuite) addAccounts() {
	domain, err := s.storage.GetDomain(s.client.ChainState.ChainID)
	s.NoError(err)

	s.wallets = make([]bls.Wallet, 0, 2)
	for i := uint32(0); i < 2; i++ {
		wallet, err := bls.NewRandomWallet(*domain)
		s.NoError(err)
		s.wallets = append(s.wallets, *wallet)

		err = s.storage.AddAccountIfNotExists(&models.Account{
			PubKeyID:  i,
			PublicKey: *wallet.PublicKey(),
		})
		s.NoError(err)
		err = s.tree.Set(i, &models.UserState{
			PubKeyID:   i,
			TokenIndex: models.MakeUint256(0),
			Balance:    models.MakeUint256(1000),
			Nonce:      models.MakeUint256(0),
		})
		s.NoError(err)
	}
}

func signTransfer(t *testing.T, wallet *bls.Wallet, transfer *models.Transfer) {
	encodedTransfer, err := encoder.EncodeTransferForSigning(transfer)
	require.NoError(t, err)
	signature, err := wallet.Sign(encodedTransfer)
	require.NoError(t, err)
	transfer.Signature = *signature.ModelsSignature()
}

func TestSyncTransferCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(SyncTransferCommitmentsTestSuite))
}
