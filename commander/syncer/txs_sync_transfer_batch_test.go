package syncer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type SyncTransferBatchTestSuite struct {
	syncTestSuite
}

func (s *SyncTransferBatchTestSuite) SetupTest() {
	s.testSuiteWithSyncAndRollupContext.SetupTestWithConfig(batchtype.Transfer, &syncTestSuiteConfig)
	s.syncTestSuite.setupTest()
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_TwoBatches() {
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
	s.setTxHashAndSign(txs...)
	for i := range txs {
		err := s.storage.AddTransaction(txs[i])
		s.NoError(err)
	}

	commitments, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 2)
	accountRoots := make([]common.Hash, len(commitments))
	expectedCommitments := make([]models.TxCommitment, 0, len(commitments))
	for i := range commitments {
		var pendingBatch *models.Batch
		pendingBatch, err = s.txsCtx.NewPendingBatch(batchtype.Transfer)
		s.NoError(err)
		commitments[i].ToTxCommitmentWithTxs().ID.BatchID = pendingBatch.ID
		commitments[i].ToTxCommitmentWithTxs().ID.IndexInBatch = 0

		err = s.txsCtx.SubmitBatch(pendingBatch, []models.CommitmentWithTxs{commitments[i]})
		s.NoError(err)
		s.client.GetBackend().Commit()

		accountRoots[i] = s.getAccountTreeRoot()
		commitments[i].CalcAndSetBodyHash(accountRoots[i])
		expectedCommitments = append(expectedCommitments, commitments[i].ToTxCommitmentWithTxs().TxCommitment)
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
		commitment, err := s.storage.GetCommitment(&commitments[i].ToTxCommitmentWithTxs().ID)
		s.NoError(err)
		s.Equal(expectedCommitments[i], *commitment.ToTxCommitment())

		actualTx, err := s.storage.GetTransfer(txs[i].Hash)
		s.NoError(err)
		txs[i].CommitmentID = &commitment.ToTxCommitment().ID
		txs[i].Signature = models.Signature{}
		s.Equal(txs[i], actualTx)
	}
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_SyncsExistingBatch() {
	accountRoot := s.getAccountTreeRoot()
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTxHashAndSign(&tx)
	s.submitBatch(&tx)

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

	commitments, err := s.storage.GetCommitmentsByBatchID(batches[0].ID)
	s.NoError(err)
	s.Len(commitments, 1)
	s.NotNil(commitments[0].ToTxCommitment().BodyHash)
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_TooManyTxsInCommitment() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTxHashAndSign(&tx)
	s.submitBatch(&tx)

	tx2 := testutils.MakeTransfer(0, 1, 1, 400)
	s.setTxHashAndSign(&tx2)
	s.submitInvalidBatch(&tx2)

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.syncCtx.SyncBatch(remoteBatches[0])
	s.NoError(err)

	var disputableErr *DisputableError
	err = s.syncCtx.SyncBatch(remoteBatches[1])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(ErrTooManyTxs.Reason, disputableErr.Reason)

	_, err = s.storage.GetBatch(remoteBatches[0].GetID())
	s.NoError(err)
	_, err = s.storage.GetBatch(remoteBatches[1].GetID())
	s.NoError(err)
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_InvalidCommitmentStateRoot() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTxHashAndSign(&tx)
	s.submitBatch(&tx)

	tx2 := testutils.MakeTransfer(0, 1, 1, 400)
	s.setTxHashAndSign(&tx2)

	batch, commitments := s.createBatch(&tx2)
	commitments[0].ToTxCommitmentWithTxs().PostStateRoot = utils.RandomHash()

	err := s.txsCtx.SubmitBatch(batch, commitments)
	s.NoError(err)
	s.client.GetBackend().Commit()

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.syncCtx.SyncBatch(remoteBatches[0])
	s.NoError(err)

	var disputableErr *DisputableError
	err = s.syncCtx.SyncBatch(remoteBatches[1])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(invalidStateRootMessage, disputableErr.Reason)

	_, err = s.storage.GetBatch(remoteBatches[0].GetID())
	s.NoError(err)
	_, err = s.storage.GetBatch(remoteBatches[1].GetID())
	s.NoError(err)
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_InvalidTxSignature() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	signTransfer(s.T(), &s.wallets[1], &tx)
	s.setTxHash(&tx)

	s.submitBatch(&tx)

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	var disputableErr *DisputableError
	err = s.syncCtx.SyncBatch(remoteBatches[0])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Signature, disputableErr.Type)
	s.Equal(InvalidSignatureMessage, disputableErr.Reason)
	s.Equal(0, disputableErr.CommitmentIndex)
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_NotValidBLSSignature() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTxHash(&tx)

	pendingBatch, commitments := s.createBatch(&tx)
	commitments[0].ToTxCommitmentWithTxs().CombinedSignature = models.Signature{1, 2, 3}

	err := s.txsCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)
	s.client.GetBackend().Commit()

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	var disputableErr *DisputableError
	err = s.syncCtx.SyncBatch(remoteBatches[0])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Signature, disputableErr.Type)
	s.Equal(0, disputableErr.CommitmentIndex)
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_CommitmentWithoutTxs() {
	commitment := s.createCommitmentWithEmptyTransactions(batchtype.Transfer)

	_, err := s.client.SubmitTransfersBatchAndWait(models.NewUint256(1), []models.CommitmentWithTxs{&commitment})
	s.NoError(err)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.syncCtx.SyncBatch(remoteBatches[0])
	s.NoError(err)
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_CommitmentWithNonexistentFeeReceiver() {
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
	s.setTxHashAndSign(&tx)
	s.submitTransferBatchWithNonexistentFeeReceiver(&tx, feeReceiverStateID)

	s.recreateDatabase()
	s.syncAllBatches()

	expectedNewlyCreatedFeeReceiver, err := st.NewStateLeaf(feeReceiverStateID, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(100),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	feeReceiver, err := s.storage.StateTree.Leaf(feeReceiverStateID)
	s.NoError(err)
	sender, err := s.storage.StateTree.Leaf(0)
	s.NoError(err)
	receiver, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)

	s.Equal(expectedNewlyCreatedFeeReceiver, feeReceiver)
	s.Equal(models.MakeUint256(1000-400-100), sender.Balance)
	s.Equal(models.MakeUint256(400), receiver.Balance)
}

func (s *SyncTransferBatchTestSuite) submitInvalidBatch(tx *models.Transfer) *models.Batch {
	pendingBatch, commitments := s.createBatch(tx)

	commitment := commitments[0].ToTxCommitmentWithTxs()
	commitments[0].ToTxCommitmentWithTxs().Transactions = utils.RandomBytes(uint64(2 * len(commitment.Transactions)))

	err := s.txsCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.GetBackend().Commit()
	return pendingBatch
}

func (s *SyncTransferBatchTestSuite) submitTransferBatchWithNonexistentFeeReceiver(tx *models.Transfer, feeReceiverStateID uint32) {
	commitmentTokenID := models.MakeUint256(0)

	receiverLeaf, err := s.storage.StateTree.Leaf(tx.ToStateID)
	s.NoError(err)
	txErr, appErr := s.txsCtx.Applier.ApplyTx(tx, receiverLeaf, commitmentTokenID)
	s.NoError(txErr)
	s.NoError(appErr)

	_, commitmentErr, appErr := s.syncCtx.batchCtx.(*TxsContext).Syncer.ApplyFee(feeReceiverStateID, &commitmentTokenID, &tx.Fee)
	s.NoError(commitmentErr)
	s.NoError(appErr)

	serializedTxs, err := encoder.SerializeTransfers([]models.Transfer{*tx})
	s.NoError(err)

	combinedSignature, err := executor.CombineSignatures(models.MakeTransferArray(*tx), s.domain)
	s.NoError(err)

	postStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	nextBatchID, err := s.storage.GetNextBatchID()
	s.NoError(err)

	commitment := models.TxCommitmentWithTxs{
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
	_, err = s.client.SubmitTransfersBatchAndWait(models.NewUint256(1), []models.CommitmentWithTxs{&commitment})
	s.NoError(err)
}

func (s *SyncTransferBatchTestSuite) setTxHash(tx *models.Transfer) {
	hash, err := encoder.HashTransfer(tx)
	s.NoError(err)
	tx.Hash = *hash
}

func (s *SyncTransferBatchTestSuite) setTxHashAndSign(txs ...*models.Transfer) {
	for i := range txs {
		signTransfer(s.T(), &s.wallets[txs[i].FromStateID], txs[i])
		s.setTxHash(txs[i])
	}
}

func TestSyncTransferBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SyncTransferBatchTestSuite))
}
