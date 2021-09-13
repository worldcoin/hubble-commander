package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncTransferBatchTestSuite struct {
	SyncTestSuite
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
		err := s.storage.AddTransfer(txs[i])
		s.NoError(err)
	}

	expectedCommitments, err := s.executionCtx.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(expectedCommitments, 2)
	accountRoots := make([]common.Hash, 2)
	for i := range expectedCommitments {
		var pendingBatch *models.Batch
		pendingBatch, err = s.executionCtx.NewPendingBatch(txtype.Transfer)
		s.NoError(err)
		expectedCommitments[i].ID.BatchID = pendingBatch.ID
		expectedCommitments[i].ID.IndexInBatch = 0
		err = s.executionCtx.SubmitBatch(pendingBatch, []models.Commitment{expectedCommitments[i]})
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
		commitment, err := s.storage.GetCommitment(&expectedCommitments[i].ID)
		s.NoError(err)
		s.Equal(expectedCommitments[i], *commitment)

		actualTx, err := s.storage.GetTransfer(txs[i].Hash)
		s.NoError(err)
		txs[i].CommitmentID = &commitment.ID
		txs[i].Signature = models.Signature{}
		s.Equal(txs[i], actualTx)
	}
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_PendingBatch() {
	accountRoot := s.getAccountTreeRoot()
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTxHashAndSign(&tx)
	submitTransferBatch(s.Assertions, s.client, s.executionCtx, &tx)

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

func (s *SyncTransferBatchTestSuite) TestSyncBatch_TooManyTxsInCommitment() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTxHashAndSign(&tx)
	submitTransferBatch(s.Assertions, s.client, s.executionCtx, &tx)

	tx2 := testutils.MakeTransfer(0, 1, 1, 400)
	s.setTxHashAndSign(&tx2)
	s.submitInvalidBatch(&tx2)

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.executionCtx.SyncBatch(&remoteBatches[0])
	s.NoError(err)

	var disputableErr *DisputableError
	err = s.executionCtx.SyncBatch(&remoteBatches[1])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(ErrTooManyTxs.Reason, disputableErr.Reason)

	_, err = s.storage.GetBatch(remoteBatches[0].ID)
	s.NoError(err)
	_, err = s.storage.GetBatch(remoteBatches[1].ID)
	s.NoError(err)
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_InvalidCommitmentStateRoot() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTxHashAndSign(&tx)
	submitTransferBatch(s.Assertions, s.client, s.executionCtx, &tx)

	tx2 := testutils.MakeTransfer(0, 1, 1, 400)
	s.setTxHashAndSign(&tx2)

	batch, commitments := createTransferBatch(s.Assertions, s.executionCtx, &tx2, testDomain)
	commitments[0].PostStateRoot = utils.RandomHash()

	err := s.executionCtx.SubmitBatch(batch, commitments)
	s.NoError(err)
	s.client.Commit()

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.executionCtx.SyncBatch(&remoteBatches[0])
	s.NoError(err)

	var disputableErr *DisputableError
	err = s.executionCtx.SyncBatch(&remoteBatches[1])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(ErrInvalidCommitmentStateRoot.Error(), disputableErr.Reason)

	_, err = s.storage.GetBatch(remoteBatches[0].ID)
	s.NoError(err)
	_, err = s.storage.GetBatch(remoteBatches[1].ID)
	s.NoError(err)
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_InvalidTxSignature() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	signTransfer(s.T(), &s.wallets[1], &tx)
	s.setTxHash(&tx)

	submitTransferBatch(s.Assertions, s.client, s.executionCtx, &tx)

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	var disputableErr *DisputableError
	err = s.executionCtx.SyncBatch(&remoteBatches[0])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Signature, disputableErr.Type)
	s.Equal(InvalidSignature, disputableErr.Reason)
	s.Equal(0, disputableErr.CommitmentIndex)
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_NotValidBLSSignature() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	s.setTxHash(&tx)

	pendingBatch, commitments := createTransferBatch(s.Assertions, s.executionCtx, &tx, s.domain)
	commitments[0].CombinedSignature = models.Signature{1, 2, 3}

	err := s.executionCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)
	s.client.Commit()

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	var disputableErr *DisputableError
	err = s.executionCtx.SyncBatch(&remoteBatches[0])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Signature, disputableErr.Type)
	s.Equal(0, disputableErr.CommitmentIndex)
}

