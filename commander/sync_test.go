package commander

import (
	"context"
	"sync"
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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown            func() error
	storage             *st.Storage
	tree                *st.StateTree
	client              *eth.TestClient
	cfg                 *config.RollupConfig
	transactionExecutor *transactionExecutor
	stateMutex          *sync.Mutex
}

func (s *SyncTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SyncTestSuite) SetupTest() {
	var err error
	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		TxsPerCommitment:       1,
	}

	s.stateMutex = &sync.Mutex{}

	s.setupDB()
}

func (s *SyncTestSuite) setupDB() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.tree = st.NewStateTree(s.storage)
	s.transactionExecutor = newTestTransactionExecutor(s.storage, s.client.Client, s.cfg)

	s.seedDB()
}

func (s *SyncTestSuite) seedDB() {
	err := s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  0,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	s.NoError(err)

	err = s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  1,
		PublicKey: models.PublicKey{2, 3, 4},
	})
	s.NoError(err)

	err = s.tree.Set(0, &models.UserState{
		PubKeyID:   0,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)

	err = s.tree.Set(1, &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *SyncTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *SyncTestSuite) TestSyncBatches_TwoTransferBatches() {
	accountRoot := s.getAccountTreeRoot()

	txs := []models.Transfer{
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(400),
				Fee:         models.MakeUint256(0),
				Nonce:       models.MakeUint256(0),
				Signature:   *mockSignature(s.T()),
			},
			ToStateID: 1,
		}, {
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(100),
				Fee:         models.MakeUint256(0),
				Nonce:       models.MakeUint256(0),
				Signature:   *mockSignature(s.T()),
			},
			ToStateID: 0,
		},
	}
	for i := range txs {
		s.setTransferHash(&txs[i])
		err := s.storage.AddTransfer(&txs[i])
		s.NoError(err)
	}

	expectedCommitments := make([]models.Commitment, 2)
	for i := range expectedCommitments {
		createdCommitments, err := s.transactionExecutor.createTransferCommitments([]models.Transfer{txs[i]}, testDomain)
		s.NoError(err)
		s.Len(createdCommitments, 1)

		expectedCommitments[i] = createdCommitments[0]
		_, err = s.transactionExecutor.submitBatch(txtype.Transfer, createdCommitments)
		s.NoError(err)
		s.client.Commit()
	}

	s.recreateDatabase()
	s.syncAllBlocks()

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)
	s.Equal(models.MakeUint256(1), batches[0].Number)
	s.Equal(models.MakeUint256(2), batches[1].Number)

	for i := range expectedCommitments {
		commitment, err := s.storage.GetCommitment(expectedCommitments[i].ID)
		s.NoError(err)
		expectedCommitments[i].IncludedInBatch = &batches[i].ID
		expectedCommitments[i].AccountTreeRoot = &accountRoot
		s.Equal(expectedCommitments[i], *commitment)

		actualTx, err := s.storage.GetTransfer(txs[i].Hash)
		s.NoError(err)
		txs[i].IncludedInCommitment = &expectedCommitments[i].ID
		txs[i].Signature = models.Signature{}
		s.Equal(txs[i], *actualTx)
	}
}

func (s *SyncTestSuite) TestSyncBatches_DoesNotSyncExistingBatchTwice() {
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *mockSignature(s.T()),
		},
		ToStateID: 1,
	}
	s.createAndSubmitTransferBatch(&tx)

	s.recreateDatabase()
	s.syncAllBlocks()

	// Begin database transaction
	var err error
	s.transactionExecutor, err = newTransactionExecutorWithCtx(context.Background(), s.storage, s.client.Client, s.cfg)
	s.NoError(err)

	tx2 := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *mockSignature(s.T()),
		},
		ToStateID: 0,
	}
	s.createAndSubmitTransferBatch(&tx2)

	batches, err := s.transactionExecutor.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	// Rollback changes to the database
	s.transactionExecutor.Rollback(nil)
	s.transactionExecutor = newTestTransactionExecutor(s.storage, s.client.Client, s.cfg)

	batches, err = s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)

	s.syncAllBlocks()

	state0, err := s.storage.GetStateLeaf(0)
	s.NoError(err)
	s.Equal(models.MakeUint256(700), state0.Balance)

	state1, err := s.storage.GetStateLeaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(300), state1.Balance)

	batches, err = s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)
}

