package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *st.TestStorage
	client              *eth.TestClient
	cfg                 *config.RollupConfig
	transactionExecutor *TransactionExecutor
	transfer            models.Transfer
	wallets             []bls.Wallet
	domain              *bls.Domain
}

func (s *SyncTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.transfer = testutils.MakeTransfer(0, 1, 0, 400)
	s.setTransferHash(&s.transfer)
}

func (s *SyncTestSuite) SetupTest() {
	var err error
	s.client, err = eth.NewConfiguredTestClient(rollup.DeploymentConfig{
		Params: rollup.Params{
			MaxTxsPerCommit: models.NewUint256(1),
		},
	}, eth.ClientConfig{})
	s.NoError(err)

	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
		DisableSignatures:      false,
	}

	s.domain, err = s.client.GetDomain()
	s.NoError(err)
	s.wallets = generateWallets(s.Assertions, s.domain, 2)
	s.setupDB()
}

func (s *SyncTestSuite) setupDB() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, s.client.Client, s.cfg, context.Background())

	seedDB(s.Assertions, s.storage.Storage, s.wallets)
}

func seedDB(s *require.Assertions, storage *st.Storage, wallets []bls.Wallet) {
	err := storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: *wallets[0].PublicKey(),
	})
	s.NoError(err)

	err = storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: *wallets[1].PublicKey(),
	})
	s.NoError(err)

	_, err = storage.StateTree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, err = storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *SyncTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SyncTestSuite) TestSyncBatch_TwoTransferBatches() {
	txs := []*models.Transfer{
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
	s.setTransferHashAndSign(txs...)
	for i := range txs {
		_, err := s.storage.AddTransfer(txs[i])
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
		s.Equal(txs[i], actualTx)
	}
}

func (s *SyncTestSuite) TestSyncBatch_PendingBatch() {
	accountRoot := s.getAccountTreeRoot()
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTransferHashAndSign(&tx)
	createAndSubmitTransferBatch(s.Assertions, s.client, s.transactionExecutor, &tx)

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

func (s *SyncTestSuite) TestSyncBatch_TooManyTransfersInCommitment() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTransferHashAndSign(&tx)
	createAndSubmitTransferBatch(s.Assertions, s.client, s.transactionExecutor, &tx)

	tx2 := testutils.MakeTransfer(0, 1, 1, 400)
	s.setTransferHashAndSign(&tx2)
	s.createAndSubmitInvalidTransferBatch(&tx2)

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.transactionExecutor.SyncBatch(&remoteBatches[0])
	s.NoError(err)

	var disputableErr *DisputableError
	err = s.transactionExecutor.SyncBatch(&remoteBatches[1])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(ErrTooManyTxs.Reason, disputableErr.Reason)

	_, err = s.storage.GetBatch(remoteBatches[0].ID)
	s.NoError(err)
	_, err = s.storage.GetBatch(remoteBatches[1].ID)
	s.NoError(err)
}

func (s *SyncTestSuite) TestSyncBatch_TooManyCreate2TransfersInCommitment() {
	tx := testutils.MakeCreate2Transfer(0, ref.Uint32(5), 0, 400, s.wallets[0].PublicKey())
	s.setC2THashAndSign(&tx)
	createAndSubmitC2TBatch(s.Assertions, s.client, s.transactionExecutor, &tx)

	tx2 := testutils.MakeCreate2Transfer(0, ref.Uint32(6), 1, 400, s.wallets[0].PublicKey())
	s.setC2THashAndSign(&tx2)
	s.createAndSubmitInvalidC2TBatch(&tx2)

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.transactionExecutor.SyncBatch(&remoteBatches[0])
	s.NoError(err)

	var disputableErr *DisputableError
	err = s.transactionExecutor.SyncBatch(&remoteBatches[1])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(ErrTooManyTxs.Reason, disputableErr.Reason)

	_, err = s.storage.GetBatch(remoteBatches[0].ID)
	s.NoError(err)
	_, err = s.storage.GetBatch(remoteBatches[1].ID)
	s.NoError(err)
}

func (s *SyncTestSuite) TestSyncBatch_InvalidTransferCommitmentStateRoot() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTransferHashAndSign(&tx)
	createAndSubmitTransferBatch(s.Assertions, s.client, s.transactionExecutor, &tx)

	tx2 := testutils.MakeTransfer(0, 1, 1, 400)
	s.setTransferHashAndSign(&tx2)

	batch, commitments := createTransferBatch(s.Assertions, s.transactionExecutor, &tx2, testDomain)
	commitments[0].PostStateRoot = utils.RandomHash()

	err := s.transactionExecutor.SubmitBatch(batch, commitments)
	s.NoError(err)
	s.client.Commit()

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.transactionExecutor.SyncBatch(&remoteBatches[0])
	s.NoError(err)

	var disputableErr *DisputableError
	err = s.transactionExecutor.SyncBatch(&remoteBatches[1])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(ErrInvalidCommitmentStateRoot.Error(), disputableErr.Reason)

	_, err = s.storage.GetBatch(remoteBatches[0].ID)
	s.NoError(err)
	_, err = s.storage.GetBatch(remoteBatches[1].ID)
	s.NoError(err)
}

