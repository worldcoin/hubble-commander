package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/result"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxsBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd     *Commander
	client  *eth.TestClient
	storage *st.TestStorage
	txsCtx  *executor.TxsContext
	cfg     *config.Config
	wallets []bls.Wallet
}

func (s *TxsBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinTxsPerCommitment = 1
	s.cfg.Rollup.MaxTxsPerCommitment = 32
	s.cfg.Rollup.DisableSignatures = false
}

func (s *TxsBatchesTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.client = newClientWithGenesisState(s.T(), s.storage)

	s.cmd = NewCommander(s.cfg, s.client.Blockchain)
	s.cmd.client = s.client.Client
	s.cmd.storage = s.storage.Storage
	s.cmd.metrics = metrics.NewCommanderMetrics()
	s.cmd.workersContext, s.cmd.stopWorkersContext = context.WithCancel(context.Background())

	executionCtx := executor.NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg.Rollup)
	s.txsCtx = executor.NewTestTxsContext(executionCtx, batchtype.Transfer)

	err = s.cmd.addGenesisBatch()
	s.NoError(err)

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)
	setAccountLeaves(s.T(), s.storage.Storage, s.wallets)
}

func (s *TxsBatchesTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TxsBatchesTestSuite) TestUnsafeSyncBatches_DoesNotSyncExistingBatchTwice() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	signTransfer(s.T(), &s.wallets[tx.FromStateID], &tx)
	s.submitBatchInTx(&tx)

	s.syncAllBlocks()

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	tx2 := testutils.MakeTransfer(1, 0, 0, 100)
	signTransfer(s.T(), &s.wallets[tx2.FromStateID], &tx2)
	s.submitBatchInTx(&tx2)

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

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_ReplaceLocalBatchWithRemoteOne() {
	initialStateRoot, err := s.cmd.storage.StateTree.Root()
	s.NoError(err)

	remoteTx := testutils.MakeTransfer(0, 1, 0, 100)
	s.setTransferHashAndSign(&remoteTx)
	s.submitBatchInTx(&remoteTx)

	localTx := testutils.MakeTransfer(0, 1, 0, 200)
	s.setTransferHashAndSign(&localTx)
	s.createTransferBatchLocally(&localTx)

	batches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(batches, 1)
	remoteBatch := batches[0]

	err = s.cmd.syncRemoteBatch(remoteBatch)
	s.NoError(err)

	// Correct batch stored
	expectedBatch := remoteBatch.ToBatch(*initialStateRoot)
	storedBatch, err := s.cmd.storage.GetBatch(remoteBatch.GetID())
	s.NoError(err)
	s.Equal(*expectedBatch, *storedBatch)

	// Correct commitment stored
	remoteTxBatch := remoteBatch.ToDecodedTxBatch()
	remoteCommitment := remoteTxBatch.Commitments[0].ToDecodedCommitment()
	expectedCommitment := models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      remoteTxBatch.ID,
				IndexInBatch: 0,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: remoteCommitment.StateRoot,
		},
		FeeReceiver:       remoteCommitment.FeeReceiver,
		CombinedSignature: remoteCommitment.CombinedSignature,
		BodyHash:          remoteTxBatch.Commitments[0].BodyHash(remoteTxBatch.AccountTreeRoot),
	}
	storedCommitment, err := s.cmd.storage.GetTxCommitment(&expectedCommitment.ID)
	s.NoError(err)
	s.Equal(expectedCommitment, *storedCommitment)

	// Correct tx stored
	expectedTx := remoteTx
	expectedTx.Signature = models.Signature{}
	expectedTx.CommitmentID = &expectedCommitment.ID
	transfer, err := s.cmd.storage.GetTransfer(remoteTx.Hash)
	s.NoError(err)
	s.Equal(expectedTx, *transfer)

	// Previously stored tx moved back to mempool
	pendingTransfers, err := s.cmd.storage.GetPendingTransfers()
	s.NoError(err)
	s.Len(pendingTransfers, 1)
	s.Equal(localTx, pendingTransfers[0])
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithTooManyTxs() {
	tx := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&tx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		commitment.Transactions = append(commitment.Transactions, make([]byte, 32*encoder.TransferLength)...)
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.TooManyTx, getDisputeResult(s.Assertions, s.client))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithInvalidPostStateRoot() {
	tx := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&tx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		commitment.PostStateRoot = utils.RandomHash()
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.Ok, getDisputeResult(s.Assertions, s.client)) // invalid post state root emits Result.Ok
}

