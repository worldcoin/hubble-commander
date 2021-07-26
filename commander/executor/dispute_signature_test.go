package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DisputeSignatureTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *st.Storage
	teardown            func() error
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
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage, s.client.Client, s.cfg, context.Background())

	s.domain, err = bls.DomainFromBytes(crypto.Keccak256(s.client.ChainState.Rollup.Bytes()))
	s.NoError(err)
}

func (s *DisputeSignatureTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *DisputeSignatureTestSuite) TestUserStateProof() {
	userState := createUserState(1, 300, 1)
	witness, err := s.transactionExecutor.stateTree.Set(1, userState)
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

//TODO: add similar test for Create2Transfer
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

func (s *DisputeSignatureTestSuite) TestDisputeSignature_Transfer() {
	wallets := s.setUserStatesAndAddAccounts()

	transfer := testutils.MakeTransfer(1, 2, 0, 50)
	signTransfer(s.T(), &wallets[1], &transfer)
	pendingBatch, commitments := createTransferBatch(s.Assertions, s.transactionExecutor, &transfer, s.domain)

	err := s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)
	s.client.Commit()

	remoteBatches, err := s.client.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(remoteBatches, 1)

	//TODO: reverted because BNPairingPrecompileCostEstimator is not deployed
	err = s.transactionExecutor.DisputeSignature(&remoteBatches[0], 0)
	s.NoError(err)
}

func (s *DisputeSignatureTestSuite) setUserStatesAndAddAccounts() []bls.Wallet {
	wallets := setUserStates(s.Assertions, s.transactionExecutor)
	for i := range wallets {
		_, err := s.transactionExecutor.accountTree.Set(&models.AccountLeaf{
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