func (s *SyncTestSuite) TestSyncBatch_InvalidCreate2TransferCommitmentStateRoot() {
	tx := testutils.MakeCreate2Transfer(0, ref.Uint32(5), 0, 400, s.wallets[0].PublicKey())
	s.setC2THashAndSign(&tx)
	createAndSubmitC2TBatch(s.Assertions, s.client, s.transactionExecutor, &tx)

	tx2 := testutils.MakeCreate2Transfer(0, ref.Uint32(6), 1, 400, s.wallets[0].PublicKey())
	s.setC2THashAndSign(&tx2)

	batch, commitments := createC2TBatch(s.Assertions, s.transactionExecutor, &tx2, testDomain)
	commitments[0].PostStateRoot = utils.RandomHash()

	err := s.transactionExecutor.SubmitBatch(batch, commitments)
	s.NoError(err)
	s.client.Commit()

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.transactionExecutor.SyncBatch(&remoteBatches[0])
	s.NoError(err)

	var disputableErr *DisputableError
	err = s.transactionExecutor.SyncBatch(&remoteBatches[1])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(ErrInvalidCommitmentStateRoot.Error(), disputableErr.Reason)

	_, err = s.storage.GetBatch(remoteBatches[0].ID)
	s.NoError(err)
	_, err = s.storage.GetBatch(remoteBatches[1].ID)
	s.NoError(err)
}

func (s *SyncTestSuite) TestSyncBatch_InvalidTransferSignature() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	signTransfer(s.T(), &s.wallets[1], &tx)
	s.setTransferHash(&tx)

	createAndSubmitTransferBatch(s.Assertions, s.client, s.transactionExecutor, &tx)

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	var disputableErr *DisputableError
	err = s.transactionExecutor.SyncBatch(&remoteBatches[0])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Signature, disputableErr.Type)
	s.Equal(InvalidSignature, disputableErr.Reason)
	s.Equal(0, disputableErr.CommitmentIndex)
}

func (s *SyncTestSuite) TestSyncBatch_InvalidCreate2TransferSignature() {
	tx := testutils.MakeCreate2Transfer(0, ref.Uint32(5), 0, 400, s.wallets[0].PublicKey())
	signCreate2Transfer(s.T(), &s.wallets[1], &tx)
	s.setCreate2TransferHash(&tx)

	createAndSubmitC2TBatch(s.Assertions, s.client, s.transactionExecutor, &tx)

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	var disputableErr *DisputableError
	err = s.transactionExecutor.SyncBatch(&remoteBatches[0])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Signature, disputableErr.Type)
	s.Equal(InvalidSignature, disputableErr.Reason)
	s.Equal(0, disputableErr.CommitmentIndex)
}

