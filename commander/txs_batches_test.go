package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
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
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxsBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd         *Commander
	client      *eth.TestClient
	testStorage *st.TestStorage
	txsCtx      *executor.TxsContext
	cfg         *config.Config
	wallets     []bls.Wallet
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
	s.testStorage, err = st.NewTestStorage()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.cmd = NewCommander(s.cfg, nil)
	s.cmd.client = s.client.Client
	s.cmd.storage = s.testStorage.Storage
	s.cmd.metrics = metrics.NewCommanderMetrics()
	s.cmd.workersContext, s.cmd.stopWorkers = context.WithCancel(context.Background())

	executionCtx := executor.NewTestExecutionContext(s.testStorage.Storage, s.client.Client, s.cfg.Rollup)
	s.txsCtx = executor.NewTestTxsContext(executionCtx, batchtype.Transfer)

	err = s.cmd.addGenesisBatch()
	s.NoError(err)

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)
	seedDB(s.T(), s.testStorage.Storage, s.wallets)
}

func (s *TxsBatchesTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.client.Close()
	err := s.testStorage.Teardown()
	s.NoError(err)
}

func (s *TxsBatchesTestSuite) TestUnsafeSyncBatches_DoesNotSyncExistingBatchTwice() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	signTransfer(s.T(), &s.wallets[tx.FromStateID], &tx)
	clonedStorage, txsCtx := s.cloneStorage()
	s.submitBatch(clonedStorage.Storage, txsCtx, &tx)
	teardown(s.Assertions, clonedStorage.Teardown)

	s.syncAllBlocks()

	tx2 := testutils.MakeTransfer(1, 0, 0, 100)
	signTransfer(s.T(), &s.wallets[tx2.FromStateID], &tx2)
	clonedStorage, txsCtx = s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitBatch(clonedStorage.Storage, txsCtx, &tx2)

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

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_ReplaceLocalBatchWithRemoteOne() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 100),
		testutils.MakeTransfer(0, 1, 0, 200),
	}
	for i := range transfers {
		s.setTransferHashAndSign(&transfers[i])
	}

	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitBatch(clonedStorage.Storage, txsCtx, &transfers[0])

	root, err := s.cmd.storage.StateTree.Root()
	s.NoError(err)

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
	expectedCommitment := models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      batch.ID,
				IndexInBatch: 0,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: txBatch.Commitments[0].StateRoot,
		},
		FeeReceiver:       txBatch.Commitments[0].FeeReceiver,
		CombinedSignature: txBatch.Commitments[0].CombinedSignature,
	}
	expectedCommitment.BodyHash = txBatch.Commitments[0].BodyHash(*batch.AccountTreeRoot)
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
	transfer := testutils.MakeTransfer(0, 1, 0, 50)
	s.submitBatch(s.testStorage.Storage, s.txsCtx, &transfer)

	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)

	transfer = testutils.MakeTransfer(0, 1, 1, 100)
	s.submitInvalidBatch(clonedStorage.Storage, txsCtx, &transfer, func(commitment *models.CommitmentWithTxs) {
		commitment.Transactions = append(commitment.Transactions, make([]byte, 32*encoder.TransferLength)...)
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	s.updateBatchAfterSubmission(remoteBatches[0].ToDecodedTxBatch())

	err = s.cmd.syncRemoteBatch(remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].GetID())
	s.Equal(result.TooManyTx, s.getDisputeResult())
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesBatchWithInvalidPostStateRoot() {
	transfer := testutils.MakeTransfer(0, 1, 0, 50)
	s.submitBatch(s.testStorage.Storage, s.txsCtx, &transfer)

	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)

	invalidTransfer := testutils.MakeTransfer(0, 1, 1, 100)
	s.submitInvalidBatch(clonedStorage.Storage, txsCtx, &invalidTransfer, func(commitment *models.CommitmentWithTxs) {
		commitment.PostStateRoot = utils.RandomHash()
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	s.updateBatchAfterSubmission(remoteBatches[0].ToDecodedTxBatch())

	err = s.cmd.syncRemoteBatch(remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].GetID())
	s.Equal(result.Ok, s.getDisputeResult()) // invalid post state root emits Result.Ok
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithInvalidSignature() {
	s.registerAccounts([]uint32{0, 1})

	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)

	invalidTransfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatch(clonedStorage.Storage, txsCtx, &invalidTransfer, func(commitment *models.CommitmentWithTxs) {
		commitment.CombinedSignature[0] = 0 // break the signature
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[0].GetID())
	s.Equal(result.BadSignature, s.getDisputeResult())
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithSignatureInBadFormat() {
	s.registerAccounts([]uint32{0, 1})

	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)

	invalidTransfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatch(clonedStorage.Storage, txsCtx, &invalidTransfer, func(commitment *models.CommitmentWithTxs) {
		commitment.CombinedSignature = models.Signature{1, 2, 3}
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[0].GetID())
	s.Equal(result.BadPrecompileCall, s.getDisputeResult())
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_RemovesExistingBatchAndDisputesFraudulentOne() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 250),
		testutils.MakeTransfer(0, 1, 1, 100),
	}

	s.submitBatch(s.testStorage.Storage, s.txsCtx, &transfers[0])

	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitInvalidBatch(clonedStorage.Storage, txsCtx, &transfers[1], func(commitment *models.CommitmentWithTxs) {
		commitment.Transactions = append(commitment.Transactions, 1, 2, 3, 4)
	})

	localBatch := s.createTransferBatch(&transfers[2])

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	s.updateBatchAfterSubmission(remoteBatches[0].ToDecodedTxBatch())

	s.client.Account = s.client.Accounts[1]
	err = s.cmd.syncRemoteBatch(remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].GetID())
	s.Equal(result.BadCompression, s.getDisputeResult())
	_, err = s.cmd.storage.GetBatch(localBatch.ID)
	s.True(st.IsNotFoundError(err))
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesFraudulentCommitmentAfterGenesisOne() {
	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)

	invalidTransfer := testutils.MakeTransfer(0, 1, 0, 100)
	s.submitInvalidBatch(clonedStorage.Storage, txsCtx, &invalidTransfer, func(commitment *models.CommitmentWithTxs) {
		commitment.Transactions = append(commitment.Transactions, 1, 2, 3, 4)
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.cmd.syncRemoteBatch(remoteBatches[0])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[0].GetID())
	s.Equal(result.BadCompression, s.getDisputeResult())
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithInvalidFeeReceiverTokenID() {
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

	s.submitBatch(s.testStorage.Storage, s.txsCtx, &transfers[0])

	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitInvalidBatch(clonedStorage.Storage, txsCtx, &transfers[1], func(commitment *models.CommitmentWithTxs) {
		commitment.FeeReceiver = 2
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	s.updateBatchAfterSubmission(remoteBatches[0].ToDecodedTxBatch())

	err = s.cmd.syncRemoteBatch(remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].GetID())
	s.Equal(result.BadToTokenID, s.getDisputeResult())
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithoutTransfersAndInvalidPostStateRoot() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 100),
	}

	s.submitBatch(s.testStorage.Storage, s.txsCtx, &transfers[0])

	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitInvalidBatch(clonedStorage.Storage, txsCtx, &transfers[1], func(commitment *models.CommitmentWithTxs) {
		commitment.Transactions = []byte{}
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	s.updateBatchAfterSubmission(remoteBatches[0].ToDecodedTxBatch())

	err = s.cmd.syncRemoteBatch(remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].GetID())
	s.Equal(result.Ok, s.getDisputeResult()) // invalid post state root emits Result.Ok
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesCommitmentWithNonexistentSender() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 100),
	}

	s.submitBatch(s.testStorage.Storage, s.txsCtx, &transfers[0])

	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitInvalidBatch(clonedStorage.Storage, txsCtx, &transfers[1], func(commitment *models.CommitmentWithTxs) {
		transfers[1].FromStateID = 10
		encodedTx, err := encoder.EncodeTransferForCommitment(&transfers[1])
		s.NoError(err)
		commitment.Transactions = encodedTx
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	s.updateBatchAfterSubmission(remoteBatches[0].ToDecodedTxBatch())

	err = s.cmd.syncRemoteBatch(remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].GetID())
	s.Equal(result.NotEnoughTokenBalance, s.getDisputeResult())
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_DisputesC2TWithNonRegisteredReceiverPublicKey() {
	// Change batch type used in TxsContext
	executionCtx := executor.NewTestExecutionContext(s.testStorage.Storage, s.client.Client, s.cfg.Rollup)
	s.txsCtx = executor.NewTestTxsContext(executionCtx, batchtype.Create2Transfer)

	// Register public keys added to the account tree for signature disputes to work
	for i := range s.wallets {
		_, err := s.client.RegisterAccountAndWait(s.wallets[i].PublicKey())
		s.NoError(err)
	}

	// Submit a single batch to override genesis state root
	firstC2T := testutils.MakeCreate2Transfer(0, nil, 0, 100, &models.PublicKey{1, 2, 3})
	s.submitBatch(s.testStorage.Storage, s.txsCtx, &firstC2T)

	invalidC2T := testutils.MakeCreate2Transfer(0, nil, 1, 100, &models.PublicKey{1, 2, 3})
	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)
	s.submitInvalidBatch(clonedStorage.Storage, txsCtx, &invalidC2T, func(commitment *models.CommitmentWithTxs) {
		// Fix post state root
		_, err := clonedStorage.Storage.StateTree.Set(3, &models.UserState{
			PubKeyID: 1234,
			TokenID:  models.MakeUint256(0),
			Balance:  invalidC2T.Amount,
			Nonce:    models.MakeUint256(0),
		})
		s.NoError(err)
		root, err := clonedStorage.Storage.StateTree.Root()
		s.NoError(err)
		commitment.PostStateRoot = *root

		// Replace toStateID and toPubKeyID in C2T
		invalidC2T.ToStateID = ref.Uint32(3)
		encodedTx, err := encoder.EncodeCreate2TransferForCommitment(&invalidC2T, 1234)
		s.NoError(err)
		commitment.Transactions = encodedTx
	})

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	s.updateBatchAfterSubmission(remoteBatches[0].ToDecodedTxBatch())

	err = s.cmd.syncRemoteBatch(remoteBatches[1])
	s.ErrorIs(err, ErrRollbackInProgress)

	s.checkBatchAfterDispute(remoteBatches[1].GetID())
	s.Equal(result.NonexistentReceiver, s.getDisputeResult())
}

func (s *TxsBatchesTestSuite) TestSyncRemoteBatch_AllowsTransferToNonexistentReceiver() {
	transfer := testutils.MakeTransfer(0, 2, 0, 100)
	s.setTransferHashAndSign(&transfer)

	encodedTx, err := encoder.EncodeTransferForCommitment(&transfer)
	s.NoError(err)

	stateRoot := common.HexToHash("0x09de852e52fff821a7384b6bce2d5c51e9f0d32484e14c2fa29fb140d54ae8e8")

	batch := &eth.DecodedTxBatch{
		DecodedBatchBase: eth.DecodedBatchBase{
			ID:              models.MakeUint256(1),
			Type:            batchtype.Transfer,
			TransactionHash: common.Hash{1, 2, 3},
			AccountTreeRoot: common.Hash{1, 2, 3},
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

func (s *TxsBatchesTestSuite) TestUnsafeSyncBatches_SyncsBatchesBeforeInvalidOne() {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 1, 0, 50),
		testutils.MakeTransfer(0, 1, 1, 250),
		testutils.MakeTransfer(0, 1, 2, 100),
	}

	s.submitBatch(s.testStorage.Storage, s.txsCtx, &transfers[0])

	clonedStorage, txsCtx := s.cloneStorage()
	defer teardown(s.Assertions, clonedStorage.Teardown)

	invalidBatch := s.submitBatch(clonedStorage.Storage, txsCtx, &transfers[1])
	s.submitBatch(clonedStorage.Storage, txsCtx, &transfers[2])

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

	commitments, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)
	err = s.cmd.storage.AddTxCommitment(&commitments[0].TxCommitment)
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
	return s.submitInvalidBatch(storage, txsCtx, tx, func(commitment *models.CommitmentWithTxs) {})
}

// Make sure that the commander and the rollup context uses the same storage
func (s *TxsBatchesTestSuite) submitInvalidBatch(
	storage *st.Storage,
	txsCtx *executor.TxsContext,
	tx models.GenericTransaction,
	modifier func(commitment *models.CommitmentWithTxs),
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

	commitments, err := txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)

	modifier(&commitments[0])

	err = txsCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.GetBackend().Commit()
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

func (s *TxsBatchesTestSuite) checkBatchAfterDispute(batchID models.Uint256) {
	_, err := s.client.GetBatch(&batchID)
	s.Error(err)
	s.Equal(eth.MsgInvalidBatchID, err.Error())

	batch, err := s.cmd.storage.GetBatch(batchID)
	s.Nil(batch)
	s.True(st.IsNotFoundError(err))
}

func (s *TxsBatchesTestSuite) getDisputeResult() result.DisputeResult {
	it, err := s.client.Rollup.FilterRollbackTriggered(nil)
	s.NoError(err)
	it.Next()
	s.NoError(it.Error())
	res := result.DisputeResult(it.Event.Result)
	s.False(it.Next())
	return res
}

func (s *TxsBatchesTestSuite) registerAccounts(pubKeyIDs []uint32) {
	for i := range pubKeyIDs {
		leaf, err := s.testStorage.AccountTree.Leaf(pubKeyIDs[i])
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
	clonedStorage, err := s.testStorage.Clone()
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
