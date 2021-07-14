package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd                 *Commander
	testClient          *eth.TestClient
	transactionExecutor *executor.TransactionExecutor
	cfg                 *config.Config
	teardown            func() error
	wallets             []bls.Wallet
}

func (s *BatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinCommitmentsPerBatch = 1
	s.cfg.Rollup.MaxCommitmentsPerBatch = 32
	s.cfg.Rollup.MinTxsPerCommitment = 1
	s.cfg.Rollup.MaxTxsPerCommitment = 1
	s.cfg.Rollup.DevMode = false
}

func (s *BatchesTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.testClient, err = eth.NewConfiguredTestClient(rollup.DeploymentConfig{
		Params: rollup.Params{
			MaxTxsPerCommit: models.NewUint256(1),
		},
	})
	s.NoError(err)
	err = testStorage.SetChainState(&s.testClient.ChainState)
	s.NoError(err)

	s.cmd = NewCommander(s.cfg)
	s.cmd.client = s.testClient.Client
	s.cmd.storage = testStorage.Storage
	s.cmd.stopChannel = make(chan bool)

	s.transactionExecutor = executor.NewTestTransactionExecutor(testStorage.Storage, s.testClient.Client, s.cfg.Rollup, context.Background())

	s.wallets = generateWallets(s.T(), s.testClient.ChainState.Rollup, 2)
	seedDB(s.T(), testStorage.Storage, st.NewStateTree(testStorage.Storage), s.wallets)
}

func (s *BatchesTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.testClient.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *BatchesTestSuite) TestUnsafeSyncBatches_DoesNotSyncExistingBatchTwice() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	signTransfer(s.T(), &s.wallets[tx.FromStateID], &tx)
	createAndSubmitTransferBatch(s.T(), s.cmd, &tx)
	s.testClient.Commit()

	s.syncAllBlocks()

	tx2 := testutils.MakeTransfer(1, 0, 0, 100)
	signTransfer(s.T(), &s.wallets[tx2.FromStateID], &tx2)
	createAndSubmitTransferBatch(s.T(), s.cmd, &tx2)
	s.testClient.Commit()

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)

	s.syncAllBlocks()

	batches, err = s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	state0, err := s.cmd.storage.GetStateLeaf(0)
	s.NoError(err)
	s.Equal(models.MakeUint256(710), state0.Balance)

	state1, err := s.cmd.storage.GetStateLeaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(290), state1.Balance)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_ReplaceLocalBatchWithRemoteOne() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 100),
		testutils.MakeTransfer(0, 1, 0, 200),
	}
	for i := range transfers {
		s.setTransferHashAndSign(&transfers[i])
	}

	s.runInTransaction(func() {
		s.createAndSubmitTransferBatch(&transfers[0])
	})

	s.createTransferBatch(&transfers[1])

	batches, err := s.testClient.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(batches, 1)

	s.testClient.Account = s.testClient.Accounts[1]
	err = s.cmd.syncRemoteBatch(&batches[0])
	s.NoError(err)

	batch, err := s.cmd.storage.GetBatch(batches[0].ID)
	s.NoError(err)
	s.Equal(batches[0].Batch, *batch)

	expectedCommitment := models.Commitment{
		ID:                3,
		Type:              txtype.Transfer,
		Transactions:      batches[0].Commitments[0].Transactions,
		FeeReceiver:       batches[0].Commitments[0].FeeReceiver,
		CombinedSignature: batches[0].Commitments[0].CombinedSignature,
		PostStateRoot:     batches[0].Commitments[0].StateRoot,
		IncludedInBatch:   &batch.ID,
	}
	commitment, err := s.cmd.storage.GetCommitment(3)
	s.NoError(err)
	s.Equal(expectedCommitment, *commitment)

	expectedTx := transfers[0]
	expectedTx.Signature = models.Signature{}
	expectedTx.IncludedInCommitment = &commitment.ID
	transfer, err := s.cmd.storage.GetTransfer(transfers[0].Hash)
	s.NoError(err)
	s.Equal(expectedTx, *transfer)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_DisputesFraudulentBatch() {
	transfer := testutils.MakeTransfer(0, 1, 0, 50)
	s.createAndSubmitTransferBatch(&transfer)

	s.runInTransaction(func() {
		invalidTransfer := testutils.MakeTransfer(0, 1, 1, 100)
		s.createAndSubmitInvalidTransferBatch(&invalidTransfer)
	})

	remoteBatches, err := s.testClient.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.cmd.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.cmd.syncRemoteBatch(&remoteBatches[1])
	s.NoError(err)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_RemovesExistingBatchAndDisputesFraudulentOne() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 250),
		testutils.MakeTransfer(0, 1, 1, 100),
	}

	s.createAndSubmitTransferBatch(&transfers[0])
	s.runInTransaction(func() {
		s.createAndSubmitInvalidTransferBatch(&transfers[1])
	})

	localBatch := s.createTransferBatch(&transfers[2])

	remoteBatches, err := s.testClient.GetBatches(&bind.FilterOpts{})
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.cmd.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	s.testClient.Account = s.testClient.Accounts[1]
	err = s.cmd.syncRemoteBatch(&remoteBatches[1])
	s.NoError(err)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
	_, err = s.cmd.storage.GetBatch(localBatch.ID)
	s.True(st.IsNotFoundError(err))
}