func (s *SyncTestSuite) TestSyncBatch_NotValidBLSSignature() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTransferHash(&tx)

	pendingBatch, commitments := createTransferBatch(s.Assertions, s.transactionExecutor, &tx, s.domain)
	commitments[0].CombinedSignature = models.Signature{1, 2, 3}

	err := s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)
	s.client.Commit()

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	var disputableErr *DisputableError
	err = s.transactionExecutor.SyncBatch(&remoteBatches[0])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Signature, disputableErr.Type)
	s.Equal(0, disputableErr.CommitmentIndex)
}

func (s *SyncTestSuite) TestSyncBatch_Create2TransferBatch() {
	tx := testutils.MakeCreate2Transfer(0, nil, 0, 400, s.wallets[0].PublicKey())
	s.setC2THashAndSign(&tx)
	expectedCommitment := createAndSubmitC2TBatch(s.Assertions, s.client, s.transactionExecutor, &tx)

	s.recreateDatabase()
	s.syncAllBatches()

	state0, err := s.storage.StateTree.Leaf(0)
	s.NoError(err)
	s.Equal(models.MakeUint256(600), state0.Balance)

	state2, err := s.storage.StateTree.Leaf(2)
	s.NoError(err)
	s.Equal(models.MakeUint256(400), state2.Balance)
	s.Equal(uint32(0), state2.PubKeyID)

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
	tx.ToStateID = transfer.ToStateID
	s.Equal(tx, *transfer)
}

func (s *SyncTestSuite) TestSyncBatch_CommitmentWithoutTransfers() {
	commitment := s.createCommitmentWithEmptyTransactions(txtype.Transfer)

	_, err := s.transactionExecutor.client.SubmitTransfersBatchAndWait([]models.Commitment{commitment})
	s.NoError(err)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.transactionExecutor.SyncBatch(&remoteBatches[0])
	s.NoError(err)
}

func (s *SyncTestSuite) TestSyncBatch_CommitmentWithoutCreate2Transfers() {
	commitment := s.createCommitmentWithEmptyTransactions(txtype.Create2Transfer)

	_, err := s.transactionExecutor.client.SubmitCreate2TransfersBatchAndWait([]models.Commitment{commitment})
	s.NoError(err)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.transactionExecutor.SyncBatch(&remoteBatches[0])
	s.NoError(err)
}

func (s *SyncTestSuite) TestSyncBatch_CommitmentWithNonexistentFeeReceiver() {
	feeReceiverStateID := uint32(1234)
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 1,
	}
	s.setTransferHashAndSign(&tx)
	s.createAndSubmitTransferBatchWithNonexistentFeeReceiver(&tx, feeReceiverStateID)

	s.recreateDatabase()
	s.syncAllBatches()

	expectedNewlyCreatedFeeReceiver, err := st.NewStateLeaf(feeReceiverStateID, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(100),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	feeReceiver, err := s.transactionExecutor.storage.StateTree.Leaf(feeReceiverStateID)
	s.NoError(err)
	sender, err := s.transactionExecutor.storage.StateTree.Leaf(0)
	s.NoError(err)
	receiver, err := s.transactionExecutor.storage.StateTree.Leaf(1)
	s.NoError(err)

	s.Equal(expectedNewlyCreatedFeeReceiver, feeReceiver)
	s.Equal(models.MakeUint256(1000-400-100), sender.Balance)
	s.Equal(models.MakeUint256(400), receiver.Balance)
}

func (s *SyncTestSuite) TestRevertBatch_RevertsState() {
	initialStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	signTransfer(s.T(), &s.wallets[s.transfer.FromStateID], &s.transfer)
	pendingBatch := createAndSubmitTransferBatch(s.Assertions, s.client, s.transactionExecutor, &s.transfer)

	err = s.transactionExecutor.RevertBatches(pendingBatch)
	s.NoError(err)

	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)
	s.Equal(*initialStateRoot, *stateRoot)

	state0, err := s.storage.StateTree.Leaf(s.transfer.FromStateID)
	s.NoError(err)
	s.Equal(uint64(1000), state0.Balance.Uint64())

	state1, err := s.storage.StateTree.Leaf(s.transfer.ToStateID)
	s.NoError(err)
	s.Equal(uint64(0), state1.Balance.Uint64())
}

