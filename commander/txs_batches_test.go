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
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
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

	s.cmd = NewCommander(s.cfg, nil)
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
	//seedDB(s.T(), s.storage.Storage, s.wallets)
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
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 100),
		testutils.MakeTransfer(0, 1, 0, 200),
	}
	for i := range transfers {
		s.setTransferHashAndSign(&transfers[i])
	}

	s.submitBatchInTx(&transfers[0])

	root, err := s.cmd.storage.StateTree.Root()
	s.NoError(err)

	//TODO-ref: change to different function
	s.createTransferBatch(&transfers[1])

	batches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(batches, 1)

	s.client.Account = s.client.Accounts[1]
	err = s.cmd.syncRemoteBatch(batches[0])
	s.NoError(err)

	batch, err := s.cmd.storage.GetBatch(batches[0].GetID())
	s.NoError(err)
	s.Equal(*batches[0].ToBatch(*root), *batch)

	txBatch := batches[0].ToDecodedTxBatch()
	decodedCommitment := txBatch.Commitments[0].ToDecodedCommitment()
	expectedCommitment := models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      batch.ID,
				IndexInBatch: 0,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: decodedCommitment.StateRoot,
		},
		FeeReceiver:       decodedCommitment.FeeReceiver,
		CombinedSignature: decodedCommitment.CombinedSignature,
		BodyHash:          txBatch.Commitments[0].BodyHash(*batch.AccountTreeRoot),
	}
	commitment, err := s.cmd.storage.GetTxCommitment(&expectedCommitment.ID)
	s.NoError(err)
	s.Equal(expectedCommitment, *commitment)

	expectedTx := transfers[0]
	expectedTx.Signature = models.Signature{}
	expectedTx.CommitmentID = &commitment.ID
	transfer, err := s.cmd.storage.GetTransfer(transfers[0].Hash)
	s.NoError(err)
	s.Equal(expectedTx, *transfer)
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithTooManyTxs() {
	transfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&transfer, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
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
	transfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&transfer, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
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

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_CanDisputeFraudulentBatchWithTransferToSelf() {
	//TODO-ref: replace with submitInvalidTransferInTx
	s.submitFraudulentBatchWithTransferToSelf()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.NotEnoughTokenBalance, getDisputeResult(s.Assertions, s.client))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithInvalidSignature() {
	s.registerAccounts([]uint32{0, 1})

	transfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&transfer, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
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

	transfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&transfer, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
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
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 250),
		testutils.MakeTransfer(0, 1, 0, 100),
	}

	s.submitInvalidBatchInTx(&transfers[0], func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		commitment.Transactions = append(commitment.Transactions, 1, 2, 3, 4)
	})

	localBatch := s.createTransferBatch(&transfers[1])

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	s.client.Account = s.client.Accounts[1]
	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	checkBatchAfterDispute(s.Assertions, s.cmd, remoteBatches[0].GetID())
	s.Equal(result.BadCompression, getDisputeResult(s.Assertions, s.client))
	_, err = s.cmd.storage.GetBatch(localBatch.ID)
	s.True(st.IsNotFoundError(err))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesFraudulentCommitmentAfterGenesisOne() {
	transfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&transfer, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
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
	transfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&transfer, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
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
	transfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&transfer, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
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
	transfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatchInTx(&transfer, func(_ *st.Storage, commitment *models.CommitmentWithTxs) {
		transfer.FromStateID = 10
		encodedTx, err := encoder.EncodeTransferForCommitment(&transfer)
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

	c2t := testutils.MakeCreate2Transfer(0, nil, 0, 100, s.wallets[0].PublicKey())
	s.submitInvalidBatchInTx(&c2t, func(storage *st.Storage, commitment *models.CommitmentWithTxs) {
		// Fix post state root
		_, err := storage.StateTree.Set(3, &models.UserState{
			PubKeyID: 1234,
			TokenID:  models.MakeUint256(0),
			Balance:  c2t.Amount,
			Nonce:    models.MakeUint256(0),
		})
		s.NoError(err)
		root, err := storage.StateTree.Root()
		s.NoError(err)
		commitment.PostStateRoot = *root

		// Replace toStateID and toPubKeyID in C2T
		c2t.ToStateID = ref.Uint32(3)
		encodedTx, err := encoder.EncodeCreate2TransferForCommitment(&c2t, 1234)
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

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_AllowsTransferToNonexistentReceiver() {
	transfer := testutils.MakeTransfer(0, 3, 0, 100)
	s.setTransferHashAndSign(&transfer)

	txController, txStorage, txsCtx := s.beginTransaction()
	defer txController.Rollback(nil)

	// Temporary set receiver StateLeaf for successful batch submission
	_, err := txStorage.StateTree.Set(transfer.ToStateID, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	s.submitBatch(txStorage, txsCtx, &transfer)

	s.syncAllBlocks()

	expectedReceiverState := models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  transfer.Amount,
		Nonce:    models.MakeUint256(0),
	}

	leaf, err := s.storage.StateTree.Leaf(transfer.ToStateID)
	s.NoError(err)
	s.Equal(expectedReceiverState, leaf.UserState)
}

func (s *TxsBatchesTestSuite) TestUnsafeSyncBatches_SyncsBatchesBeforeInvalidOne() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 250),
		testutils.MakeTransfer(0, 1, 1, 100),
	}
	s.setTransferHashAndSign(&transfers[0])

	txController, txStorage, txsCtx := s.beginTransaction()
	defer txController.Rollback(nil)

	s.submitBatch(txStorage, txsCtx, &transfers[0])
	invalidBatch := s.submitBatch(txStorage, txsCtx, &transfers[1])

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

// Make sure that the commander and the rollup context uses the same storage
func (s *TxsBatchesTestSuite) createTransferBatch(tx *models.Transfer) *models.Batch {
	err := s.cmd.storage.AddTransfer(tx)
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

func (s *TxsBatchesTestSuite) submitFraudulentBatchWithTransferToSelf() *models.Batch {
	storage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, storage.Teardown)

	selfTransfer := testutils.MakeTransfer(0, 0, 0, 100)
	invalidTransfer := testutils.MakeTransfer(1, 0, 0, 100)
	txs := []models.Transfer{selfTransfer, invalidTransfer}

	// Apply self transfer
	commitmentTokenID := models.MakeUint256(0)
	_, txErr, appErr := txsCtx.Applier.ApplyTransferForSync(&selfTransfer, commitmentTokenID)
	s.NoError(txErr)
	s.NoError(appErr)

	// Apply invalid transfer
	receiver, err := storage.StateTree.LeafOrEmpty(invalidTransfer.ToStateID)
	s.NoError(err)

	receiver.Balance = *receiver.Balance.Add(&invalidTransfer.Amount)
	_, err = storage.StateTree.Set(invalidTransfer.ToStateID, &receiver.UserState)
	s.NoError(err)

	// Apply fee
	feeReceiverStateID := uint32(0)
	_, err = txsCtx.Applier.ApplyFee(feeReceiverStateID, *selfTransfer.Fee.Add(&invalidTransfer.Fee))
	s.NoError(err)

	// Create commitment
	nextBatchID, err := storage.GetNextBatchID()
	s.NoError(err)

	postStateRoot, err := storage.StateTree.Root()
	s.NoError(err)

	domain, err := s.client.GetDomain()
	s.NoError(err)

	combinedSignature, err := executor.CombineSignatures(models.MakeTransferArray(txs...), domain)
	s.NoError(err)

	serializedTxs, err := encoder.SerializeTransfers(txs)
	s.NoError(err)

	commitment := models.CommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				ID: models.CommitmentID{
					BatchID:      *nextBatchID,
					IndexInBatch: 0,
				},
				Type:          batchtype.Transfer,
				PostStateRoot: *postStateRoot,
			},
			FeeReceiver:       feeReceiverStateID,
			CombinedSignature: *combinedSignature,
		},
		Transactions: serializedTxs,
	}

	// Submit batch
	batch, err := s.client.SubmitTransfersBatchAndWait(nextBatchID, []models.CommitmentWithTxs{commitment})
	s.NoError(err)

	return batch
}