// This test checks that state witnesses needed for dispute tx are gathered correctly in case of a self transfer
func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesFraudulentBatchWithSelfTransfer() {
	tx := testutils.MakeTransfer(0, 0, 0, 100)
	s.submitInvalidBatchInTx(&tx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		commitment.PostStateRoot = common.Hash{1, 2, 3, 4}
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())

	// post state root is checked after applying the transfer in the SC dispute method, so we know that the witnesses were correct
	s.Equal(result.Ok, getDisputeResult(s.Assertions, s.client)) // invalid post state root emits Result.Ok
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithInvalidSignature() {
	s.registerAccounts([]uint32{0, 1})

	tx := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&tx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		commitment.CombinedSignature[0] = 0 // break the signature
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.BadSignature, getDisputeResult(s.Assertions, s.client))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithSignatureInBadFormat() {
	s.registerAccounts([]uint32{0, 1})

	tx := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&tx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		commitment.CombinedSignature = models.Signature{1, 2, 3}
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.BadPrecompileCall, getDisputeResult(s.Assertions, s.client))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_RemovesExistingBatchAndDisputesFraudulentOne() {
	remoteTx := testutils.MakeTransfer(0, 1, 0, 250)
	s.submitInvalidBatchInTx(&remoteTx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		commitment.Transactions = append(commitment.Transactions, 1, 2, 3, 4)
	})

	localTx := testutils.MakeTransfer(0, 1, 0, 100)
	localBatch := s.createTransferBatchLocally(&localTx)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)
	remoteBatch := remoteBatches[0]

	err = s.cmd.syncRemoteBatch(remoteBatch)
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatch.GetID())
	s.Equal(result.BadCompression, getDisputeResult(s.Assertions, s.client))

	_, err = s.cmd.storage.GetBatch(localBatch.ID)
	s.True(st.IsNotFoundError(err))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesFraudulentCommitmentAfterGenesisOne() {
	tx := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&tx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		commitment.Transactions = append(commitment.Transactions, 1, 2, 3, 4)
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.BadCompression, getDisputeResult(s.Assertions, s.client))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithInvalidFeeReceiverTokenID() {
	tx := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&tx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		commitment.FeeReceiver = 2
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.BadToTokenID, getDisputeResult(s.Assertions, s.client))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithoutTransfersAndInvalidPostStateRoot() {
	tx := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&tx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		commitment.Transactions = []byte{}
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.Ok, getDisputeResult(s.Assertions, s.client)) // invalid post state root emits Result.Ok
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithNonexistentSender() {
	tx := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&tx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		tx.FromStateID = 10
		encodedTx, err := encoder.EncodeTransferForCommitment(&tx)
		s.NoError(err)
		commitment.Transactions = encodedTx
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.NotEnoughTokenBalance, getDisputeResult(s.Assertions, s.client))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesC2TWithNonRegisteredReceiverPublicKey() {
	// Change batch type used in TxsContext
	executionCtx := executor.NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg.Rollup)
	s.txsCtx = executor.NewTestTxsContext(executionCtx, batchtype.Create2Transfer)

	// Register public keys added to the account tree for signature disputes to work
	s.registerAccounts([]uint32{0, 1})

	tx := testutils.MakeCreate2Transfer(0, nil, 0, 100, s.wallets[0].PublicKey())
	s.submitInvalidBatchInTx(&tx, func(storage *st.Storage, commitment *models.CommitmentWithTxs) {
		// Fix post state root
		_, err := storage.StateTree.Set(3, &models.UserState{
			PubKeyID: 1234,
			TokenID:  models.MakeUint256(0),
			Balance:  tx.Amount,
			Nonce:    models.MakeUint256(0),
		})
		s.NoError(err)
		root, err := storage.StateTree.Root()
		s.NoError(err)
		commitment.PostStateRoot = *root

		// Replace toStateID and toPubKeyID in C2T
		tx.ToStateID = ref.Uint32(3)
		encodedTx, err := encoder.EncodeCreate2TransferForCommitment(&tx, 1234)
		s.NoError(err)
		commitment.Transactions = encodedTx
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.NonexistentReceiver, getDisputeResult(s.Assertions, s.client))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithInvalidTokenAmount() {
	tx := testutils.MakeTransfer(0, 1, 0, 10)

	s.submitInvalidBatchInTx(&tx, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		tx.Amount = models.MakeUint256(0)
		encodedTx, err := encoder.EncodeTransferForCommitment(&tx)
		s.NoError(err)

		commitment.Transactions = encodedTx
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.InvalidTokenAmount, getDisputeResult(s.Assertions, s.client))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_AllowsTransferToNonexistentReceiver() {
	tx := s.submitTransferToNonexistentReceiver()

	s.syncAllBlocks()

	expectedReceiverState := models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  tx.Amount,
		Nonce:    models.MakeUint256(0),
	}

	leaf, err := s.storage.StateTree.Leaf(tx.ToStateID)
	s.NoError(err)
	s.Equal(expectedReceiverState, leaf.UserState)
}

func (s *TxsBatchesTestSuite) submitTransferToNonexistentReceiver() models.Transfer {
	tx := testutils.MakeTransfer(0, 3, 0, 100)
	s.setTransferHashAndSign(&tx)

	txController, txStorage, txsCtx := s.beginTransaction()
	defer txController.Rollback(nil)

	// Temporary set receiver StateLeaf for successful batch submission
	_, err := txStorage.StateTree.Set(tx.ToStateID, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	s.submitBatch(txStorage, txsCtx, &tx)

	return tx
}

func (s *TxsBatchesTestSuite) TestUnsafeSyncBatches_SyncsBatchesBeforeInvalidOne() {
	txs := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 250),
		testutils.MakeTransfer(0, 1, 1, 100),
	}
	s.setTransferHashAndSign(&txs[0])

	txController, txStorage, txsCtx := s.beginTransaction()
	defer txController.Rollback(nil)

	s.submitBatch(txStorage, txsCtx, &txs[0])
	invalidBatch := s.submitBatch(txStorage, txsCtx, &txs[1])

	s.cmd.invalidBatchID = &invalidBatch.ID

	s.syncAllBlocks()

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)
	s.EqualValues(1, batches[1].ID.Uint64())
}

