package executor

import (
	"context"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DisputeSignatureTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *st.TestStorage
	client              *eth.TestClient
	cfg                 *config.RollupConfig
	transactionExecutor *TransactionExecutor
	domain              *bls.Domain
}

func (s *DisputeSignatureTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
		DevMode:                false,
	}
}

func (s *DisputeSignatureTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)

	s.client, err = eth.NewConfiguredTestClient(
		rollup.DeploymentConfig{},
		eth.ClientConfig{TxTimeout: ref.Duration(2 * time.Second)},
	)
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, s.client.Client, s.cfg, context.Background())

	s.domain, err = s.client.GetDomain()
	s.NoError(err)
}

func (s *DisputeSignatureTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DisputeSignatureTestSuite) TestUserStateProof() {
	userState := createUserState(1, 300, 1)
	witness, err := s.transactionExecutor.storage.StateTree.Set(1, userState)
	s.NoError(err)

	stateProof, err := s.transactionExecutor.userStateProof(1)
	s.NoError(err)
	s.Equal(userState, stateProof.UserState)
	s.Equal(witness, stateProof.Witness)
}

func (s *DisputeSignatureTestSuite) TestPublicKeyProof() {
	account := &models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.storage.AddAccountLeafIfNotExists(account)
	s.NoError(err)

	publicKeyProof, err := s.transactionExecutor.publicKeyProof(account.PubKeyID)
	s.NoError(err)
	s.Equal(account.PublicKey, *publicKeyProof.PublicKey)
	s.Len(publicKeyProof.Witness, 32)
}

func (s *DisputeSignatureTestSuite) TestReceiverPublicKeyProof() {
	account := &models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.storage.AddAccountLeafIfNotExists(account)
	s.NoError(err)

	publicKeyHash := crypto.Keccak256Hash(account.PublicKey.Bytes())

	publicKeyProof, err := s.transactionExecutor.receiverPublicKeyProof(account.PubKeyID)
	s.NoError(err)
	s.Equal(publicKeyHash, publicKeyProof.PublicKeyHash)
	s.Len(publicKeyProof.Witness, 32)
}

func (s *DisputeSignatureTestSuite) TestSignatureProof() {
	s.setUserStatesAndAddAccounts()

	transfers := []models.Transfer{
		testutils.MakeTransfer(1, 2, 0, 50),
		testutils.MakeTransfer(0, 2, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 75),
	}

	expectedUserStates := make([]models.UserState, 0, len(transfers))
	expectedPublicKeys := make([]models.PublicKey, 0, len(transfers))
	for i := range transfers {
		leaf, err := s.storage.GetStateLeaf(transfers[i].FromStateID)
		s.NoError(err)
		expectedUserStates = append(expectedUserStates, leaf.UserState)

		publicKey, err := s.storage.GetPublicKey(leaf.PubKeyID)
		s.NoError(err)
		expectedPublicKeys = append(expectedPublicKeys, *publicKey)
	}

	serializedTxs, err := encoder.SerializeTransfers(transfers)
	s.NoError(err)

	signatureProof, err := s.transactionExecutor.signatureProof(&encoder.DecodedCommitment{Transactions: serializedTxs})
	s.NoError(err)
	s.Len(signatureProof.UserStates, 3)
	s.Len(signatureProof.PublicKeys, 3)

	for i := range signatureProof.UserStates {
		s.Equal(expectedUserStates[i], *signatureProof.UserStates[i].UserState)
		s.Equal(expectedPublicKeys[i], *signatureProof.PublicKeys[i].PublicKey)
	}
}

