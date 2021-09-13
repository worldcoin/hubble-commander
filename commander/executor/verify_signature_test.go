package executor

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
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VerifySignatureTestSuite struct {
	*require.Assertions
	suite.Suite
	executionCtx *ExecutionContext
	storage      *st.TestStorage
	client       *eth.TestClient
	cfg          *config.RollupConfig
	wallets      []bls.Wallet
}

func (s *VerifySignatureTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MaxTxsPerCommitment:    1,
		DisableSignatures:      false,
	}
}

func (s *VerifySignatureTestSuite) SetupTest() {
	var err error
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.executionCtx = NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)
	s.addAccounts()
}

func (s *VerifySignatureTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
	s.client.Close()
}

func (s *VerifySignatureTestSuite) TestVerifyTransferSignature_ValidSignature() {
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

	combinedSignature, err := CombineSignatures(models.TransferArray(transfers), testDomain)
	s.NoError(err)
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: *combinedSignature,
	}

	err = s.executionCtx.verifyTransferSignature(commitment, transfers)
	s.NoError(err)
}

func (s *VerifySignatureTestSuite) TestVerifyTransferSignature_InvalidSignature() {
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

	combinedSignature, err := CombineSignatures(models.TransferArray(transfers), testDomain)
	s.NoError(err)
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: *combinedSignature,
	}

	err = s.executionCtx.verifyTransferSignature(commitment, transfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Signature, disputableErr.Type)
	s.Equal(InvalidSignature, disputableErr.Reason)
}

func (s *VerifySignatureTestSuite) TestVerifyTransferSignature_EmptyTransactions() {
	var transfers []models.Transfer
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: models.Signature{1, 2, 3},
	}

	err := s.executionCtx.verifyTransferSignature(commitment, transfers)
	s.NoError(err)
}

func (s *VerifySignatureTestSuite) TestVerifyCreate2TransferSignature_ValidSignature() {
	transfers := []models.Create2Transfer{
		{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(200),
				Fee:         models.MakeUint256(100),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID:   ref.Uint32(1),
			ToPublicKey: *s.wallets[0].PublicKey(),
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
			ToStateID:   ref.Uint32(0),
			ToPublicKey: *s.wallets[1].PublicKey(),
		},
	}
	for i := range transfers {
		signCreate2Transfer(s.T(), &s.wallets[i], &transfers[i])
	}

	combinedSignature, err := CombineSignatures(models.Create2TransferArray(transfers), testDomain)
	s.NoError(err)
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: *combinedSignature,
	}

	err = s.executionCtx.verifyCreate2TransferSignature(commitment, transfers)
	s.NoError(err)
}

func (s *VerifySignatureTestSuite) TestVerifyCreate2TransfersSignature_EmptyTransactions() {
	var transfers []models.Create2Transfer
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: models.Signature{1, 2, 3},
	}

	err := s.executionCtx.verifyCreate2TransferSignature(commitment, transfers)
	s.NoError(err)
}

func (s *VerifySignatureTestSuite) addAccounts() {
	domain, err := s.client.GetDomain()
	s.NoError(err)

	s.wallets = make([]bls.Wallet, 0, 2)
	for i := uint32(0); i < 2; i++ {
		wallet, err := bls.NewRandomWallet(*domain)
		s.NoError(err)
		s.wallets = append(s.wallets, *wallet)

		err = s.storage.AccountTree.SetSingle(&models.AccountLeaf{
			PubKeyID:  i,
			PublicKey: *wallet.PublicKey(),
		})
		s.NoError(err)
		_, err = s.storage.StateTree.Set(i, &models.UserState{
			PubKeyID: i,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(1000),
			Nonce:    models.MakeUint256(0),
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

func signCreate2Transfer(t *testing.T, wallet *bls.Wallet, transfer *models.Create2Transfer) {
	encodedTransfer, err := encoder.EncodeCreate2TransferForSigning(transfer)
	require.NoError(t, err)
	signature, err := wallet.Sign(encodedTransfer)
	require.NoError(t, err)
	transfer.Signature = *signature.ModelsSignature()
}

func TestSyncTransferCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(VerifySignatureTestSuite))
}