func (s *TxsBatchesTestSuite) syncAllBlocks() {
	latestBlockNumber, err := s.client.GetLatestBlockNumber()
	s.NoError(err)

	err = s.cmd.unsafeSyncBatches(0, *latestBlockNumber)
	s.NoError(err)
}

func (s *TxsBatchesTestSuite) createTransferBatchLocally(tx *models.Transfer) *models.Batch {
	err := s.cmd.storage.AddTransaction(tx)
	s.NoError(err)

	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.Transfer)
	s.NoError(err)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
	err = s.cmd.storage.AddTxCommitment(&batchData.Commitments()[0].TxCommitment)
	s.NoError(err)

	pendingBatch.TransactionHash = utils.RandomHash()
	err = s.cmd.storage.AddBatch(pendingBatch)
	s.NoError(err)

	return pendingBatch
}

func (s *TxsBatchesTestSuite) submitBatch(
	storage *st.Storage,
	txsCtx *executor.TxsContext,
	tx models.GenericTransaction,
) *models.Batch {
	pendingBatch := submitInvalidTxsBatch(s.Assertions, storage, txsCtx, tx,
		func(_ *st.Storage, _ *models.CommitmentWithTxs) {},
	)

	s.client.GetBackend().Commit()
	return pendingBatch
}

