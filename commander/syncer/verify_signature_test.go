package syncer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
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
	storage *st.TestStorage
	cfg     *config.RollupConfig
	client  *eth.Client
	syncCtx *TxsContext
	wallets []bls.Wallet
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
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client = eth.DomainOnlyTestClient

	s.addAccounts()
}

func (s *VerifySignatureTestSuite) TearDownTest() {
	err := s.storage.Close()
	s.NoError(err)
}

func (s *VerifySignatureTestSuite) TestVerifyTransferSignature_ValidSignature() {
	s.syncCtx = NewTestTxsContext(s.storage.Storage, s.client, s.cfg, txtype.Transfer)

	txs := models.TransferArray{
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
	for i := range txs {
		signTransfer(s.T(), &s.wallets[i], &txs[i])
	}

	combinedSignature, err := executor.CombineSignatures(txs, &bls.TestDomain)
	s.NoError(err)
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: *combinedSignature,
	}

	err = s.syncCtx.verifyTxSignature(commitment, txs)
	s.NoError(err)
}

func (s *VerifySignatureTestSuite) TestVerifyTransferSignature_InvalidSignature() {
	s.syncCtx = NewTestTxsContext(s.storage.Storage, s.client, s.cfg, txtype.Transfer)

	txs := models.TransferArray{
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
	for i := range txs {
		invalidTransfer := txs[i]
		invalidTransfer.Nonce = models.MakeUint256(4)
		signTransfer(s.T(), &s.wallets[i], &invalidTransfer)
	}

	combinedSignature, err := executor.CombineSignatures(txs, &bls.TestDomain)
	s.NoError(err)
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: *combinedSignature,
	}

	err = s.syncCtx.verifyTxSignature(commitment, txs)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Signature, disputableErr.Type)
	s.Equal(InvalidSignatureMessage, disputableErr.Reason)
}

func (s *VerifySignatureTestSuite) TestVerifyTransferSignature_EmptyTransactions() {
	s.syncCtx = NewTestTxsContext(s.storage.Storage, s.client, s.cfg, txtype.Transfer)

	txs := make(models.TransferArray, 0)
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: models.Signature{1, 2, 3},
	}

	err := s.syncCtx.verifyTxSignature(commitment, txs)
	s.NoError(err)
}

func (s *VerifySignatureTestSuite) TestVerifyCreate2TransferSignature_ValidSignature() {
	s.syncCtx = NewTestTxsContext(s.storage.Storage, s.client, s.cfg, txtype.Create2Transfer)

	txs := models.Create2TransferArray{
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
	for i := range txs {
		signCreate2Transfer(s.T(), &s.wallets[i], &txs[i])
	}

	combinedSignature, err := executor.CombineSignatures(txs, &bls.TestDomain)
	s.NoError(err)
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: *combinedSignature,
	}

	err = s.syncCtx.verifyTxSignature(commitment, txs)
	s.NoError(err)
}

func (s *VerifySignatureTestSuite) TestVerifyCreate2TransfersSignature_EmptyTransactions() {
	s.syncCtx = NewTestTxsContext(s.storage.Storage, s.client, s.cfg, txtype.Create2Transfer)

	txs := make(models.Create2TransferArray, 0)
	commitment := &encoder.DecodedCommitment{
		CombinedSignature: models.Signature{1, 2, 3},
	}

	err := s.syncCtx.verifyTxSignature(commitment, txs)
	s.NoError(err)
}

func (s *VerifySignatureTestSuite) TestUserStateProof() {
	s.syncCtx = NewTestTxsContext(s.storage.Storage, s.client, s.cfg, txtype.Transfer)

	userState := &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(300),
		Nonce:    models.MakeUint256(1),
	}
	witness, err := s.storage.StateTree.Set(1, userState)
	s.NoError(err)

	stateProof, err := s.syncCtx.UserStateProof(1)
	s.NoError(err)
	s.Equal(userState, stateProof.UserState)
	s.Equal(witness, stateProof.Witness)
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

func signMassMigration(t *testing.T, wallet *bls.Wallet, massMigration *models.MassMigration) {
	encodedTransfer := encoder.EncodeMassMigrationForSigning(massMigration)
	signature, err := wallet.Sign(encodedTransfer)
	require.NoError(t, err)
	massMigration.Signature = *signature.ModelsSignature()
}

func TestVerifySignatureTestSuite(t *testing.T) {
	suite.Run(t, new(VerifySignatureTestSuite))
}