func (s *SyncTransferBatchTestSuite) TestSyncBatch_CommitmentWithoutTxs() {
	commitment := s.createCommitmentWithEmptyTransactions(txtype.Transfer)

	_, err := s.executionCtx.client.SubmitTransfersBatchAndWait([]models.Commitment{commitment})
	s.NoError(err)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.executionCtx.SyncBatch(&remoteBatches[0])
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

	feeReceiver, err := s.executionCtx.storage.StateTree.Leaf(feeReceiverStateID)
	s.NoError(err)
	sender, err := s.executionCtx.storage.StateTree.Leaf(0)
	s.NoError(err)
	receiver, err := s.executionCtx.storage.StateTree.Leaf(1)
	s.NoError(err)

	s.Equal(expectedNewlyCreatedFeeReceiver, feeReceiver)
	s.Equal(models.MakeUint256(1000-400-100), sender.Balance)
	s.Equal(models.MakeUint256(400), receiver.Balance)
}

func (s *SyncTransferBatchTestSuite) submitInvalidBatch(tx *models.Transfer) *models.Batch {
	pendingBatch, commitments := createTransferBatch(s.Assertions, s.executionCtx, tx, testDomain)

	commitments[0].Transactions = append(commitments[0].Transactions, commitments[0].Transactions...)

	err := s.executionCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *SyncTransferBatchTestSuite) submitTransferBatchWithNonexistentFeeReceiver(tx *models.Transfer, feeReceiverStateID uint32) {
	commitmentTokenID := models.MakeUint256(0)

	receiverLeaf, err := s.executionCtx.storage.StateTree.Leaf(tx.ToStateID)
	s.NoError(err)
	txErr, appErr := s.executionCtx.ApplyTransfer(tx, receiverLeaf, commitmentTokenID)
	s.NoError(txErr)
	s.NoError(appErr)

	_, commitmentErr, appErr := s.executionCtx.ApplyFeeForSync(feeReceiverStateID, &commitmentTokenID, &tx.Fee)
	s.NoError(commitmentErr)
	s.NoError(appErr)

	serializedTxs, err := encoder.SerializeTransfers([]models.Transfer{*tx})
	s.NoError(err)

	combinedSignature, err := CombineSignatures(models.MakeTransferArray(*tx), s.domain)
	s.NoError(err)

	postStateRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	nextBatchID, err := s.executionCtx.storage.GetNextBatchID()
	s.NoError(err)

	commitment := models.Commitment{
		ID: models.CommitmentID{
			BatchID:      *nextBatchID,
			IndexInBatch: 0,
		},
		Type:              txtype.Transfer,
		Transactions:      serializedTxs,
		FeeReceiver:       feeReceiverStateID,
		CombinedSignature: *combinedSignature,
		PostStateRoot:     *postStateRoot,
	}
	_, err = s.client.SubmitTransfersBatchAndWait([]models.Commitment{commitment})
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

func submitTransferBatch(
	s *require.Assertions,
	client *eth.TestClient,
	executionCtx *ExecutionContext,
	tx *models.Transfer,
) *models.Batch {
	domain, err := client.GetDomain()
	s.NoError(err)
	pendingBatch, commitments := createTransferBatch(s, executionCtx, tx, domain)

	err = executionCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	client.Commit()
	return pendingBatch
}

func createTransferBatch(
	s *require.Assertions,
	executionCtx *ExecutionContext,
	tx *models.Transfer,
	domain *bls.Domain,
) (*models.Batch, []models.Commitment) {
	err := executionCtx.storage.AddTransfer(tx)
	s.NoError(err)

	pendingBatch, err := executionCtx.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitments, err := executionCtx.CreateTransferCommitments(domain)
	s.NoError(err)
	s.Len(commitments, 1)

	return pendingBatch, commitments
}

func TestSyncTransferBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SyncTransferBatchTestSuite))
}