// Submits batch on chain without adding it to storage
func (s *TxsBatchesTestSuite) submitBatchInTx(tx models.GenericTransaction) {
	s.submitInvalidBatchInTx(tx, func(_ *st.Storage, _ *models.CommitmentWithTxs) {})
}

// Submits invalid batch on chain without adding it to storage
func (s *TxsBatchesTestSuite) submitInvalidBatchInTx(
	tx models.GenericTransaction,
	modifier func(storage *st.Storage, commitment *models.CommitmentWithTxs),
) {
	txController, txStorage, txsCtx := s.beginTransaction()
	defer txController.Rollback(nil)

	submitInvalidTxsBatch(s.Assertions, txStorage, txsCtx, tx, modifier)

	s.client.GetBackend().Commit()
}

func (s *TxsBatchesTestSuite) beginTransaction() (*db.TxController, *st.Storage, *executor.TxsContext) {
	txController, txStorage := s.storage.BeginTransaction(st.TxOptions{})

	executionCtx := executor.NewTestExecutionContext(txStorage, s.client.Client, s.cfg.Rollup)
	return txController, txStorage, executor.NewTestTxsContext(executionCtx, s.txsCtx.BatchType)
}

func submitInvalidTxsBatch(
	s *require.Assertions,
	storage *st.Storage,
	txsCtx *executor.TxsContext,
	tx models.GenericTransaction,
	modifier func(storage *st.Storage, commitment *models.CommitmentWithTxs),
) *models.Batch {
	err := storage.AddTransaction(tx)
	s.NoError(err)

	pendingBatch, err := txsCtx.NewPendingBatch(txsCtx.BatchType)
	s.NoError(err)

	batchData, err := txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)

	modifier(storage, &batchData.Commitments()[0])

	err = txsCtx.SubmitBatch(pendingBatch, batchData)
	s.NoError(err)

	return pendingBatch
}

func (s *TxsBatchesTestSuite) setTransferHashAndSign(txs ...*models.Transfer) {
	for i := range txs {
		signTransfer(s.T(), &s.wallets[txs[i].FromStateID], txs[i])
		hash, err := encoder.HashTransfer(txs[i])
		s.NoError(err)
		txs[i].Hash = *hash
	}
}

func checkBatchAfterDispute(s *require.Assertions, cmd *Commander, batchID models.Uint256) {
	_, err := cmd.client.GetBatch(&batchID)
	s.Error(err)
	s.Equal(eth.MsgInvalidBatchID, err.Error())

	batch, err := cmd.storage.GetBatch(batchID)
	s.Nil(batch)
	s.True(st.IsNotFoundError(err))
}

func getDisputeResult(s *require.Assertions, client *eth.TestClient) result.DisputeResult {
	it, err := client.Rollup.FilterRollbackTriggered(nil)
	s.NoError(err)
	it.Next()
	s.NoError(it.Error())
	res := result.DisputeResult(it.Event.Result)
	s.False(it.Next())
	return res
}

func (s *TxsBatchesTestSuite) registerAccounts(pubKeyIDs []uint32) {
	for i := range pubKeyIDs {
		leaf, err := s.storage.AccountTree.Leaf(pubKeyIDs[i])
		s.NoError(err)

		pubKeyID, err := s.client.RegisterAccountAndWait(&leaf.PublicKey)
		s.NoError(err)
		s.Equal(pubKeyIDs[i], *pubKeyID)
	}
}

func TestTxsBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(TxsBatchesTestSuite))
}
