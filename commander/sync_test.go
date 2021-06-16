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
	transfer            models.Transfer
}

func (s *SyncTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.transfer = models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   mockSignature(s.Assertions),
		},
		ToStateID: 1,
	}
	s.setTransferHash(&s.transfer)
}

func (s *SyncTestSuite) SetupTest() {
	var err error
	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch:    1,
		MaxCommitmentsPerBatch:    32,
		TxsPerCommitment:          1,
		PendingTxsCountMultiplier: 1,
	}

	s.setupDB()
}

func (s *SyncTestSuite) setupDB() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.tree = st.NewStateTree(s.storage)
	s.transactionExecutor = newTestTransactionExecutor(s.storage, s.client.Client, s.cfg, transactionExecutorOpts{AssumeNonces: true})

	seedDB(s.T(), s.storage, s.tree)
}

func seedDB(t *testing.T, storage *st.Storage, tree *st.StateTree) {
	err := storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  0,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	require.NoError(t, err)

	err = storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  1,
		PublicKey: models.PublicKey{2, 3, 4},
	})
	require.NoError(t, err)

	err = tree.Set(0, &models.UserState{
		PubKeyID:   0,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(0),
	})
	require.NoError(t, err)

	err = tree.Set(1, &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	})
	require.NoError(t, err)
}

func (s *SyncTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *SyncTestSuite) TestSyncBatches_TwoTransferBatches() {
	txs := []models.Transfer{
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(400),
				Fee:         models.MakeUint256(0),
				Nonce:       models.MakeUint256(0),
				Signature:   mockSignature(s.Assertions),
			},
			ToStateID: 1,
		}, {
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(100),
				Fee:         models.MakeUint256(0),
				Nonce:       models.MakeUint256(1),
				Signature:   mockSignature(s.Assertions),
			},
			ToStateID: 1,
		},
	}
	for i := range txs {
		s.setTransferHash(&txs[i])
		err := s.storage.AddTransfer(&txs[i])
		s.NoError(err)
	}

	expectedCommitments, err := s.transactionExecutor.createTransferCommitments([]models.Transfer{txs[0]}, testDomain)
	s.NoError(err)
	s.Len(expectedCommitments, 2)
	accountRoots := make([]common.Hash, 2)
	for i := range expectedCommitments {
		var pendingBatch *models.Batch
		pendingBatch, err = newPendingBatch(s.storage, txtype.Transfer)
		s.NoError(err)
		err = s.transactionExecutor.submitBatch(pendingBatch, []models.Commitment{expectedCommitments[i]})
		s.NoError(err)
		s.client.Commit()

		accountRoots[i] = s.getAccountTreeRoot()
	}

	s.recreateDatabase()
	s.syncAllBlocks()

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)
	s.Equal(models.MakeUint256(1), batches[0].ID)
	s.Equal(models.MakeUint256(2), batches[1].ID)
	s.Equal(accountRoots[0], *batches[0].AccountTreeRoot)
	s.Equal(accountRoots[1], *batches[1].AccountTreeRoot)

	for i := range expectedCommitments {
		commitment, err := s.storage.GetCommitment(expectedCommitments[i].ID)
		s.NoError(err)
		expectedCommitments[i].IncludedInBatch = &batches[i].ID
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
			Signature:   mockSignature(s.Assertions),
		},
		ToStateID: 1,
	}
	s.createAndSubmitTransferBatch(&tx)

	s.recreateDatabase()
	s.syncAllBlocks()

	// Begin database transaction
	var err error
	s.transactionExecutor, err = newTransactionExecutor(s.storage, s.client.Client, s.cfg, transactionExecutorOpts{})
	s.NoError(err)

	tx2 := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   mockSignature(s.Assertions),
		},
		ToStateID: 0,
	}
	s.createAndSubmitTransferBatch(&tx2)

	batches, err := s.transactionExecutor.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	// Rollback changes to the database
	s.transactionExecutor.Rollback(nil)
	s.transactionExecutor = newTestTransactionExecutor(s.storage, s.client.Client, s.cfg, transactionExecutorOpts{})

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
			Signature:   mockSignature(s.Assertions),
		},
		ToStateID: 1,
	}
	s.setTransferHash(&tx)
	s.createAndSubmitTransferBatch(&tx)

	pendingBatch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)
	s.Nil(pendingBatch.Hash)
	s.Nil(pendingBatch.FinalisationBlock)
	s.Nil(pendingBatch.AccountTreeRoot)

	s.syncAllBlocks()

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
	s.NotNil(batches[0].Hash)
	s.NotNil(batches[0].FinalisationBlock)

	s.Equal(accountRoot, *batches[0].AccountTreeRoot)
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
			Signature:   mockSignature(s.Assertions),
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

	treeRoot := s.getAccountTreeRoot()
	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
	s.Equal(treeRoot, *batches[0].AccountTreeRoot)

	commitment, err := s.storage.GetCommitment(expectedCommitment.ID)
	s.NoError(err)
	expectedCommitment.IncludedInBatch = &batches[0].ID
	s.Equal(expectedCommitment, *commitment)

	transfer, err := s.storage.GetCreate2Transfer(tx.Hash)
	s.NoError(err)
	transfer.Signature = tx.Signature
	transfer.ToPublicKey = models.PublicKey{}
	tx.IncludedInCommitment = &commitment.ID
	s.Equal(tx, *transfer)
}

