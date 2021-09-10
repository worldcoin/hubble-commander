package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncC2TBatchTestSuite struct {
	SyncTestSuite
}

func (s *SyncC2TBatchTestSuite) TestSyncBatch_TooManyTxsInCommitment() {
	tx := testutils.MakeCreate2Transfer(0, nil, 0, 400, s.wallets[0].PublicKey())
	s.setTxHashAndSign(&tx)
	createAndSubmitC2TBatch(s.Assertions, s.client, s.executionCtx, &tx)

	tx2 := testutils.MakeCreate2Transfer(0, nil, 1, 400, s.wallets[0].PublicKey())
	s.setTxHashAndSign(&tx2)
	s.createAndSubmitInvalidBatch(&tx2)

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

func (s *SyncC2TBatchTestSuite) TestSyncBatch_InvalidCommitmentStateRoot() {
	tx := testutils.MakeCreate2Transfer(0, nil, 0, 400, s.wallets[0].PublicKey())
	s.setTxHashAndSign(&tx)
	createAndSubmitC2TBatch(s.Assertions, s.client, s.executionCtx, &tx)

	tx2 := testutils.MakeCreate2Transfer(0, nil, 1, 400, s.wallets[0].PublicKey())
	s.setTxHashAndSign(&tx2)

	batch, commitments := createC2TBatch(s.Assertions, s.executionCtx, &tx2, testDomain)
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

func (s *SyncC2TBatchTestSuite) TestSyncBatch_InvalidTxSignature() {
	tx := testutils.MakeCreate2Transfer(0, nil, 0, 400, s.wallets[0].PublicKey())
	signCreate2Transfer(s.T(), &s.wallets[1], &tx)
	s.setTxHash(&tx)

	createAndSubmitC2TBatch(s.Assertions, s.client, s.executionCtx, &tx)

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

func (s *SyncC2TBatchTestSuite) TestSyncBatch_SingleBatch() {
	tx := testutils.MakeCreate2Transfer(0, nil, 0, 400, s.wallets[0].PublicKey())
	s.setTxHashAndSign(&tx)
	expectedCommitment := createAndSubmitC2TBatch(s.Assertions, s.client, s.executionCtx, &tx)

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

	commitment, err := s.storage.GetCommitment(&expectedCommitment.ID)
	s.NoError(err)
	s.Equal(expectedCommitment, *commitment)

	transfer, err := s.storage.GetCreate2Transfer(tx.Hash)
	s.NoError(err)
	transfer.Signature = tx.Signature
	tx.CommitmentID = &commitment.ID
	tx.ToStateID = transfer.ToStateID
	s.Equal(tx, *transfer)
}

func (s *SyncC2TBatchTestSuite) TestSyncBatch_CommitmentWithoutTxs() {
	commitment := s.createCommitmentWithEmptyTransactions(txtype.Create2Transfer)

	_, err := s.executionCtx.client.SubmitCreate2TransfersBatchAndWait([]models.Commitment{commitment})
	s.NoError(err)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.executionCtx.SyncBatch(&remoteBatches[0])
	s.NoError(err)
}

func (s *SyncC2TBatchTestSuite) createAndSubmitInvalidBatch(tx *models.Create2Transfer) models.Commitment {
	pendingBatch, commitments := createC2TBatch(s.Assertions, s.executionCtx, tx, testDomain)

	commitments[0].Transactions = append(commitments[0].Transactions, commitments[0].Transactions...)

	err := s.executionCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return commitments[0]
}

func (s *SyncC2TBatchTestSuite) setTxHash(tx *models.Create2Transfer) {
	hash, err := encoder.HashCreate2Transfer(tx)
	s.NoError(err)
	tx.Hash = *hash
}

func (s *SyncC2TBatchTestSuite) setTxHashAndSign(txs ...*models.Create2Transfer) {
	for i := range txs {
		signCreate2Transfer(s.T(), &s.wallets[txs[i].FromStateID], txs[i])
		s.setTxHash(txs[i])
	}
}

func createAndSubmitC2TBatch(
	s *require.Assertions,
	client *eth.TestClient,
	executionCtx *ExecutionContext,
	tx *models.Create2Transfer,
) models.Commitment {
	domain, err := client.GetDomain()
	s.NoError(err)
	pendingBatch, commitments := createC2TBatch(s, executionCtx, tx, domain)

	err = executionCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	client.Commit()
	return commitments[0]
}

func createC2TBatch(
	s *require.Assertions,
	executionCtx *ExecutionContext,
	tx *models.Create2Transfer,
	domain *bls.Domain,
) (*models.Batch, []models.Commitment) {
	err := executionCtx.storage.AddCreate2Transfer(tx)
	s.NoError(err)

	pendingBatch, err := executionCtx.NewPendingBatch(txtype.Create2Transfer)
	s.NoError(err)

	commitments, err := executionCtx.CreateCreate2TransferCommitments(domain)
	s.NoError(err)
	s.Len(commitments, 1)
	return pendingBatch, commitments
}

func TestSyncC2TBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SyncC2TBatchTestSuite))
}
