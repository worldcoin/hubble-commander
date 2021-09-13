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
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd          *Commander
	testClient   *eth.TestClient
	testStorage  *st.TestStorage
	executionCtx *executor.ExecutionContext
	cfg          *config.Config
	wallets      []bls.Wallet
}

func (s *BatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinCommitmentsPerBatch = 1
	s.cfg.Rollup.MaxCommitmentsPerBatch = 32
	s.cfg.Rollup.MinTxsPerCommitment = 1
	s.cfg.Rollup.MaxTxsPerCommitment = 1
	s.cfg.Rollup.DisableSignatures = false
}

func (s *BatchesTestSuite) SetupTest() {
	var err error
	s.testStorage, err = st.NewTestStorage()
	s.NoError(err)
	s.testClient, err = eth.NewConfiguredTestClient(rollup.DeploymentConfig{
		Params: rollup.Params{
			MaxTxsPerCommit: models.NewUint256(1),
		},
	}, eth.ClientConfig{})
	s.NoError(err)

	s.cmd = NewCommander(s.cfg, nil)
	s.cmd.client = s.testClient.Client
	s.cmd.storage = s.testStorage.Storage
	s.cmd.workersContext, s.cmd.stopWorkers = context.WithCancel(context.Background())

	s.executionCtx = executor.NewTestExecutionContext(s.testStorage.Storage, s.testClient.Client, s.cfg.Rollup)

	err = s.cmd.addGenesisBatch()
	s.NoError(err)

	domain, err := s.testClient.GetDomain()
	s.NoError(err)
	s.wallets = generateWallets(s.T(), *domain, 2)
	seedDB(s.T(), s.testStorage.Storage, s.wallets)
}

func (s *BatchesTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.testClient.Close()
	err := s.testStorage.Teardown()
	s.NoError(err)
}