func (s *SyncTestSuite) TestRevertBatch_RevertsState() {
	initialStateRoot, err := s.tree.Root()
	s.NoError(err)

	pendingBatch := s.createAndSubmitTransferBatch(&s.transfer)
	decodedBatch := &eth.DecodedBatch{
		Batch: models.Batch{
			ID:                models.MakeUint256(1),
			Type:              txtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(20),
		},
	}
	err = s.transactionExecutor.revertBatches(decodedBatch, pendingBatch)
	s.NoError(err)

	stateRoot, err := s.tree.Root()
	s.NoError(err)
	s.Equal(*initialStateRoot, *stateRoot)

	state0, err := s.storage.GetStateLeaf(s.transfer.FromStateID)
	s.NoError(err)
	s.Equal(uint64(1000), state0.Balance.Uint64())
	state1, err := s.storage.GetStateLeaf(s.transfer.ToStateID)
	s.NoError(err)
	s.Equal(uint64(0), state1.Balance.Uint64())

	transfer, err := s.storage.GetTransfer(s.transfer.Hash)
	s.NoError(err)
	s.Nil(transfer.IncludedInCommitment)
}

func (s *SyncTestSuite) TestRevertBatch_DeletesCommitmentsAndBatches() {
	initialStateRoot, err := s.tree.Root()
	s.NoError(err)

	transfers := make([]models.Transfer, 2)
	transfers[0] = s.transfer
	transfers[1] = models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(200),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(1),
			Signature:   mockSignature(s.Assertions),
		},
		ToStateID: 1,
	}

	pendingBatches := make([]models.Batch, 2)
	for i := range pendingBatches {
		pendingBatches[i] = *s.createAndSubmitTransferBatch(&transfers[i])
	}

	decodedBatch := &eth.DecodedBatch{
		Batch: models.Batch{
			ID:                models.MakeUint256(1),
			Type:              txtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(133742069),
		},
	}
	err = s.transactionExecutor.revertBatches(decodedBatch, &pendingBatches[0])
	s.NoError(err)

	stateRoot, err := s.tree.Root()
	s.NoError(err)
	s.Equal(*initialStateRoot, *stateRoot)

	syncedBatch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)
	s.Equal(decodedBatch.TransactionHash, syncedBatch.TransactionHash)
	s.Equal(decodedBatch.Hash, syncedBatch.Hash)
	s.Equal(decodedBatch.FinalisationBlock, syncedBatch.FinalisationBlock)

	_, err = s.storage.GetBatch(models.MakeUint256(2))
	s.Equal(st.NewNotFoundError("batch"), err)
}

func (s *SyncTestSuite) TestRevertBatch_SyncsCorrectBatch() {
	startBlock, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	pendingBatch := s.createAndSubmitTransferBatch(&s.transfer)
	s.recreateDatabase()

	localTransfer := s.transfer
	localTransfer.Hash = utils.RandomHash()
	localTransfer.Amount = models.MakeUint256(200)
	localTransfer.Fee = models.MakeUint256(10)
	localBatch := s.createTransferBatch(&localTransfer)

	batches, err := s.client.GetBatches(&bind.FilterOpts{Start: *startBlock})
	s.NoError(err)
	s.Len(batches, 1)

	err = s.transactionExecutor.revertBatches(&batches[0], localBatch)
	s.NoError(err)

	batch, err := s.storage.GetBatch(pendingBatch.ID)
	s.NoError(err)
	s.Equal(batches[0].Batch, *batch)

	expectedCommitment := models.Commitment{
		ID:                2,
		Type:              txtype.Transfer,
		Transactions:      batches[0].Commitments[0].Transactions,
		FeeReceiver:       batches[0].Commitments[0].FeeReceiver,
		CombinedSignature: batches[0].Commitments[0].CombinedSignature,
		PostStateRoot:     batches[0].Commitments[0].StateRoot,
		IncludedInBatch:   &batch.ID,
	}
	commitment, err := s.storage.GetCommitment(2)
	s.NoError(err)
	s.Equal(expectedCommitment, *commitment)

	expectedTx := s.transfer
	expectedTx.Signature = models.Signature{}
	expectedTx.IncludedInCommitment = &commitment.ID
	transfer, err := s.storage.GetTransfer(s.transfer.Hash)
	s.NoError(err)
	s.Equal(expectedTx, *transfer)
}

func (s *SyncTestSuite) createAndSubmitTransferBatch(tx *models.Transfer) *models.Batch {
	err := s.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := newPendingBatch(s.storage, txtype.Transfer)
	s.NoError(err)

	commitments, err := s.transactionExecutor.createTransferCommitments([]models.Transfer{*tx}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = s.transactionExecutor.submitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *SyncTestSuite) createTransferBatch(tx *models.Transfer) *models.Batch {
	err := s.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := newPendingBatch(s.storage, txtype.Transfer)
	s.NoError(err)

	commitments, err := s.transactionExecutor.createTransferCommitments([]models.Transfer{*tx}, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	pendingBatch.TransactionHash = utils.RandomHash()
	err = s.storage.AddBatch(pendingBatch)
	s.NoError(err)

	err = s.transactionExecutor.markCommitmentsAsIncluded(commitments, pendingBatch.ID)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
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

	pendingBatch, err := newPendingBatch(s.storage, txtype.Create2Transfer)
	s.NoError(err)
	err = s.transactionExecutor.submitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return commitments[0]
}

func (s *SyncTestSuite) syncAllBlocks() {
	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)
	err = s.transactionExecutor.SyncBatches(0, *latestBlockNumber)
	s.NoError(err)
}

func mockSignature(s *require.Assertions) models.Signature {
	wallet, err := bls.NewRandomWallet(*testDomain)
	s.NoError(err)
	signature, err := wallet.Sign(utils.RandomBytes(4))
	s.NoError(err)
	return *signature.ModelsSignature()
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