func (s *SyncTestSuite) TestRevertBatch_ExcludesTransactionsFromCommitments() {
	signTransfer(s.T(), &s.wallets[s.transfer.FromStateID], &s.transfer)
	pendingBatch := createAndSubmitTransferBatch(s.Assertions, s.client, s.transactionExecutor, &s.transfer)

	err := s.transactionExecutor.RevertBatches(pendingBatch)
	s.NoError(err)

	transfer, err := s.storage.GetTransfer(s.transfer.Hash)
	s.NoError(err)
	s.Nil(transfer.IncludedInCommitment)
}

func (s *SyncTestSuite) TestRevertBatch_DeletesCommitmentsAndBatches() {
	transfers := make([]models.Transfer, 2)
	transfers[0] = s.transfer
	transfers[1] = testutils.MakeTransfer(0, 1, 1, 200)

	pendingBatches := make([]models.Batch, 2)
	for i := range pendingBatches {
		signTransfer(s.T(), &s.wallets[transfers[i].FromStateID], &transfers[i])
		pendingBatches[i] = *createAndSubmitTransferBatch(s.Assertions, s.client, s.transactionExecutor, &transfers[i])
	}

	latestCommitment, err := s.transactionExecutor.storage.GetLatestCommitment()
	s.NoError(err)
	s.EqualValues(2, latestCommitment.ID)

	err = s.transactionExecutor.RevertBatches(&pendingBatches[0])
	s.NoError(err)

	_, err = s.transactionExecutor.storage.GetLatestCommitment()
	s.Equal(st.NewNotFoundError("commitment"), err)

	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 0)
}

func createAndSubmitTransferBatch(
	s *require.Assertions,
	client *eth.TestClient,
	txExecutor *TransactionExecutor,
	tx *models.Transfer,
) *models.Batch {
	domain, err := client.GetDomain()
	s.NoError(err)
	pendingBatch, commitments := createTransferBatch(s, txExecutor, tx, domain)

	err = txExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	client.Commit()
	return pendingBatch
}

func (s *SyncTestSuite) createAndSubmitInvalidTransferBatch(tx *models.Transfer) *models.Batch {
	pendingBatch, commitments := createTransferBatch(s.Assertions, s.transactionExecutor, tx, testDomain)

	commitments[0].Transactions = append(commitments[0].Transactions, commitments[0].Transactions...)

	err := s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *SyncTestSuite) createAndSubmitTransferBatchWithNonexistentFeeReceiver(tx *models.Transfer, feeReceiverStateID uint32) {
	commitmentTokenID := models.MakeUint256(0)

	receiverLeaf, err := s.transactionExecutor.storage.StateTree.Leaf(tx.ToStateID)
	s.NoError(err)
	txErr, appErr := s.transactionExecutor.ApplyTransfer(tx, receiverLeaf, commitmentTokenID)
	s.NoError(txErr)
	s.NoError(appErr)

	_, commitmentErr, appErr := s.transactionExecutor.ApplyFeeForSync(feeReceiverStateID, &commitmentTokenID, &tx.Fee)
	s.NoError(commitmentErr)
	s.NoError(appErr)

	serializedTxs, err := encoder.SerializeTransfers([]models.Transfer{*tx})
	s.NoError(err)

	combinedSignature, err := combineTransferSignatures([]models.Transfer{*tx}, s.domain)
	s.NoError(err)

	postStateRoot, err := s.transactionExecutor.storage.StateTree.Root()
	s.NoError(err)

	commitment := models.Commitment{
		ID:                0,
		Type:              txtype.Transfer,
		Transactions:      serializedTxs,
		FeeReceiver:       feeReceiverStateID,
		CombinedSignature: *combinedSignature,
		PostStateRoot:     *postStateRoot,
	}
	_, err = s.client.SubmitTransfersBatchAndWait([]models.Commitment{commitment})
	s.NoError(err)
}