func (s *BatchesTestSuite) TestUnsafeSyncBatches_DoesNotSyncExistingBatchTwice() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	signTransfer(s.T(), &s.wallets[tx.FromStateID], &tx)
	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	s.submitTransferBatch(clonedStorage.Storage, executionCtx, &tx)
	teardown(s.Assertions, clonedStorage.Teardown)

	s.syncAllBlocks()

	tx2 := testutils.MakeTransfer(1, 0, 0, 100)
	signTransfer(s.T(), &s.wallets[tx2.FromStateID], &tx2)
	clonedStorage, executionCtx = cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitTransferBatch(clonedStorage.Storage, executionCtx, &tx2)

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	s.syncAllBlocks()

	batches, err = s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 3)

	state0, err := s.cmd.storage.StateTree.Leaf(0)
	s.NoError(err)
	s.Equal(models.MakeUint256(710), state0.Balance)

	state1, err := s.cmd.storage.StateTree.Leaf(1)
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

	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitTransferBatch(clonedStorage.Storage, executionCtx, &transfers[0])

	s.createTransferBatch(&transfers[1])

	batches, err := s.testClient.GetAllBatches()
	s.NoError(err)
	s.Len(batches, 1)

	s.testClient.Account = s.testClient.Accounts[1]
	err = s.cmd.syncRemoteBatch(&batches[0])
	s.NoError(err)

	batch, err := s.cmd.storage.GetBatch(batches[0].ID)
	s.NoError(err)
	s.Equal(batches[0].Batch, *batch)

	expectedCommitment := models.Commitment{
		ID: models.CommitmentID{
			BatchID:      batch.ID,
			IndexInBatch: 0,
		},
		Type:              txtype.Transfer,
		Transactions:      batches[0].Commitments[0].Transactions,
		FeeReceiver:       batches[0].Commitments[0].FeeReceiver,
		CombinedSignature: batches[0].Commitments[0].CombinedSignature,
		PostStateRoot:     batches[0].Commitments[0].StateRoot,
	}
	commitment, err := s.cmd.storage.GetCommitment(&expectedCommitment.ID)
	s.NoError(err)
	s.Equal(expectedCommitment, *commitment)

	expectedTx := transfers[0]
	expectedTx.Signature = models.Signature{}
	expectedTx.CommitmentID = &commitment.ID
	transfer, err := s.cmd.storage.GetTransfer(transfers[0].Hash)
	s.NoError(err)
	s.Equal(expectedTx, *transfer)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithTooManyTxs() {
	transfer := testutils.MakeTransfer(0, 1, 0, 50)
	s.submitTransferBatch(s.testStorage.Storage, s.executionCtx, &transfer)

	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)

	transfer = testutils.MakeTransfer(0, 1, 1, 100)
	s.submitInvalidTransferBatch(clonedStorage.Storage, executionCtx, &transfer, func(commitment *models.Commitment) {
		commitment.Transactions = append(commitment.Transactions, commitment.Transactions...)
	})

	remoteBatches, err := s.testClient.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.cmd.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.cmd.syncRemoteBatch(&remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithInvalidPostStateRoot() {
	transfer := testutils.MakeTransfer(0, 1, 0, 50)
	s.submitTransferBatch(s.testStorage.Storage, s.executionCtx, &transfer)

	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)

	invalidTransfer := testutils.MakeTransfer(0, 1, 1, 100)
	s.submitInvalidTransferBatch(clonedStorage.Storage, executionCtx, &invalidTransfer, func(commitment *models.Commitment) {
		commitment.PostStateRoot = utils.RandomHash()
	})

	remoteBatches, err := s.testClient.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.cmd.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.cmd.syncRemoteBatch(&remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithInvalidSignature() {
	s.registerAccounts([]uint32{0, 1})

	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)

	invalidTransfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidTransferBatch(clonedStorage.Storage, executionCtx, &invalidTransfer, func(commitment *models.Commitment) {
		commitment.CombinedSignature = models.Signature{1, 2, 3}
	})

	remoteBatches, err := s.testClient.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(&remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[0].ID)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_RemovesExistingBatchAndDisputesFraudulentOne() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 250),
		testutils.MakeTransfer(0, 1, 1, 100),
	}

	s.submitTransferBatch(s.testStorage.Storage, s.executionCtx, &transfers[0])

	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitInvalidTransferBatch(clonedStorage.Storage, executionCtx, &transfers[1], func(commitment *models.Commitment) {
		commitment.Transactions = append(commitment.Transactions, commitment.Transactions...)
	})

	localBatch := s.createTransferBatch(&transfers[2])

	remoteBatches, err := s.testClient.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.cmd.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	s.testClient.Account = s.testClient.Accounts[1]
	err = s.cmd.syncRemoteBatch(&remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
	_, err = s.cmd.storage.GetBatch(localBatch.ID)
	s.True(st.IsNotFoundError(err))
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_DisputesFraudulentCommitmentAfterGenesisOne() {
	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)

	invalidTransfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidTransferBatch(clonedStorage.Storage, executionCtx, &invalidTransfer, func(commitment *models.Commitment) {
		commitment.Transactions = append(commitment.Transactions, commitment.Transactions...)
	})

	remoteBatches, err := s.testClient.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(&remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[0].ID)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithInvalidFeeReceiverTokenID() {
	_, err := s.testStorage.StateTree.Set(2, &models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(2),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 100),
	}

	s.submitTransferBatch(s.testStorage.Storage, s.executionCtx, &transfers[0])

	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitInvalidTransferBatch(clonedStorage.Storage, executionCtx, &transfers[1], func(commitment *models.Commitment) {
		commitment.FeeReceiver = 2
	})

	remoteBatches, err := s.testClient.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.cmd.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.cmd.syncRemoteBatch(&remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithoutTransfersAndInvalidPostStateRoot() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 100),
	}

	s.submitTransferBatch(s.testStorage.Storage, s.executionCtx, &transfers[0])

	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitInvalidTransferBatch(clonedStorage.Storage, executionCtx, &transfers[1], func(commitment *models.Commitment) {
		commitment.Transactions = []byte{}
	})

	remoteBatches, err := s.testClient.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.cmd.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.cmd.syncRemoteBatch(&remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithNotExistingSender() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 100),
	}

	s.submitTransferBatch(s.testStorage.Storage, s.executionCtx, &transfers[0])

	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitInvalidTransferBatch(clonedStorage.Storage, executionCtx, &transfers[1], func(commitment *models.Commitment) {
		transfers[1].FromStateID = 10
		encodedTx, err := encoder.EncodeTransferForCommitment(&transfers[1])
		s.NoError(err)
		commitment.Transactions = encodedTx
	})

	remoteBatches, err := s.testClient.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.cmd.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.cmd.syncRemoteBatch(&remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
}

func (s *BatchesTestSuite) TestSyncRemoteBatch_AllowsNonexistentReceiver() {
	transfer := testutils.MakeTransfer(0, 2, 0, 100)
	s.setTransferHashAndSign(&transfer)

	encodedTx, err := encoder.EncodeTransferForCommitment(&transfer)
	s.NoError(err)

	stateRoot := common.HexToHash("0x09de852e52fff821a7384b6bce2d5c51e9f0d32484e14c2fa29fb140d54ae8e8")

	batch := &eth.DecodedBatch{
		Batch: models.Batch{
			ID:              models.MakeUint256(1),
			Type:            txtype.Transfer,
			TransactionHash: common.Hash{1, 2, 3},
		},
		Commitments: []encoder.DecodedCommitment{{
			StateRoot:         stateRoot,
			CombinedSignature: *s.getTransferCombinedSignature(&transfer),
			FeeReceiver:       0,
			Transactions:      encodedTx,
		}},
	}

	err = s.cmd.syncRemoteBatch(batch)
	s.NoError(err)

	expectedReceiverState := models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  transfer.Amount,
		Nonce:    models.MakeUint256(0),
	}

	leaf, err := s.testStorage.StateTree.Leaf(transfer.ToStateID)
	s.NoError(err)
	s.Equal(expectedReceiverState, leaf.UserState)
}

