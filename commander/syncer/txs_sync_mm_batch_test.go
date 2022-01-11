package syncer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type SyncMMBatchTestSuite struct {
	syncTestSuite
}

func (s *SyncMMBatchTestSuite) SetupTest() {
	s.testSuiteWithSyncAndRollupContext.SetupTestWithConfig(batchtype.MassMigration, &syncTestSuiteConfig)
	s.syncTestSuite.setupTest()
}

func (s *SyncMMBatchTestSuite) TestSyncBatch_SingleBatch() {
	tx := testutils.MakeMassMigration(0, 1, 0, 400)
	s.setTxHashAndSign(&tx)

	commitments := s.submitBatch(&tx)
	expectedCommitment := commitments[0].ToMMCommitmentWithTxs().MMCommitment

	s.recreateDatabase()
	s.syncAllBatches()

	senderState, err := s.storage.StateTree.Leaf(tx.FromStateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(600), senderState.Balance)

	treeRoot := s.getAccountTreeRoot()
	batches, err := s.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
	s.Equal(models.MakeUint256(1), batches[0].ID)
	s.Equal(treeRoot, *batches[0].AccountTreeRoot)

	decodedBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(decodedBatches, 1)

	commitment, err := s.storage.GetMMCommitment(&expectedCommitment.ID)
	s.NoError(err)
	expectedCommitment.BodyHash = decodedBatches[0].ToDecodedTxBatch().Commitments[0].BodyHash(s.getAccountTreeRoot())
	s.Equal(expectedCommitment, *commitment)

	massMigration, err := s.storage.GetMassMigration(tx.Hash)
	s.NoError(err)
	massMigration.Signature = tx.Signature
	tx.CommitmentID = &commitment.ID
	s.Equal(tx, *massMigration)
}

func (s *SyncMMBatchTestSuite) TestSyncBatch_InvalidCommitmentTokenID() {
	tx := testutils.MakeMassMigration(0, 1, 0, 400)
	s.setTxHashAndSign(&tx)

	s.submitInvalidBatch(&tx, func(commitments []models.CommitmentWithTxs) {
		commitments[0].ToMMCommitmentWithTxs().Meta.TokenID = models.MakeUint256(1234)
	})

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	var disputableErr *DisputableError
	err = s.syncCtx.SyncBatch(remoteBatches[0])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(invalidTokenID, disputableErr.Reason)
}

func (s *SyncMMBatchTestSuite) TestSyncBatch_InvalidCommitmentTotalAmount() {
	tx := testutils.MakeMassMigration(0, 1, 0, 400)
	s.setTxHashAndSign(&tx)

	s.submitInvalidBatch(&tx, func(commitments []models.CommitmentWithTxs) {
		commitments[0].ToMMCommitmentWithTxs().Meta.Amount = models.MakeUint256(100)
	})

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	var disputableErr *DisputableError
	err = s.syncCtx.SyncBatch(remoteBatches[0])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(mismatchedTotalAmountMessage, disputableErr.Reason)
}

func (s *SyncMMBatchTestSuite) TestSyncBatch_InvalidCommitmentWithdrawRoot() {
	tx := testutils.MakeMassMigration(0, 1, 0, 400)
	s.setTxHashAndSign(&tx)

	s.submitInvalidBatch(&tx, func(commitments []models.CommitmentWithTxs) {
		commitments[0].ToMMCommitmentWithTxs().WithdrawRoot = common.Hash{1, 2, 3}
	})

	s.recreateDatabase()

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	var disputableErr *DisputableError
	err = s.syncCtx.SyncBatch(remoteBatches[0])
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(invalidWithdrawRootMessage, disputableErr.Reason)
}

func (s *SyncMMBatchTestSuite) submitInvalidBatch(tx models.GenericTransaction, modifier func(commitments []models.CommitmentWithTxs)) {
	pendingBatch, commitments := s.createBatch(tx)

	modifier(commitments)

	err := s.txsCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.GetBackend().Commit()
}

func (s *SyncMMBatchTestSuite) setTxHash(tx *models.MassMigration) {
	hash, err := encoder.HashMassMigration(tx)
	s.NoError(err)
	tx.Hash = *hash
}

func (s *SyncMMBatchTestSuite) setTxHashAndSign(txs ...*models.MassMigration) {
	for i := range txs {
		signMassMigration(s.T(), &s.wallets[txs[i].FromStateID], txs[i])
		s.setTxHash(txs[i])
	}
}

func TestSyncMMBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SyncMMBatchTestSuite))
}
