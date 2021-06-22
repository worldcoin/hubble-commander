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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	transactionExecutor *TransactionExecutor
	transfer            models.Transfer
	wallets             []bls.Wallet
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
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		TxsPerCommitment:       1,
		DevMode:                false,
	}

	s.wallets = generateWallets(s.T(), s.client.ChainState.Rollup, 2)
	s.setupDB()
}

func (s *SyncTestSuite) setupDB() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.tree = st.NewStateTree(s.storage)
	s.transactionExecutor = NewTestTransactionExecutor(s.storage, s.client.Client, s.cfg, TransactionExecutorOpts{AssumeNonces: true})
	err = s.storage.SetChainState(&s.client.ChainState)
	s.NoError(err)

	seedDB(s.T(), s.storage, s.tree, s.wallets)
}

func seedDB(t *testing.T, storage *st.Storage, tree *st.StateTree, wallets []bls.Wallet) {
	err := storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  0,
		PublicKey: *wallets[0].PublicKey(),
	})
	require.NoError(t, err)

	err = storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  1,
		PublicKey: *wallets[1].PublicKey(),
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

func (s *SyncTestSuite) TestSyncBatch_TwoTransferBatches() {
	txs := []models.Transfer{
		{
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(400),
				Fee:         models.MakeUint256(0),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 1,
		}, {
			TransactionBase: models.TransactionBase{
				TxType:      txtype.Transfer,
				FromStateID: 0,
				Amount:      models.MakeUint256(100),
				Fee:         models.MakeUint256(0),
				Nonce:       models.MakeUint256(1),
			},
			ToStateID: 1,
		},
	}
	for i := range txs {
		signTransfer(s.T(), &s.wallets[txs[i].FromStateID], &txs[i])
		s.setTransferHash(&txs[i])
		err := s.storage.AddTransfer(&txs[i])
		s.NoError(err)
	}

	expectedCommitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(expectedCommitments, 2)
	accountRoots := make([]common.Hash, 2)
	for i := range expectedCommitments {
		var pendingBatch *models.Batch
		pendingBatch, err = s.transactionExecutor.NewPendingBatch(txtype.Transfer)
		s.NoError(err)
		err = s.transactionExecutor.SubmitBatch(pendingBatch, []models.Commitment{expectedCommitments[i]})
		s.NoError(err)
		s.client.Commit()

		accountRoots[i] = s.getAccountTreeRoot()
	}

	s.recreateDatabase()
	s.syncAllBatches()

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

func (s *SyncTestSuite) TestSyncBatch_PendingBatch() {
	accountRoot := s.getAccountTreeRoot()
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 1,
	}
	s.setTransferHash(&tx)
	signTransfer(s.T(), &s.wallets[tx.FromStateID], &tx)
	s.createAndSubmitTransferBatch(&tx)

	pendingBatch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)
	s.Nil(pendingBatch.Hash)
	s.Nil(pendingBatch.FinalisationBlock)
	s.Nil(pendingBatch.AccountTreeRoot)

	s.syncAllBatches()

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
	s.NotNil(batches[0].Hash)
	s.NotNil(batches[0].FinalisationBlock)

	s.Equal(accountRoot, *batches[0].AccountTreeRoot)
}

func (s *SyncTestSuite) TestSyncBatch_Create2TransferBatch() {
	tx := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Create2Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID:   ref.Uint32(5),
		ToPublicKey: *s.wallets[0].PublicKey(),
	}
	s.setCreate2TransferHash(&tx)
	signCreate2Transfer(s.T(), &s.wallets[tx.FromStateID], &tx)
	expectedCommitment := s.createAndSubmitC2TBatch(&tx)

	s.recreateDatabase()
	s.syncAllBatches()

	state0, err := s.storage.GetStateLeaf(0)
	s.NoError(err)
	s.Equal(models.MakeUint256(600), state0.Balance)

	state5, err := s.storage.GetStateLeaf(5)
	s.NoError(err)
	s.Equal(models.MakeUint256(400), state5.Balance)
	s.Equal(uint32(0), state5.PubKeyID)

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
	tx.IncludedInCommitment = &commitment.ID
	s.Equal(tx, *transfer)
}

func (s *SyncTestSuite) TestRevertBatch_RevertsState() {
	initialStateRoot, err := s.tree.Root()
	s.NoError(err)

	signTransfer(s.T(), &s.wallets[s.transfer.FromStateID], &s.transfer)
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
		},
		ToStateID: 1,
	}

	pendingBatches := make([]models.Batch, 2)
	for i := range pendingBatches {
		signTransfer(s.T(), &s.wallets[transfers[i].FromStateID], &transfers[i])
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

	signTransfer(s.T(), &s.wallets[s.transfer.FromStateID], &s.transfer)
	pendingBatch := s.createAndSubmitTransferBatch(&s.transfer)
	s.recreateDatabase()

	localTransfer := s.transfer
	localTransfer.Hash = utils.RandomHash()
	localTransfer.Amount = models.MakeUint256(200)
	localTransfer.Fee = models.MakeUint256(10)
	signTransfer(s.T(), &s.wallets[localTransfer.FromStateID], &localTransfer)
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

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *SyncTestSuite) createTransferBatch(tx *models.Transfer) *models.Batch {
	err := s.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
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

func (s *SyncTestSuite) createAndSubmitC2TBatch(tx *models.Create2Transfer) models.Commitment {
	err := s.storage.AddCreate2Transfer(tx)
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Create2Transfer)
	s.NoError(err)
	err = s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return commitments[0]
}

func (s *SyncTestSuite) syncAllBatches() {
	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	newRemoteBatches, err := s.client.GetBatches(&bind.FilterOpts{
		Start: 0,
		End:   latestBlockNumber,
	})
	s.NoError(err)

	for i := range newRemoteBatches {
		remoteBatch := &newRemoteBatches[i]
		err = s.transactionExecutor.SyncBatch(remoteBatch)
		s.NoError(err)
	}
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

func generateWallets(t *testing.T, rollupAddress common.Address, walletsAmount int) []bls.Wallet {
	domain, err := bls.DomainFromBytes(crypto.Keccak256(rollupAddress.Bytes()))
	require.NoError(t, err)

	wallets := make([]bls.Wallet, 0, walletsAmount)
	for i := 0; i < walletsAmount; i++ {
		wallet, err := bls.NewRandomWallet(*domain)
		require.NoError(t, err)
		wallets = append(wallets, *wallet)
	}
	return wallets
}

func TestSyncTestSuite(t *testing.T) {
	suite.Run(t, new(SyncTestSuite))
}