func (s *SyncTestSuite) createCommitmentWithEmptyTransactions(commitmentType txtype.TransactionType) models.Commitment {
	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	feeReceiver, err := s.transactionExecutor.getCommitmentFeeReceiver()
	s.NoError(err)

	return models.Commitment{
		Type:              commitmentType,
		Transactions:      []byte{},
		FeeReceiver:       feeReceiver.StateID,
		CombinedSignature: models.Signature{},
		PostStateRoot:     *stateRoot,
	}
}

func createTransferBatch(
	s *require.Assertions,
	txExecutor *TransactionExecutor,
	tx *models.Transfer,
	domain *bls.Domain,
) (*models.Batch, []models.Commitment) {
	_, err := txExecutor.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := txExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitments, err := txExecutor.CreateTransferCommitments(domain)
	s.NoError(err)
	s.Len(commitments, 1)

	return pendingBatch, commitments
}

func createAndSubmitC2TBatch(
	s *require.Assertions,
	client *eth.TestClient,
	txExecutor *TransactionExecutor,
	tx *models.Create2Transfer,
) models.Commitment {
	domain, err := client.GetDomain()
	s.NoError(err)
	pendingBatch, commitments := createC2TBatch(s, txExecutor, tx, domain)

	err = txExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	client.Commit()
	return commitments[0]
}

func (s *SyncTestSuite) createAndSubmitInvalidC2TBatch(tx *models.Create2Transfer) models.Commitment {
	pendingBatch, commitments := createC2TBatch(s.Assertions, s.transactionExecutor, tx, testDomain)

	commitments[0].Transactions = append(commitments[0].Transactions, commitments[0].Transactions...)

	err := s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return commitments[0]
}

func createC2TBatch(
	s *require.Assertions,
	txExecutor *TransactionExecutor,
	tx *models.Create2Transfer,
	domain *bls.Domain,
) (*models.Batch, []models.Commitment) {
	_, err := txExecutor.storage.AddCreate2Transfer(tx)
	s.NoError(err)

	pendingBatch, err := txExecutor.NewPendingBatch(txtype.Create2Transfer)
	s.NoError(err)

	commitments, err := txExecutor.CreateCreate2TransferCommitments(domain)
	s.NoError(err)
	s.Len(commitments, 1)
	return pendingBatch, commitments
}

func (s *SyncTestSuite) syncAllBatches() {
	newRemoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)

	for i := range newRemoteBatches {
		remoteBatch := &newRemoteBatches[i]
		err = s.transactionExecutor.SyncBatch(remoteBatch)
		s.NoError(err)
	}
}

func (s *SyncTestSuite) recreateDatabase() {
	err := s.storage.Teardown()
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

func (s *SyncTestSuite) setTransferHashAndSign(txs ...*models.Transfer) {
	for i := range txs {
		signTransfer(s.T(), &s.wallets[txs[i].FromStateID], txs[i])
		s.setTransferHash(txs[i])
	}
}

func (s *SyncTestSuite) setC2THashAndSign(txs ...*models.Create2Transfer) {
	for i := range txs {
		signCreate2Transfer(s.T(), &s.wallets[txs[i].FromStateID], txs[i])
		s.setCreate2TransferHash(txs[i])
	}
}

func generateWallets(s *require.Assertions, domain *bls.Domain, walletsAmount int) []bls.Wallet {
	wallets := make([]bls.Wallet, 0, walletsAmount)
	for i := 0; i < walletsAmount; i++ {
		wallet, err := bls.NewRandomWallet(*domain)
		s.NoError(err)
		wallets = append(wallets, *wallet)
	}
	return wallets
}

func TestSyncTestSuite(t *testing.T) {
	suite.Run(t, new(SyncTestSuite))
}