func (s *DisputeSignatureTestSuite) TestSignatureProofWithReceiver() {
	wallets := s.setUserStatesAndAddAccounts()

	receiverAccounts := []models.AccountLeaf{
		{PubKeyID: 3, PublicKey: *wallets[2].PublicKey()},
		{PubKeyID: 4, PublicKey: *wallets[2].PublicKey()},
		{PubKeyID: 5, PublicKey: *wallets[1].PublicKey()},
	}

	transfers := []models.Create2Transfer{
		testutils.MakeCreate2Transfer(1, ref.Uint32(3), 0, 50, &receiverAccounts[0].PublicKey),
		testutils.MakeCreate2Transfer(0, ref.Uint32(4), 0, 50, &receiverAccounts[1].PublicKey),
		testutils.MakeCreate2Transfer(0, ref.Uint32(5), 1, 75, &receiverAccounts[2].PublicKey),
	}
	pubKeyIDs := []uint32{3, 4, 5}

	expectedUserStates := make([]models.UserState, 0, len(transfers))
	senderPublicKeys := make([]models.PublicKey, 0, len(transfers))
	receiverPublicKeys := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		leaf, err := s.storage.GetStateLeaf(transfers[i].FromStateID)
		s.NoError(err)
		expectedUserStates = append(expectedUserStates, leaf.UserState)

		publicKey, err := s.storage.GetPublicKey(leaf.PubKeyID)
		s.NoError(err)
		senderPublicKeys = append(senderPublicKeys, *publicKey)

		err = s.transactionExecutor.storage.AccountTree.SetSingle(&receiverAccounts[i])
		s.NoError(err)
		receiverPublicKeys = append(receiverPublicKeys, crypto.Keccak256Hash(transfers[i].ToPublicKey.Bytes()))
	}

	serializedTxs, err := encoder.SerializeCreate2Transfers(transfers, pubKeyIDs)
	s.NoError(err)

	signatureProof, err := s.transactionExecutor.signatureProofWithReceiver(&encoder.DecodedCommitment{Transactions: serializedTxs})
	s.NoError(err)
	s.Len(signatureProof.UserStates, 3)
	s.Len(signatureProof.SenderPublicKeys, 3)
	s.Len(signatureProof.ReceiverPublicKeys, 3)

	for i := range signatureProof.UserStates {
		s.Equal(expectedUserStates[i], *signatureProof.UserStates[i].UserState)
		s.Equal(senderPublicKeys[i], *signatureProof.SenderPublicKeys[i].PublicKey)
		s.Equal(receiverPublicKeys[i], signatureProof.ReceiverPublicKeys[i].PublicKeyHash)
	}
}

func (s *DisputeSignatureTestSuite) TestDisputeSignature_DisputesTransferBatchWithInvalidSignature() {
	wallets := s.setUserStatesAndAddAccounts()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	signTransfer(s.T(), &wallets[0], &transfer)

	createAndSubmitTransferBatch(s.Assertions, s.client, s.transactionExecutor, &transfer)

	remoteBatches, err := s.client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.transactionExecutor.DisputeSignature(&remoteBatches[0], 0)
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].ID)
}

func (s *DisputeSignatureTestSuite) TestDisputeSignature_DisputesC2TBatchWithInvalidSignature() {
	wallets := s.setUserStatesAndAddAccounts()

	receiver := &models.AccountLeaf{
		PubKeyID:  3,
		PublicKey: *wallets[2].PublicKey(),
	}

	transfer := testutils.MakeCreate2Transfer(0, &receiver.PubKeyID, 0, 100, &receiver.PublicKey)
	signCreate2Transfer(s.T(), &wallets[1], &transfer)

	createAndSubmitC2TBatch(s.Assertions, s.client, s.transactionExecutor, &transfer)

	err := s.transactionExecutor.storage.AccountTree.SetSingle(receiver)
	s.NoError(err)

	remoteBatches, err := s.client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.transactionExecutor.DisputeSignature(&remoteBatches[0], 0)
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].ID)
}

func (s *DisputeSignatureTestSuite) TestDisputeSignature_Transfer_ValidBatch() {
	wallets := s.setUserStatesAndAddAccounts()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	signTransfer(s.T(), &wallets[1], &transfer)

	createAndSubmitTransferBatch(s.Assertions, s.client, s.transactionExecutor, &transfer)

	remoteBatches, err := s.client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.transactionExecutor.DisputeSignature(&remoteBatches[0], 0)
	s.ErrorIs(err, eth.ErrWaitForRollbackTimeout)
}

func (s *DisputeSignatureTestSuite) TestDisputeSignature_Create2Transfer_ValidBatch() {
	wallets := s.setUserStatesAndAddAccounts()

	receiver := &models.AccountLeaf{
		PubKeyID:  3,
		PublicKey: *wallets[2].PublicKey(),
	}

	transfer := testutils.MakeCreate2Transfer(0, &receiver.PubKeyID, 0, 100, &receiver.PublicKey)
	signCreate2Transfer(s.T(), &wallets[0], &transfer)

	createAndSubmitC2TBatch(s.Assertions, s.client, s.transactionExecutor, &transfer)

	err := s.transactionExecutor.storage.AccountTree.SetSingle(receiver)
	s.NoError(err)

	remoteBatches, err := s.client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.transactionExecutor.DisputeSignature(&remoteBatches[0], 0)
	s.ErrorIs(err, eth.ErrWaitForRollbackTimeout)
}

func (s *DisputeSignatureTestSuite) setUserStatesAndAddAccounts() []bls.Wallet {
	wallets := setUserStates(s.Assertions, s.transactionExecutor, s.domain)
	for i := range wallets {
		err := s.transactionExecutor.storage.AccountTree.SetSingle(&models.AccountLeaf{
			PubKeyID:  uint32(i),
			PublicKey: *wallets[i].PublicKey(),
		})
		s.NoError(err)
	}
	return wallets
}

func TestDisputeSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeSignatureTestSuite))
}