func (s *BatchesTestSuite) TestUnsafeSyncBatches_SyncsBatchesBeforeInvalidOne() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 250),
		testutils.MakeTransfer(0, 1, 2, 100),
	}

	s.submitTransferBatch(s.testStorage.Storage, s.executionCtx, &transfers[0])

	clonedStorage, executionCtx := cloneStorage(s.Assertions, s.cfg, s.testStorage, s.testClient.Client)
	defer teardown(s.Assertions, clonedStorage.Teardown)

	invalidBatch := s.submitTransferBatch(clonedStorage.Storage, executionCtx, &transfers[1])
	s.submitTransferBatch(clonedStorage.Storage, executionCtx, &transfers[2])

	s.cmd.invalidBatchID = &invalidBatch.ID

	s.syncAllBlocks()

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)
	s.EqualValues(1, batches[1].ID.Uint64())
}

func (s *BatchesTestSuite) syncAllBlocks() {
	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	err = s.cmd.unsafeSyncBatches(0, *latestBlockNumber)
	s.NoError(err)
}

// Make sure that the commander and the execution context uses the same storage
func (s *BatchesTestSuite) submitTransferBatch(
	storage *st.Storage,
	executionCtx *executor.ExecutionContext,
	tx *models.Transfer,
) *models.Batch {
	err := storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := executionCtx.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	domain, err := s.testClient.GetDomain()
	s.NoError(err)
	commitments, err := executionCtx.CreateTransferCommitments(domain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = executionCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.testClient.Commit()
	return pendingBatch
}

// Make sure that the commander and the execution context uses the same storage
func (s *BatchesTestSuite) createTransferBatch(tx *models.Transfer) *models.Batch {
	err := s.cmd.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	domain, err := s.testClient.GetDomain()
	s.NoError(err)
	commitments, err := s.executionCtx.CreateTransferCommitments(domain)
	s.NoError(err)
	s.Len(commitments, 1)
	err = s.cmd.storage.AddCommitment(&commitments[0])
	s.NoError(err)

	pendingBatch.TransactionHash = utils.RandomHash()
	err = s.cmd.storage.AddBatch(pendingBatch)
	s.NoError(err)

	return pendingBatch
}

// Make sure that the commander and the execution context uses the same storage
func (s *BatchesTestSuite) submitInvalidTransferBatch(
	storage *st.Storage,
	executionCtx *executor.ExecutionContext,
	tx *models.Transfer,
	modifier func(commitment *models.Commitment),
) *models.Batch {
	err := storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := executionCtx.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	domain, err := s.testClient.GetDomain()
	s.NoError(err)
	commitments, err := executionCtx.CreateTransferCommitments(domain)
	s.NoError(err)
	s.Len(commitments, 1)

	modifier(&commitments[0])

	err = executionCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.testClient.Commit()
	return pendingBatch
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
	s.Equal(eth.MsgInvalidBatchID, err.Error())

	batch, err := s.cmd.storage.GetBatch(batchID)
	s.Nil(batch)
	s.True(st.IsNotFoundError(err))
}

func (s *BatchesTestSuite) registerAccounts(pubKeyIDs []uint32) {
	registrations, unsubscribe, err := s.testClient.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	for i := range pubKeyIDs {
		leaf, err := s.testStorage.AccountTree.Leaf(pubKeyIDs[i])
		s.NoError(err)

		pubKeyID, err := s.testClient.RegisterAccount(&leaf.PublicKey, registrations)
		s.NoError(err)
		s.Equal(pubKeyIDs[i], *pubKeyID)
	}
}

func (s *BatchesTestSuite) getTransferCombinedSignature(transfer *models.Transfer) *models.Signature {
	domain, err := s.testClient.GetDomain()
	s.NoError(err)
	sig, err := executor.CombineSignatures(models.MakeTransferArray(*transfer), domain)
	s.NoError(err)
	return sig
}

func cloneStorage(
	s *require.Assertions,
	cfg *config.Config,
	storage *st.TestStorage,
	client *eth.Client,
) (*st.TestStorage, *executor.ExecutionContext) {
	clonedStorage, err := storage.Clone()
	s.NoError(err)

	executionCtx := executor.NewTestExecutionContext(clonedStorage.Storage, client, cfg.Rollup)

	return clonedStorage, executionCtx
}

func teardown(s *require.Assertions, teardown func() error) {
	err := teardown()
	s.NoError(err)
}

func TestBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(BatchesTestSuite))
}