func (s *TxsBatchesTestSuite) submitBatch(
	storage *st.Storage,
	txsCtx *executor.TxsContext,
	tx models.GenericTransaction,
) *models.Batch {
	pendingBatch := submitInvalidTxsBatch(s.Assertions, storage, txsCtx, tx,
		func(storage *st.Storage, commitment *models.CommitmentWithTxs) {},
	)

	s.client.GetBackend().Commit()
	return pendingBatch
}

func (s *TxsBatchesTestSuite) submitBatchInTx(tx models.GenericTransaction) *models.Batch {
	return s.submitInvalidBatchInTx(tx, func(_ *st.Storage, _ *models.CommitmentWithTxs) {})
}

// Make sure that the commander and the rollup context uses the same storage
func (s *TxsBatchesTestSuite) submitInvalidBatchInTx(
	tx models.GenericTransaction,
	modifier func(storage *st.Storage, commitment *models.CommitmentWithTxs),
) *models.Batch {
	txController, txStorage, txsCtx := s.beginTransaction()
	defer txController.Rollback(nil)

	pendingBatch := submitInvalidTxsBatch(s.Assertions, txStorage, txsCtx, tx, modifier)

	s.client.GetBackend().Commit()
	return pendingBatch
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
	if tx.Type() == txtype.Transfer {
		err := storage.AddTransfer(tx.ToTransfer())
		s.NoError(err)
	} else {
		err := storage.AddCreate2Transfer(tx.ToCreate2Transfer())
		s.NoError(err)
	}

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

func (s *TxsBatchesTestSuite) updateBatchAfterSubmission(batch *eth.DecodedTxBatch) {
	err := s.cmd.storage.UpdateBatch(batch.ToBatch(utils.RandomHash()))
	s.NoError(err)

	commitments, err := s.cmd.storage.GetTxCommitmentsByBatchID(batch.ID)
	s.NoError(err)
	for i := range commitments {
		commitments[i].BodyHash = batch.Commitments[i].BodyHash(batch.AccountTreeRoot)
	}

	err = s.cmd.storage.UpdateCommitments(commitments)
	s.NoError(err)
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

func (s *TxsBatchesTestSuite) getTransferCombinedSignature(transfer *models.Transfer) *models.Signature {
	domain, err := s.client.GetDomain()
	s.NoError(err)
	sig, err := executor.CombineSignatures(models.MakeTransferArray(*transfer), domain)
	s.NoError(err)
	return sig
}

func (s *TxsBatchesTestSuite) cloneStorage() (*st.TestStorage, *executor.TxsContext) {
	clonedStorage, err := s.storage.Clone()
	s.NoError(err)

	executionCtx := executor.NewTestExecutionContext(clonedStorage.Storage, s.client.Client, s.cfg.Rollup)
	txsCtx := executor.NewTestTxsContext(executionCtx, s.txsCtx.BatchType)

	return clonedStorage, txsCtx
}

func teardown(s *require.Assertions, teardown func() error) {
	err := teardown()
	s.NoError(err)
}

func TestTxsBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(TxsBatchesTestSuite))
}