func (s *BatchesTestSuite) syncAllBlocks() {
	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	err = s.cmd.unsafeSyncBatches(0, *latestBlockNumber)
	s.NoError(err)
}

func (s *BatchesTestSuite) createAndSubmitTransferBatch(tx *models.Transfer) *models.Batch {
	_, err := s.cmd.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.testClient.Commit()
	return pendingBatch
}

func (s *BatchesTestSuite) createTransferBatch(tx *models.Transfer) *models.Batch {
	_, err := s.cmd.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	pendingBatch.TransactionHash = utils.RandomHash()
	err = s.cmd.storage.AddBatch(pendingBatch)
	s.NoError(err)

	err = s.cmd.storage.MarkCommitmentAsIncluded(commitments[0].ID, pendingBatch.ID)
	s.NoError(err)

	return pendingBatch
}

func (s *BatchesTestSuite) createAndSubmitInvalidTransferBatch(tx *models.Transfer) *models.Batch {
	_, err := s.cmd.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	commitments[0].Transactions = append(commitments[0].Transactions, commitments[0].Transactions...)

	err = s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.testClient.Commit()
	return pendingBatch
}

func (s *BatchesTestSuite) runInTransaction(handler func()) {
	storage := *s.cmd.storage
	txController, txStorage, err := s.cmd.storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	s.NoError(err)
	s.cmd.storage = txStorage

	defer func() {
		txController.Rollback(nil)
		s.cmd.storage = &storage
		s.transactionExecutor = executor.NewTestTransactionExecutor(s.cmd.storage, s.testClient.Client, s.cfg.Rollup, context.Background())
	}()

	s.transactionExecutor = executor.NewTestTransactionExecutor(s.cmd.storage, s.testClient.Client, s.cfg.Rollup, context.Background())
	handler()
}

func (s *BatchesTestSuite) setTransferHashAndSign(txs ...*models.Transfer) {
	for i := range txs {
		signTransfer(s.T(), &s.wallets[txs[i].FromStateID], txs[i])
		hash, err := encoder.HashTransfer(txs[i])
		s.NoError(err)
		txs[i].Hash = *hash
	}
}

func (s *BatchesTestSuite) checkBatchAfterDispute(batchID models.Uint256) {
	_, err := s.testClient.GetBatch(&batchID)
	s.Error(err)
	s.Equal("execution reverted: Batch id greater than total number of batches, invalid batch id", err.Error())

	batch, err := s.cmd.storage.GetBatch(batchID)
	s.Nil(batch)
	s.True(st.IsNotFoundError(err))
}

func TestBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(BatchesTestSuite))
}