func (s *SyncTestSuite) TestSyncBatches_PendingBatch() {
	accountRoot := s.getAccountTreeRoot()
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *mockSignature(s.T()),
		},
		ToStateID: 1,
	}
	s.setTransferHash(&tx)
	s.createAndSubmitTransferBatch(&tx)

	pendingBatch, err := s.storage.GetBatch(1)
	s.NoError(err)
	s.Nil(pendingBatch.Hash)
	s.Nil(pendingBatch.FinalisationBlock)

	s.syncAllBlocks()

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
	s.NotNil(batches[0].Hash)
	s.NotNil(batches[0].FinalisationBlock)

	commitment, err := s.storage.GetCommitment(1)
	s.NoError(err)
	s.Equal(accountRoot, *commitment.AccountTreeRoot)
}

func (s *SyncTestSuite) TestSyncBatches_Create2Transfer() {
	s.registerAccountOnChain(&models.PublicKey{1, 2, 3}, 0)
	tx := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Create2Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *mockSignature(s.T()),
		},
		ToStateID:   ref.Uint32(5),
		ToPublicKey: models.PublicKey{},
	}
	s.setCreate2TransferHash(&tx)
	expectedCommitment := s.createAndSubmitC2TBatch(&tx)

	s.recreateDatabase()
	s.syncAllBlocks()

	state0, err := s.storage.GetStateLeaf(0)
	s.NoError(err)
	s.Equal(models.MakeUint256(600), state0.Balance)

	state5, err := s.storage.GetStateLeaf(5)
	s.NoError(err)
	s.Equal(models.MakeUint256(400), state5.Balance)
	s.Equal(uint32(1), state5.PubKeyID)

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)

	commitment, err := s.storage.GetCommitment(expectedCommitment.ID)
	s.NoError(err)
	expectedCommitment.IncludedInBatch = &batches[0].ID
	treeRoot := s.getAccountTreeRoot()
	expectedCommitment.AccountTreeRoot = &treeRoot
	s.Equal(expectedCommitment, *commitment)

	transfer, err := s.storage.GetCreate2Transfer(tx.Hash)
	s.NoError(err)
	transfer.Signature = tx.Signature
	transfer.ToPublicKey = models.PublicKey{}
	tx.IncludedInCommitment = &commitment.ID
	s.Equal(tx, *transfer)
}

func (s *SyncTestSuite) createAndSubmitTransferBatch(tx *models.Transfer) {
	err := s.storage.AddTransfer(tx)
	s.NoError(err)

	commitments, err := s.transactionExecutor.createTransferCommitments([]models.Transfer{*tx}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	_, err = s.transactionExecutor.submitBatch(txtype.Transfer, commitments)
	s.NoError(err)

	s.client.Commit()
}

func (s *SyncTestSuite) registerAccountOnChain(publicKey *models.PublicKey, expectedPubKeyID uint32) {
	registrations, unsubscribe, err := s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	senderPubKeyID, err := s.client.RegisterAccount(publicKey, registrations)
	s.NoError(err)
	s.Equal(expectedPubKeyID, *senderPubKeyID)
}

func (s *SyncTestSuite) createAndSubmitC2TBatch(tx *models.Create2Transfer) models.Commitment {
	err := s.storage.AddCreate2Transfer(tx)
	s.NoError(err)

	commitments, err := s.transactionExecutor.createCreate2TransferCommitments([]models.Create2Transfer{*tx}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	_, err = s.transactionExecutor.submitBatch(txtype.Create2Transfer, commitments)
	s.NoError(err)

	s.client.Commit()
	return commitments[0]
}

func (s *SyncTestSuite) syncAllBlocks() {
	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)
	err = s.transactionExecutor.SyncBatches(s.stateMutex, 0, *latestBlockNumber)
	s.NoError(err)
}

func mockSignature(t *testing.T) *models.Signature {
	wallet, err := bls.NewRandomWallet(*testDomain)
	require.NoError(t, err)
	signature, err := wallet.Sign(utils.RandomBytes(4))
	require.NoError(t, err)
	return signature.ModelsSignature()
}

func (s *SyncTestSuite) recreateDatabase() {
	err := s.teardown()
	s.NoError(err)
	s.setupDB()
}

func (s *SyncTestSuite) getAccountTreeRoot() common.Hash {
	rawAccountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)
	return common.BytesToHash(rawAccountRoot[:])
}

func (s *SyncTestSuite) setTransferHash(tx *models.Transfer) {
	hash, err := encoder.HashTransfer(tx)
	s.NoError(err)
	tx.Hash = *hash
}

func (s *SyncTestSuite) setCreate2TransferHash(tx *models.Create2Transfer) {
	hash, err := encoder.HashCreate2Transfer(tx)
	s.NoError(err)
	tx.Hash = *hash
}

func TestSyncTestSuite(t *testing.T) {
	suite.Run(t, new(SyncTestSuite))
}
